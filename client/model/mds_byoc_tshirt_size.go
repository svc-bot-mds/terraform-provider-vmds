package model

type MdsByocTshirtSize struct {
	Name     string `json:"name"`
	Nodes    int64  `json:"nodes"`
	Provider string `json:"provider"`
	Storage  string `json:"storage"`
	Type     string `json:"type"`
}
