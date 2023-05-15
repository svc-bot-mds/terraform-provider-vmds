package model

type MdsServiceAccount struct {
	Id     string   `json:"id"`
	Name   string   `json:"name"`
	Status string   `json:"status"`
	Tags   []string `json:"tags"`
}
