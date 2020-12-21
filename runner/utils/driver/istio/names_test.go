package istio

import "testing"

func Test_splitInternalServiceName(t *testing.T) {
	type args struct {
		n string
	}
	tests := []struct {
		name          string
		args          args
		wantServiceID string
		wantVersion   string
	}{
		{
			name: "Proper internal service name",
			args: args{
				n: getInternalServiceName("greeter", "v1"),
			},
			wantServiceID: "greeter",
			wantVersion:   "v1",
		},
		{
			name: "Proper internal service name, where service name contains hypen(-)",
			args: args{
				n: getInternalServiceName("http-bin", "v1"),
			},
			wantServiceID: "http-bin",
			wantVersion:   "v1",
		},
		{
			name: "Improper internal service name",
			args: args{
				n: "httpbin-v1",
			},
			wantServiceID: "",
			wantVersion:   "",
		},
		{
			name: "Improper internal service name",
			args: args{
				n: "http-bin-v1",
			},
			wantServiceID: "",
			wantVersion:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotServiceID, gotVersion := splitInternalServiceName(tt.args.n)
			if gotServiceID != tt.wantServiceID {
				t.Errorf("splitInternalServiceName() gotServiceID = %v, want %v", gotServiceID, tt.wantServiceID)
			}
			if gotVersion != tt.wantVersion {
				t.Errorf("splitInternalServiceName() gotVersion = %v, want %v", gotVersion, tt.wantVersion)
			}
		})
	}
}

func Test_checkIfInternalServiceDomain(t *testing.T) {
	type args struct {
		projectID             string
		serviceID             string
		internalServiceDomain string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Proper internal service domain",
			args: args{
				projectID:             "myProject",
				serviceID:             "greeter",
				internalServiceDomain: getInternalServiceDomain("myProject", "greeter", "v2"),
			},
			want: true,
		},
		{
			name: "Improper internal service domain",
			args: args{
				projectID:             "myProject",
				serviceID:             "greeter",
				internalServiceDomain: "greeter.myProject.svc.cluster.local",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkIfInternalServiceDomain(tt.args.projectID, tt.args.serviceID, tt.args.internalServiceDomain); got != tt.want {
				t.Errorf("checkIfInternalServiceDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
