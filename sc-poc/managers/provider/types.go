package provider

import "sort"

type (
	workspaceset struct {
		name      string
		providers map[string]any
	}

	providers []provider
	provider  struct {
		name     string
		priority int
	}
)

func (a providers) sort() {
	sort.SliceStable(a, func(i, j int) bool {
		return a[i].priority > a[j].priority
	})
}
