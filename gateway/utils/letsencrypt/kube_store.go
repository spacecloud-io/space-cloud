package letsencrypt

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mholt/certmagic"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	kubeErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KubeStore struct {
	kubeClient *kubernetes.Clientset
	projectId  string
	path       string
}

func NewKubeStore() (*KubeStore, error) {
	scProject := os.Getenv("LETSENCRYPT_SC_PROJECT")
	if scProject == "" {
		scProject = "space_cloud"
	}
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		logrus.Errorf("error in kubernetes store of lets encrypt unable to create in cluster config - %s", err.Error())
		return nil, err
	}
	// Create the kubernetes client
	kube, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		logrus.Errorf("error in kubernetes store of lets encrypt unable to create kubernetes client - %s", err.Error())
		return nil, err
	}

	return &KubeStore{kubeClient: kube, projectId: scProject, path: "certmagic"}, nil
}

func (s *KubeStore) Store(key string, value []byte) error {
	_, err := s.kubeClient.CoreV1().Secrets(s.projectId).Get(key, metav1.GetOptions{})
	if kubeErrors.IsNotFound(err) {
		// Create a new Secret
		logrus.Debugf("Creating secret (%s)", key)
		_, err = s.kubeClient.CoreV1().Secrets(s.projectId).Create(s.generateSecretValue(key, value))
		return err

	} else if err == nil {
		// secret already exists...update it!
		logrus.Debugf("Updating secret (%s)", key)
		_, err = s.kubeClient.CoreV1().Secrets(s.projectId).Update(s.generateSecretValue(key, value))
		return err
	}
	logrus.Errorf("error in kubernetes store of lets encrypt unable to set secret for key (%s) - %s", key, err.Error())
	return err
}

func (s *KubeStore) Load(key string) ([]byte, error) {
	secret, err := s.kubeClient.CoreV1().Secrets(s.projectId).Get(key, metav1.GetOptions{})
	if err != nil {
		logrus.Errorf("error in kubernetes store of lets encrypt unable to get secret for key (%s) - %s", key, err.Error())
		return nil, err
	}
	return secret.Data["value"], nil
}

func (s *KubeStore) Delete(key string) error {
	err := s.kubeClient.CoreV1().Secrets(s.projectId).Delete(key, &metav1.DeleteOptions{})
	if kubeErrors.IsNotFound(err) || err == nil {
		return nil
	}
	logrus.Errorf("error in kubernetes store of lets encrypt unable to delete secret for key (%s) - %s", key, err.Error())
	return err
}

func (s *KubeStore) Exists(key string) bool {
	_, err := s.kubeClient.CoreV1().Secrets(s.projectId).Get(key, metav1.GetOptions{})
	if err != nil {
		logrus.Errorf("error in kubernetes store of lets encrypt unable to check secret if exists (%s) - %s", key, err.Error())
		return false
	}
	return true
}

func (s *KubeStore) List(prefix string, recursive bool) ([]string, error) {
	// List all secrets
	kubeSecret, err := s.kubeClient.CoreV1().Secrets(s.projectId).List(metav1.ListOptions{})
	if err != nil {
		logrus.Errorf("error in kubernetes store of lets encrypt unable to list secrets - %s", err.Error())
		return nil, err
	}
	keys := make([]string, len(kubeSecret.Items))

	// Modifying SecretValue with empty []byte
	for i, v := range kubeSecret.Items {
		if strings.HasPrefix(v.Name, prefix) {
			keys[i] = v.Name
		}
	}
	return keys, nil
}

func (s *KubeStore) Stat(key string) (certmagic.KeyInfo, error) {
	secret, err := s.kubeClient.CoreV1().Secrets(s.projectId).Get(key, metav1.GetOptions{})
	if err != nil {
		return certmagic.KeyInfo{}, err
	}

	modifiedTime, err := time.Parse(time.RFC3339, string(secret.Data["modified"]))
	if err != nil {
		logrus.Errorf("error in kubernetes store of lets encrypt unable to parse string to time for key (%s) - %s", key, err.Error())
		return certmagic.KeyInfo{}, fmt.Errorf("unable to parse string to time - %v", err)
	}
	size, err := strconv.Atoi(string(secret.Data["size"]))
	if err != nil {
		logrus.Errorf("error in kubernetes store of lets encrypt unable to convert string to integer for key (%s) - %s", key, err.Error())
		return certmagic.KeyInfo{}, fmt.Errorf("unable to convert string to integer - %v", err)
	}

	return certmagic.KeyInfo{
		Modified:   modifiedTime,
		Key:        secret.Name,
		Size:       int64(size),
		IsTerminal: true,
	}, nil
}

func (s *KubeStore) Lock(key string) error {
	start := time.Now()
	lockFile := s.lockFileName(key)

	for {
		err := s.createLockFile(lockFile)
		if err == nil {
			// got the lock
			return nil
		}

		if err.Error() != lockFileExists {
			// unexpected error
			logrus.Errorf("error in kubernetes store of lets encrypt - %s", err.Error())
			return fmt.Errorf("creating lock file: %+v", err)
		}

		// lock file already exists
		info, err := s.Stat(lockFile)
		switch {
		case err != nil:
			return fmt.Errorf("accessing lock file: %v", err)

		case s.fileLockIsStale(info):
			_ = s.deleteLockFile(lockFile)
			continue

		case time.Since(start) > staleLockDuration*2:
			// should never happen, hopefully
			return fmt.Errorf("possible deadlock: %s passed trying to obtain lock for %s", time.Since(start), key)

		default:
			// lockfile exists and is not stale;
			// just wait a moment and try again
			time.Sleep(fileLockPollInterval)

		}
	}
}

// Unlock releases the lock for name.
func (s *KubeStore) Unlock(key string) error {
	return s.deleteLockFile(s.lockFileName(key))
}

func (s *KubeStore) String() string {
	return "KubeStore:" + s.path
}

func (s *KubeStore) lockFileName(key string) string {
	return filepath.Join(s.lockDir(), fmt.Sprintf("%s.lock", StorageKeys.Safe(key)))
}

func (s *KubeStore) lockDir() string {
	return filepath.Join(s.path, "locks")
}

func (s *KubeStore) fileLockIsStale(info certmagic.KeyInfo) bool {
	return time.Since(info.Modified) > staleLockDuration
}

func (s *KubeStore) createLockFile(filename string) error {
	exists := s.Exists(filename)
	if exists {
		return fmt.Errorf(lockFileExists)
	}

	err := s.Store(filename, []byte("lock"))
	if err != nil {
		logrus.Errorf("error while creating lock file in lets encrypt - %v", err)
	}
	return err
}

func (s *KubeStore) deleteLockFile(keyPath string) error {
	err := s.Delete(keyPath)
	if err != nil {
		logrus.Errorf("error while deleting lock file in lets encrypt - %v", err)
		return fmt.Errorf("error while deleting lock file in lets encrypt - %v", err)
	}
	return nil
}

func (s *KubeStore) generateSecretValue(key string, value []byte) *v1.Secret {
	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      key,
			Namespace: s.projectId,
		},
		Data: map[string][]byte{
			"value":    value,
			"size":     {byte(len(value))},
			"modified": []byte(time.Now().String()),
		},
	}
}
