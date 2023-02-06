package apis

import "testing"

func Test_apps_sort(t *testing.T) {
	tests := []struct {
		name string
		a    apps
		want []string
	}{
		{
			name: "simple sorting test",
			a: apps{
				app{name: "app2", priority: 50},
				app{name: "app1", priority: 100},
				app{name: "app3", priority: 70},
			},
			want: []string{"app1", "app3", "app2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.a.sort()

			for i, a := range tt.a {
				if a.name != tt.want[i] {
					t.Errorf("Invalid order received: wanted - %v; got - %v", tt.want, tt.a)
					return
				}
			}
		})
	}
}
