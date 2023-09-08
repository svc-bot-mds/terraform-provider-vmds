package model

type MdsByocCloudProviderRegion struct {
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	ShortName string   `json:"shortName"`
	Regions   []string `json:"regions"`
}
