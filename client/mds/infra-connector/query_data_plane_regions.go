package infra_connector

type DataPlaneRegionsQuery struct {
	Provider  string `schema:"provider"`
	CPU       string `schema:"cpu"`
	Memory    string `schema:"memory"`
	Storage   string `schema:"storage"`
	NodeCount string `schema:"nodeCount"`
	OrgId     string `schema:"orgId,omitempty"`
}
