package istio

import (
	"fmt"

	b64 "encoding/base64"
	v1 "k8s.io/api/core/v1"
	kubeErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-cloud/runner/model"
)

// helper function
func generateSecret(secret *model.Secret, projectID string) *v1.Secret {
	var encodedData map[string][]byte
	var typeOfSecret v1.SecretType

	// Base64 encoding!
	for k, v := range secret.Data {
		encValue := b64.StdEncoding.EncodeToString([]byte(v))
		encodedData[k] = []byte(encValue)
	}
	// Check what type of secret is to be created: file/env/docker
	if secret.SecretType == "file" {
		// any default path??
		typeOfSecret = "Opaque"
		return &v1.Secret{Type: typeOfSecret, TypeMeta: metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"}, ObjectMeta: metav1.ObjectMeta{Name: secret.SecretName, Namespace: projectID, Annotations: map[string]string{"annotationRoute": secret.RoutePath}}, Data: encodedData}

	} else if secret.SecretType == "env" {
		typeOfSecret = "Opaque"

	} else {
		// for secretType : docker
		typeOfSecret = v1.SecretTypeDockerConfigJson
	}

	return &v1.Secret{Type: typeOfSecret, TypeMeta: metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"}, ObjectMeta: metav1.ObjectMeta{Name: secret.SecretName, Namespace: projectID}, Data: encodedData}
}

// CreateSecret is used to upsert secret
func (i *Istio) CreateSecret(secretObj *model.Secret, projectID string) error {
	// check whether the secret type is correct!
	if secretObj.SecretType != "file" && secretObj.SecretType != "env" && secretObj.SecretType != "docker" {
		return fmt.Errorf("invalid secret type provided: (%s)", secretObj.SecretType)
	}

	_, err := i.kube.CoreV1().Secrets(projectID).Get(secretObj.SecretName, metav1.GetOptions{})
	if kubeErrors.IsNotFound(err) {
		// Create a new Secret
		logrus.Debugf("Creating secret for %s ", secretObj.SecretName)

		newSecret := generateSecret(secretObj, projectID)
		_, err := i.kube.CoreV1().Secrets(projectID).Create(newSecret)
		if err != nil {
			return fmt.Errorf("Error creating secret: (%s)", err)
		}

	} else if err == nil {
		// secret already exists...update it!
		logrus.Debugf("Updating secret for %s ", secretObj.SecretName)
		newSecret := generateSecret(secretObj, projectID)
		_, err = i.kube.CoreV1().Secrets(projectID).Update(newSecret)
		if err != nil {
			return fmt.Errorf("Error updating secret : (%s)", err)
		}

	} else {
		return fmt.Errorf("Listing Secrets failed with error: (%s)", err)
	}
	return nil
}

// ListSecrets lists all the secrets in the provided name-space!
func (i *Istio) ListSecrets(secretObj *model.Secret, projectID string) (*v1.SecretList, error) {

	kubeSecret, err := i.kube.CoreV1().Secrets(projectID).List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Error retrieving secrets😫: (%s)", err)
	}
	// Modifying SecretValue with empty []byte
	for k := range kubeSecret.Items {
		for k1 := range kubeSecret.Items[k].Data {
			for key := range k1 {
				kubeSecret.Items[key].Data[k1] = make([]byte, 0)
			}
		}
	}
	return kubeSecret, nil
}

// DelSecrets is used to delete secrets!
func (i *Istio) DelSecrets(secretObj *model.Secret, projectID string) error {
	err := i.kube.CoreV1().Secrets(projectID).Delete(secretObj.SecretName, &metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("Error deleting secret: (%s)", err)
	}
	return nil
}

// SetKey adds a new secret key-value pair
func (i *Istio) SetKey(secretValObj *model.SecretValue, projectID string, secretName string, secretKey string) error {
	// encoding secret value to base64
	encSecret := b64.StdEncoding.EncodeToString([]byte(secretValObj.Value))
	//Get secret and then check type
	kubeSecret, err := i.kube.CoreV1().Secrets(projectID).Get(secretName, metav1.GetOptions{})

	if kubeErrors.IsNotFound(err) {
		return fmt.Errorf("Secret not found error: (%s)", err)
	} else if err == nil {
		//Add secret key-value
		switch kubeSecret.Type {
		case v1.SecretTypeDockerConfigJson:
			kubeSecret.Data[v1.DockerConfigJsonKey] = []byte(encSecret)
		default:
			kubeSecret.Data[secretKey] = []byte(encSecret)
		}

		//Create a new secret(updated)
		_, err := i.kube.CoreV1().Secrets(projectID).Update(kubeSecret)
		if err != nil {
			return fmt.Errorf("Error setting secretKey: (%s)", err)
		}
	} else {
		// return unknown error
		return err
	}
	return nil
}

// DelKey is used to delete a key from the secret!
func (i *Istio) DelKey(projectID string, secretName string, secretKey string) error {
	// Get secret
	kubeSecret, err := i.kube.CoreV1().Secrets(projectID).Get(secretName, metav1.GetOptions{})

	if kubeErrors.IsNotFound(err) {
		return fmt.Errorf("Secret not found error: (%s)", err)
	} else if err == nil {
		// check the type of secret (docker/opaque)
		switch kubeSecret.Type {
		case v1.SecretTypeDockerConfigJson:
			delete(kubeSecret.Data, v1.DockerConfigJsonKey)
		default:
			// Iterate over map and delete the key
			for k := range kubeSecret.Data {
				if k == secretKey {
					delete(kubeSecret.Data, k)
					break
				}
			}
		}
		// Update the secret
		_, err := i.kube.CoreV1().Secrets(projectID).Update(kubeSecret)
		if err != nil {
			return fmt.Errorf("Error deleting secretKey: (%s)", err)
		}
	} else {
		return err
	}
	return nil
}
