package model

type MdsClusterMetaData struct {
	Id          string              `json:"id"`
	Name        string              `json:"name"`
	Provider    string              `json:"provider"`
	ServiceType string              `json:"serviceType"`
	Status      string              `json:"status"`
	VHosts      []MdsVhosts         `json:"vhosts,omitempty"`
	Queues      []MdsQueuesModel    `json:"queues,omitempty"`
	Exchanges   []MdsExchangesModel `json:"exchanges,omitempty"`
	Bindings    []MdsBindingsModel  `json:"bindings,omitempty"`
}

type MdsVhosts struct {
	Name string `json:"name"`
}

type MdsQueuesModel struct {
	Name  string `json:"name"`
	VHost string `json:"vhost"`
}

type MdsExchangesModel struct {
	Name  string `json:"name"`
	VHost string `json:"vhost"`
}

type MdsBindingsModel struct {
	Source          string `json:"source"`
	VHost           string `json:"vhost"`
	RoutingKey      string `json:"routingKey"`
	Destination     string `json:"destination"`
	DestinationType string `json:"destinationType"`
}
