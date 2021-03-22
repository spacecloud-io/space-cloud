package istio

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/spaceuptech/space-cloud/runner/model"
)

func TestIstio_ListSecrets(t *testing.T) {
	type args struct {
		ctx       context.Context
		projectID string
	}
	tests := []struct {
		name              string
		args              args
		secretToBeCreated *v1.Secret
		want              []*model.Secret
		wantErr           bool
	}{
		{
			name: "Get file secret",
			args: args{
				ctx:       context.Background(),
				projectID: "myproject",
			},
			secretToBeCreated: &v1.Secret{
				Type: v1.SecretTypeOpaque,
				ObjectMeta: metav1.ObjectMeta{
					Name:        "Gcp credentials file",
					Namespace:   "myproject",
					Labels:      map[string]string{"app": "space-cloud"},
					Annotations: map[string]string{"rootPath": "/root", "secretType": model.FileType},
				},
				Data: map[string][]byte{"credentials.json": []byte("Planets yell with history at the biological alpha quadrant!")},
			},
			want: []*model.Secret{{
				ID:       "Gcp credentials file",
				Type:     model.FileType,
				RootPath: "/root",
				Data:     map[string]string{"credentials.json": "Planets yell with history at the biological alpha quadrant!"},
			}},
			wantErr: false,
		},
		{
			name: "Get environment secret",
			args: args{
				ctx:       context.Background(),
				projectID: "myproject",
			},
			secretToBeCreated: &v1.Secret{
				Type: v1.SecretTypeOpaque,
				ObjectMeta: metav1.ObjectMeta{
					Name:        "sendgrid API key",
					Namespace:   "myproject",
					Labels:      map[string]string{"app": "space-cloud"},
					Annotations: map[string]string{"rootPath": "", "secretType": model.EnvType},
				},
				Data: map[string][]byte{"SEND_GRID_API_KEY": []byte("Fly without death, and we won’t gather a sun.")},
			},
			want: []*model.Secret{{
				ID:       "sendgrid API key",
				Type:     model.EnvType,
				RootPath: "",
				Data:     map[string]string{"SEND_GRID_API_KEY": "Fly without death, and we won’t gather a sun."},
			}},
			wantErr: false,
		},
		{
			name: "Get docker secret",
			args: args{
				ctx:       context.Background(),
				projectID: "myproject",
			},
			secretToBeCreated: &v1.Secret{
				Type: v1.SecretTypeDockerConfigJson,
				ObjectMeta: metav1.ObjectMeta{
					Name:        "google gcr secret",
					Namespace:   "myproject",
					Labels:      map[string]string{"app": "space-cloud"},
					Annotations: map[string]string{"rootPath": "", "secretType": model.DockerType},
				},
				Data: map[string][]byte{v1.DockerConfigJsonKey: []byte(`{"auths":{"http://gcr.google.com/":{"auth":"c2FtOjEyMw=="}}}`)},
			},
			want: []*model.Secret{{
				ID:       "google gcr secret",
				Type:     model.DockerType,
				RootPath: "",
				Data: map[string]string{
					"username": "sam",
					"password": "123",
					"url":      "http://gcr.google.com/",
				},
			}},
			wantErr: false,
		},
	}

	i := Istio{kube: fake.NewSimpleClientset()}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := i.kube.CoreV1().Secrets(tt.args.projectID).Create(tt.args.ctx, tt.secretToBeCreated, metav1.CreateOptions{}); err != nil {
				t.Errorf("Cannot generate secret required for this test function")
				return
			}
			got, err := i.ListSecrets(tt.args.ctx, tt.args.projectID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListSecrets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListSecrets() mismatch in got & want")
				a, _ := json.MarshalIndent(got, "", " ")
				t.Logf("got= %s", string(a))
				a, _ = json.MarshalIndent(tt.want, "", " ")
				t.Logf("want = %s", string(a))
				return
			}
			if err := i.kube.CoreV1().Secrets(tt.args.projectID).Delete(tt.args.ctx, tt.secretToBeCreated.ObjectMeta.Name, metav1.DeleteOptions{}); err != nil {
				t.Errorf("Cannot delete secret required for this test function")
				return
			}
		})
	}
}
