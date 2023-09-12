package model

type MdsDataPlaneRegion struct {
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	ShortName string   `json:"shortName"`
	Regions   []string `json:"regions"`
}
