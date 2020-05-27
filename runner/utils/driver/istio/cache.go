package istio

import (
	"context"
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type cache struct {
	db   *bolt.DB
	kube *kubernetes.Clientset
}

func newCache(kube *kubernetes.Clientset) (*cache, error) {
	db, err := bolt.Open("./runner.db", 0666, nil)
	if err != nil {
		return nil, err
	}

	return &cache{db: db, kube: kube}, nil
}

// TODO: Move the update deployment logic in here
func (c *cache) setDeployment(ns, name string, deployment *v1.Deployment) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		// There will be a bucket for each namespace
		bucket, err := tx.CreateBucketIfNotExists([]byte(ns))
		if err != nil {
			return err
		}

		// Marshal the deployment
		data, err := proto.Marshal(deployment)
		if err != nil {
			logrus.Errorf("Could not cache deployment (%s) in namespace (%s) - %s", name, ns, err.Error())
			return err
		}

		// Store it in the bucket
		return bucket.Put([]byte(getCacheDeploymentKey(name)), data)
	})
}

func (c *cache) getDeployment(ctx context.Context, ns, name string) (*v1.Deployment, error) {
	deployment := new(v1.Deployment)
	var foundInCache bool
	var err error

	if err := c.db.View(func(tx *bolt.Tx) error {
		// There will be a bucket for each namespace
		bucket := tx.Bucket([]byte(ns))

		// Exit if the bucket does not exist. This signifies that the cache is empty
		if bucket == nil {
			return nil
		}

		// Get the deployment from the cache. The value returned will be nil if the key doesn't not exist.
		// In this case we simply return nil and load the deployment config from kubernetes.
		val := bucket.Get([]byte(getCacheDeploymentKey(name)))
		if val == nil {
			return nil
		}

		// Attempt to unmarshal the stored deployment config. If the unmarshal is successful, set the found in cache flag.
		// If an error occurred in the unmarshal process, simply log the error. We will return nil in this case as well or else
		// the error will be propagated to the client as well. For now we simply read the config again directly from kubernetes
		if err := proto.Unmarshal(val, deployment); err != nil {
			logrus.Errorf("Could not unmarshal deployment (%s) in namespace (%s) - %s", name, ns, err.Error())
			return nil
		}
		foundInCache = true
		return nil
	}); err != nil {
		return nil, err
	}

	// Query kubernetes if the deployment wasn't found in the cache. Also, update the cache with the deployment received.
	if !foundInCache {
		logrus.Debugf("Service (%s) not found in cache... Querying kubernetes", name)
		deployment, err = c.kube.AppsV1().Deployments(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		// Attempt to set the deployment in the cache. We are ignoring the error since it's okay if some error occurred.
		// The deployment will eventually get cached so it's alright
		_ = c.setDeployment(ns, name, deployment)
	}

	return deployment, nil
}

func getCacheDeploymentKey(name string) string {
	return fmt.Sprintf("kube/deployments/%s", name)
}
