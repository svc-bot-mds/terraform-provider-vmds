package customer_metadata

type MdsSvcAccountUpdateRequest struct {
	Tags      []string `json:"tags"`
	PolicyIds []string `json:"policyIds"`
}
