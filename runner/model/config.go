package model

// Config is a map with key = projectId-serviceId and the value being the routes([]Route)
type Config map[string]Routes // key = projectId-serviceId
