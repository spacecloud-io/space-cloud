package letsencrypt

import "github.com/spaceuptech/space-cloud/gateway/utils"

type domainMapping map[string][]string // key is project id and array is domain

func (d domainMapping) setProjectDomains(project string, domains []string) {
	d[project] = domains
}

func (d domainMapping) deleteProject(project string) {
	delete(d, project)
}

func (d domainMapping) getUniqueDomains() []string {
	var domains []string

	// Iterate over all projects
	for _, v := range d {
		// Iterate over all domains in project
		for _, domain := range v {
			if !utils.StringExists(domains, domain) {
				domains = append(domains, domain)
			}
		}
	}

	return domains
}
