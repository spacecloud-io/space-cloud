package provider

import "testing"

func Test_providers_sort(t *testing.T) {
	tests := []struct {
		name string
		p    providers
		want []string
	}{
		{
			name: "simple sorting test",
			p: providers{
				provider{name: "provider2", priority: 50},
				provider{name: "provider1", priority: 100},
				provider{name: "provider3", priority: 70},
			},
			want: []string{"provider1", "provider3", "provider2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.sort()

			for i, a := range tt.p {
				if a.name != tt.want[i] {
					t.Errorf("Invalid order received: wanted - %v; got - %v", tt.want, tt.p)
					return
				}
			}
		})
	}
}
