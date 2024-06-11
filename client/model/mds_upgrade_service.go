package model

type MdsUpgradeServiceRequest struct {
}

// UpdateMdsClusterVersionRequest represents the request structure for updating the cluster version
type UpdateMdsClusterVersionRequest struct {
	Id            string `json:"id"`
	RequestType   string `json:"requestType"`
	TargetVersion string `json:"targetVersion"`
	Metadata      struct {
		OmitBackup bool `json:"omitBackup"`
	} `json:"metadata"`
}

// UpdateMdsClusterVersionResponse represents the response structure for updating the cluster version
type UpdateMdsClusterVersionResponse struct {
	Success bool `json:"success"`
}
