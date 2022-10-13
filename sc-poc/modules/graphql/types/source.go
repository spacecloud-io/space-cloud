package types

// Source is a graphql remote endpoint that we want to federate
type Source struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
