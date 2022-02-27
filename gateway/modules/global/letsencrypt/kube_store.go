package letsencrypt

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/caddyserver/certmagic"
	"github.com/spaceuptech/helpers"
	v1 "k8s.io/api/core/v1"
	kubeErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// KubeStore object for storing kube info
type KubeStore struct {
	kubeClient *kubernetes.Clientset
	projectID  string
	path       string
}

// NewKubeStore creates a new instance kube store
func NewKubeStore() (*KubeStore, error) {
	scProject := os.Getenv("LETSENCRYPT_SC_PROJECT")
	if scProject == "" {
		scProject = "space-cloud"
	}
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to create in cluster config for kubernetes store of lets encrypt", err, nil)
	}
	// Create the kubernetes client
	kube, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to create kubernetes client in kubernetes store of lets encrypt", err, nil)
	}

	return &KubeStore{kubeClient: kube, projectID: scProject, path: "certmagic"}, nil
}

// Store stores specified key & value in kube store
func (s *KubeStore) Store(key string, value []byte) error {
	key = s.makeKey(key)
	_, err := s.kubeClient.CoreV1().Secrets(s.projectID).Get(context.TODO(), key, metav1.GetOptions{})
	if kubeErrors.IsNotFound(err) {
		// Create a new Secret
		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Creating secret (%s)", key), nil)
		_, err = s.kubeClient.CoreV1().Secrets(s.projectID).Create(context.TODO(), s.generateSecretValue(key, value), metav1.CreateOptions{})
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("unable to create secret for key (%s) in kubernetes store of lets encrypt", key), err, nil)
		}
		return err
	} else if err == nil {
		// secret already exists...update it!
		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Updating secret (%s)", key), nil)
		_, err = s.kubeClient.CoreV1().Secrets(s.projectID).Update(context.TODO(), s.generateSecretValue(key, value), metav1.UpdateOptions{})
		if err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("unable to update secret for key (%s) in kubernetes store of lets encrypt", key), err, nil)
		}
		return err
	}
	_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unable to set secret for key (%s) in kubernetes store of lets encrypt", key), err, nil)
	return err
}

// Load loads specified key from kube store
func (s *KubeStore) Load(key string) ([]byte, error) {
	key = s.makeKey(key)
	secret, err := s.kubeClient.CoreV1().Secrets(s.projectID).Get(context.TODO(), key, metav1.GetOptions{})
	if err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unable to get secret for key (%s) in kubernetes store of lets encrypt", key), err, nil)
	}
	return secret.Data["value"], nil
}

// Delete deletes specified key from kube store
func (s *KubeStore) Delete(key string) error {
	key = s.makeKey(key)
	err := s.kubeClient.CoreV1().Secrets(s.projectID).Delete(context.TODO(), key, metav1.DeleteOptions{})
	if kubeErrors.IsNotFound(err) || err == nil {
		return nil
	}
	return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unable to delete secret for key (%s) in kubernetes store of lets encrypt", key), err, nil)
}

// Exists check if specified key exists in kube store
func (s *KubeStore) Exists(key string) bool {
	key = s.makeKey(key)
	_, err := s.kubeClient.CoreV1().Secrets(s.projectID).Get(context.TODO(), key, metav1.GetOptions{})
	if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unable to check secret if exists (%s) in kubernetes store of lets encrypt", key), err, nil)
		return false
	}
	return true
}

// List return all key having prefix
func (s *KubeStore) List(prefix string, recursive bool) ([]string, error) {
	// List all secrets
	kubeSecret, err := s.kubeClient.CoreV1().Secrets(s.projectID).List(context.TODO(), metav1.ListOptions{LabelSelector: "app=letsencrypt"})
	if err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to list secrets in kubernetes store of lets encrypt", err, map[string]interface{}{"prefix": prefix, "recursive": recursive})
	}
	keys := make([]string, len(kubeSecret.Items))

	// Modifying SecretValue with empty []byte
	for i, v := range kubeSecret.Items {
		v.Name = s.getOriginalKey(v.Name)
		if strings.HasPrefix(v.Name, prefix) {
			keys[i] = v.Name
		}
	}
	return keys, nil
}

// Stat returns stat for specified key
func (s *KubeStore) Stat(key string) (certmagic.KeyInfo, error) {
	key = s.makeKey(key)
	secret, err := s.kubeClient.CoreV1().Secrets(s.projectID).Get(context.TODO(), key, metav1.GetOptions{})
	if err != nil {
		return certmagic.KeyInfo{}, err
	}

	modifiedTime, err := time.Parse(time.RFC3339Nano, string(secret.Data["modified"]))
	if err != nil {
		return certmagic.KeyInfo{}, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unable to get stat of key (%s) in kubernetes store of lets encrypt cannot parse string to time", key), err, nil)
	}
	size, err := strconv.Atoi(string(secret.Data["size"]))
	if err != nil {
		return certmagic.KeyInfo{}, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unable to get stat of key (%s) in kubernetes store of lets encrypt cannot conver string to integer", key), err, nil)
	}

	return certmagic.KeyInfo{
		Modified:   modifiedTime,
		Key:        secret.Name,
		Size:       int64(size),
		IsTerminal: true,
	}, nil
}

// Lock implements a lock mechanism
func (s *KubeStore) Lock(ctx context.Context, key string) error {
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
			return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to create locker file in kubernetes store of lets encrypt", err, nil)
		}

		// lock file already exists
		info, err := s.Stat(lockFile)
		switch {
		case err != nil:
			return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Unable to access lock file with key (%s) in kubernetes store of lets encrypt", key), err, nil)

		case s.fileLockIsStale(info):
			_ = s.deleteLockFile(lockFile)
			continue

		case time.Since(start) > staleLockDuration*2:
			// should never happen, hopefully
			return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("possible deadlock: %s passed trying to obtain lock for %s", time.Since(start), key), err, nil)

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
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), lockFileExists, nil, nil)
	}

	err := s.Store(filename, []byte("lock"))
	if err != nil {
		_ = helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to create lock file in kubernetes store of lets encrypt", nil, nil)
	}
	return err
}

func (s *KubeStore) deleteLockFile(keyPath string) error {
	err := s.Delete(keyPath)
	if err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), "Unable to delete lock file in lets encrypt", err, nil)
	}
	return nil
}

func (s *KubeStore) generateSecretValue(key string, value []byte) *v1.Secret {
	return &v1.Secret{
		Type: v1.SecretTypeOpaque,
		ObjectMeta: metav1.ObjectMeta{
			Name:      key,
			Namespace: s.projectID,
			Labels:    map[string]string{"app": "letsencrypt"},
		},
		Data: map[string][]byte{
			"value":    value,
			"size":     {byte(len(value))},
			"modified": []byte(time.Now().Format(time.RFC3339Nano)),
		},
	}
}

func (s *KubeStore) makeKey(key string) string {
	newKey := fmt.Sprintf("letsencrypt-%s", key)
	newKey = strings.ReplaceAll(newKey, "/", "--")
	newKey = strings.ReplaceAll(newKey, "_", "---")
	newKey = strings.ReplaceAll(newKey, "@", "----")
	return newKey
}

func (s *KubeStore) getOriginalKey(key string) string {
	// Make sure you replace the maximum number of `-` first. It's in descending order
	oldKey := strings.TrimPrefix(key, "letsencrypt-")
	oldKey = strings.ReplaceAll(oldKey, "----", "@")
	oldKey = strings.ReplaceAll(oldKey, "---", "_")
	oldKey = strings.ReplaceAll(oldKey, "--", "/")
	return oldKey
}
