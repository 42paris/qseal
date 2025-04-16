package kubemodels

type Secret struct {
	APIVersion string            `json:"apiVersion"`
	Data       map[string]string `json:"data"`
	Kind       string            `json:"kind"`
	Metadata   Metadata          `json:"metadata"`
	Type       string            `json:"type"`
}

type Metadata struct {
	Name string `json:"name"`
}
