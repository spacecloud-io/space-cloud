package static

import (
	"sort"
	"strings"
	"sync"

	"github.com/spaceuptech/space-cloud/config"
)

// Module is responsible for static
type Module struct {
	sync.RWMutex
	Enabled        bool
	routes         []*config.StaticRoute
	internalRoutes []*config.StaticRoute
}

// Init returns a new instance of the static module wit default values
func Init() *Module {
	return &Module{Enabled: false}
}

// SetConfig set the config required by the static module
func (m *Module) SetConfig(s *config.Static) error {
	m.Lock()
	defer m.Unlock()

	if s != nil {
		sortStatic(s.Routes)
		sortStatic(s.InternalRoutes)
	}

	m.Enabled = true

	if m.internalRoutes == nil {
		m.internalRoutes = []*config.StaticRoute{}
	}

	if m.routes == nil {
		m.routes = []*config.StaticRoute{}
	}

	if s == nil || s.Routes == nil {
		return nil
	}

	m.routes = s.Routes
	return nil
}

func (m *Module) isEnabled() bool {
	m.RLock()
	defer m.RUnlock()

	return m.Enabled
}

// SelectRoute select the rules for a given request
func (m *Module) SelectRoute(host, url string) (*config.StaticRoute, bool) {
	m.RLock()
	defer m.RUnlock()

	for _, route := range m.routes {
		if strings.HasPrefix(url, route.URLPrefix) {
			if route.Host != "" && route.Host != host {
				continue
			}
			return route, true
		}
	}

	for _, route := range m.internalRoutes {
		if strings.HasPrefix(url, route.URLPrefix) {
			if route.Host != "" && route.Host != host {
				continue
			}
			return route, true
		}
	}

	return nil, false
}

func sortStatic(rules []*config.StaticRoute) {
	if rules == nil {
		return
	}
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].URLPrefix < rules[j].URLPrefix
	})
	var splitKey int
	for key, val := range rules {
		if strings.Index(val.URLPrefix, "{") != -1 {
			splitKey = key
			break
		}
	}
	ar1 := rules[:splitKey]
	ar2 := rules[splitKey:]
	rules = append(bubbleSortStatic(ar1), bubbleSortStatic(ar2)...)
}

func bubbleSortStatic(arr []*config.StaticRoute) []*config.StaticRoute {
	var lenArr []int
	for _, value := range arr {
		lenArr = append(lenArr, strings.Count(value.URLPrefix, "/"))
	}

	for i := 0; i < len(lenArr)-1; i++ {
		for j := 0; j < len(lenArr)-i-1; j++ {
			if lenArr[j] < lenArr[j+1] {
				temp := arr[j]
				arr[j] = arr[j+1]
				arr[j+1] = temp
				num := lenArr[j]
				lenArr[j] = lenArr[j+1]
				lenArr[j+1] = num
			}
		}
	}
	return arr
}
