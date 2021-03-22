package istio

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/spaceuptech/helpers"
	v1 "k8s.io/api/core/v1"
	kubeErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spaceuptech/space-cloud/runner/model"
)

// CreateSecret is used to upsert secret
func (i *Istio) CreateSecret(ctx context.Context, projectID string, secretObj *model.Secret) error {
	// check whether the oldSecret type is correct!
	if secretObj.Type != model.FileType && secretObj.Type != model.EnvType && secretObj.Type != model.DockerType {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Invalid old secret type (%s) provided", secretObj.Type), nil, nil)
	}

	oldSecret, err := i.kube.CoreV1().Secrets(projectID).Get(ctx, secretObj.ID, metav1.GetOptions{})
	if kubeErrors.IsNotFound(err) {
		// Create a new Secret
		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), fmt.Sprintf("Creating secret (%s)", secretObj.ID), nil)
		newSecret, err := generateSecret(ctx, projectID, secretObj)
		if err != nil {
			return err
		}

		_, err = i.kube.CoreV1().Secrets(projectID).Create(ctx, newSecret, metav1.CreateOptions{})
		return err

	} else if err == nil {
		// oldSecret already exists...update it!
		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), fmt.Sprintf("Updating secret (%s)", secretObj.ID), nil)
		oldSecretType := oldSecret.Annotations["secretType"]
		if oldSecret.Annotations["secretType"] != secretObj.Type {
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Secret type mismatch. Wanted - %s; Got - %s", oldSecretType, secretObj.Type), nil, nil)
		}
		newSecret, err := generateSecret(ctx, projectID, secretObj)
		if err != nil {
			return err
		}
		_, err = i.kube.CoreV1().Secrets(projectID).Update(ctx, newSecret, metav1.UpdateOptions{})
		return err
	}
	return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Failed to create secret (%s)", secretObj.ID), err, nil)
}

// ListSecrets lists all the secrets in the provided name-space!
func (i *Istio) ListSecrets(ctx context.Context, projectID string) ([]*model.Secret, error) {
	// List all secrets
	kubeSecret, err := i.kube.CoreV1().Secrets(projectID).List(ctx, metav1.ListOptions{LabelSelector: "app=space-cloud"})
	if err != nil {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Failed to fetch list of secrets", err, nil)
	}
	listOfSecrets := make([]*model.Secret, len(kubeSecret.Items))

	// Modifying SecretValue with empty []byte
	for i, v := range kubeSecret.Items {
		s := &model.Secret{
			ID:       v.ObjectMeta.Name,
			Type:     v.ObjectMeta.Annotations["secretType"],
			RootPath: v.ObjectMeta.Annotations["rootPath"],
			Data:     make(map[string]string, len(v.Data)),
		}
		if s.Type == model.DockerType {
			value, ok := v.Data[v1.DockerConfigJsonKey]
			if !ok {
				return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Docker secret not made according to the space cloud format", fmt.Errorf("key (%s) does not exists in secret (%s) with type docker", v1.DockerConfigJsonKey, v.ObjectMeta.Name), nil)
			}
			obj := map[string]interface{}{}
			if err := json.Unmarshal(value, &obj); err != nil {
				return nil, err
			}

			for key, data := range obj["auths"].(map[string]interface{}) {
				tempObj := data.(map[string]interface{})["auth"].(string)
				decodedString, err := b64.StdEncoding.DecodeString(tempObj)
				if err != nil {
					return nil, err
				}
				arr := strings.Split(string(decodedString), ":")
				if len(arr) < 2 {
					return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), "Docker secret not made according to the space cloud format", fmt.Errorf("auth value for secret (%s) with type docker is not seperated by (:)", v.ObjectMeta.Name), nil)
				}
				s.Data["username"] = arr[0]
				s.Data["password"] = arr[1]
				s.Data["url"] = key
				break
			}
		} else {
			for k1, data := range v.Data {
				s.Data[k1] = string(data)
			}
		}
		listOfSecrets[i] = s
	}
	return listOfSecrets, nil
}

// DeleteSecret is used to delete secrets!
func (i *Istio) DeleteSecret(ctx context.Context, projectID string, secretName string) error {
	err := i.kube.CoreV1().Secrets(projectID).Delete(ctx, secretName, metav1.DeleteOptions{})
	if kubeErrors.IsNotFound(err) || err == nil {
		return nil
	}
	return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Failed to delete secret (%s)", secretName), err, nil)
}

// SetFileSecretRootPath is used to set the file secret root path
func (i *Istio) SetFileSecretRootPath(ctx context.Context, projectID string, secretName, rootPath string) error {
	if secretName == "" || rootPath == "" {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Empty secret name (%s) or root path (%s) provided", secretName, rootPath), nil, nil)
	}
	// Get secret and then check type
	kubeSecret, err := i.kube.CoreV1().Secrets(projectID).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// update root path
	switch kubeSecret.Type {
	case v1.SecretTypeDockerConfigJson:
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "set root path operation cannot be performed on secrets with type docker", nil, nil)
	case v1.SecretTypeOpaque:
		kubeSecret.Annotations["rootPath"] = rootPath
	default:
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid secret type - %s", kubeSecret.Type), nil, nil)
	}

	// Update the secret
	_, err = i.kube.CoreV1().Secrets(projectID).Update(ctx, kubeSecret, metav1.UpdateOptions{})
	return err
}

// SetKey adds a new secret key-value pair
func (i *Istio) SetKey(ctx context.Context, projectID string, secretName string, secretKey string, secretValObj *model.SecretValue) error {
	if secretName == "" || secretValObj.Value == "" {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("key/value not provided; got (%s,%s)", secretName, secretValObj.Value), nil, nil)
	}

	// Get secret and then check type
	kubeSecret, err := i.kube.CoreV1().Secrets(projectID).Get(ctx, secretName, metav1.GetOptions{})

	if kubeErrors.IsNotFound(err) {
		return err
	} else if err == nil {
		// Add secret key-value
		switch kubeSecret.Type {
		case v1.SecretTypeDockerConfigJson:
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), "set key operation cannot be performed on secrets with type docker", nil, nil)
		case v1.SecretTypeOpaque:
			if kubeSecret.Data == nil {
				kubeSecret.Data = make(map[string][]byte, 1)
			}
			kubeSecret.Data[secretKey] = []byte(secretValObj.Value)
		default:
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid secret type - %s", kubeSecret.Type), nil, nil)
		}

		// Update the secret
		_, err := i.kube.CoreV1().Secrets(projectID).Update(ctx, kubeSecret, metav1.UpdateOptions{})
		return err
	}
	return err
}

// DeleteKey is used to delete a key from the secret!
func (i *Istio) DeleteKey(ctx context.Context, projectID string, secretName string, secretKey string) error {
	// Get secret
	kubeSecret, err := i.kube.CoreV1().Secrets(projectID).Get(ctx, secretName, metav1.GetOptions{})

	if kubeErrors.IsNotFound(err) {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("secret with name (%s) does not exist", secretName), err, nil)
	} else if err == nil {
		// Check the type of secret (docker/opaque)
		switch kubeSecret.Type {
		case v1.SecretTypeDockerConfigJson:
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), "delete key operation cannot be performed on secrets with type docker", nil, nil)
		case v1.SecretTypeOpaque:
			delete(kubeSecret.Data, secretKey)
		default:
			return helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid secret type - %s", kubeSecret.Type), nil, nil)
		}

		// Update the secret
		_, err := i.kube.CoreV1().Secrets(projectID).Update(ctx, kubeSecret, metav1.UpdateOptions{})
		return err
	}
	return err
}

// helper function
func generateSecret(ctx context.Context, projectID string, secret *model.Secret) (*v1.Secret, error) {
	encodedData := map[string][]byte{}
	var typeOfSecret v1.SecretType

	// Check what type of secret is to be created: file/env/docker
	switch secret.Type {
	case model.FileType, model.EnvType:
		typeOfSecret = v1.SecretTypeOpaque
		for k, v := range secret.Data {
			encodedData[k] = []byte(v)
		}
	case model.DockerType:
		username, p1 := secret.Data["username"]
		password, p2 := secret.Data["password"]
		url, p3 := secret.Data["url"]

		if !p1 || !p2 || !p3 {
			return nil, errors.New("incorrect secret value provided for secret type docker")
		}

		typeOfSecret = v1.SecretTypeDockerConfigJson
		authSecret := username + ":" + password
		encAuthSecret := b64.StdEncoding.EncodeToString([]byte(authSecret))
		// ref: https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/#registry-secret-existing-credentials
		dockerJSON := map[string]interface{}{
			"auths": map[string]interface{}{
				url: map[string]string{
					"auth": encAuthSecret,
				},
			},
		}
		data, _ := json.Marshal(dockerJSON)
		encodedData[v1.DockerConfigJsonKey] = []byte(data)
	default:
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid secret type (%s) provided", secret.Type), nil, nil)
	}
	return &v1.Secret{
		Type: typeOfSecret,
		ObjectMeta: metav1.ObjectMeta{
			Name:      secret.ID,
			Namespace: projectID,
			Labels: map[string]string{
				"app":                          "space-cloud",
				"app.kubernetes.io/name":       secret.ID,
				"app.kubernetes.io/managed-by": "space-cloud",
			},
			Annotations: map[string]string{"rootPath": secret.RootPath, "secretType": secret.Type},
		},
		Data: encodedData,
	}, nil
}

func (i *Istio) getSecrets(ctx context.Context, service *model.Service) (map[string]*v1.Secret, error) {
	listOfSecrets := map[string]*v1.Secret{}
	tasks := service.Tasks
	for _, task := range tasks {
		for _, secretName := range task.Secrets {
			if _, p := listOfSecrets[secretName]; p {
				continue
			}
			secrets, err := i.kube.CoreV1().Secrets(service.ProjectID).Get(ctx, secretName, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}
			listOfSecrets[secretName] = secrets
		}
	}
	return listOfSecrets, nil
}
