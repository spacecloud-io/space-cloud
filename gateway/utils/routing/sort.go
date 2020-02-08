package routing

import (
	"github.com/spaceuptech/space-cloud/gateway/config"
	"sort"
	"strings"
)

func sortRoutes(rules []*config.Route) {
	var newRoutes routes
	newRoutes = rules
	sort.Stable(newRoutes)
}

type routes []*config.Route

func (a routes) Len() int      { return len(a) }
func (a routes) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a routes) Less(i, j int) bool {
	return len(strings.Split(a[i].Source.URL, "/")) >= len(strings.Split(a[j].Source.URL, "/"))
}
