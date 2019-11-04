package syncman

import (
	"errors"
	"math"
	"sort"

	"github.com/getlantern/deepcopy"
	"github.com/hashicorp/consul/api"

	"github.com/spaceuptech/space-cloud/config"
)

func (s *Manager) setSpaceCloudInstances(nodes memRange) {
	s.lock.Lock()
	defer s.lock.Unlock()

	var passingNodes memRange

	// Filter out failing nodes
	for _, node := range nodes {
		if isCheckNotPassing(node.Checks) {
			continue
		}

		passingNodes = append(passingNodes, node)
	}

	// Sort and store
	sort.Stable(passingNodes)
	s.services = passingNodes
}

func isCheckNotPassing(checks api.HealthChecks) bool {
	for _, check := range checks {
		if check.Status != api.HealthPassing {
			return true
		}
	}
	return false
}

func (s *Manager) setProjectConfig(conf *config.Project) {
	for i, p := range s.projectConfig.Projects {
		if p.ID == conf.ID {
			s.projectConfig.Projects[i] = conf
			return
		}
	}

	s.projectConfig.Projects = append(s.projectConfig.Projects, conf)
}

func (s *Manager) delete(projectID string) {
	for i, p := range s.projectConfig.Projects {
		if p.ID == projectID {
			remove(s.projectConfig.Projects, i)
			break
		}
	}
}

func remove(s []*config.Project, i int) []*config.Project {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}

type memRange []*api.ServiceEntry

func (a memRange) Len() int           { return len(a) }
func (a memRange) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a memRange) Less(i, j int) bool { return a[i].Service.Address < a[j].Service.Address }

func calcTokens(n int, tokens int, i int) (start int, end int) {
	tokensPerMember := int(math.Ceil(float64(tokens) / float64(n)))
	start = tokensPerMember * i
	end = start + tokensPerMember - 1
	if end > tokens {
		end = tokens - 1
	}
	return
}

func calcIndex(token, totalTokens, n int) int {
	bucketSize := totalTokens / n
	return token / bucketSize
}

// getConfigWithoutLock returns the config present in the state
func (s *Manager) getConfigWithoutLock(projectID string) (*config.Project, error) {
	// Iterate over all projects stored
	for _, p := range s.projectConfig.Projects {
		if projectID == p.ID {
			proj := new(config.Project)
			deepcopy.Copy(proj, p)
			return proj, nil
		}
	}

	return nil, errors.New("given project is not present in state")
}
