package upgrade_service

// UpdateMdsClusterVersionRequest represents the request structure for updating the cluster version
type UpdateMdsClusterVersionRequest struct {
	Id            string                                 `json:"id"`
	RequestType   string                                 `json:"requestType"`
	TargetVersion string                                 `json:"targetVersion"`
	Metadata      UpdateMdsClusterVersionRequestMetadata `json:"metadata"`
}

// upgrade_service.UpdateMdsClusterVersionRequestMetadata represents the metadata for the version update request
type UpdateMdsClusterVersionRequestMetadata struct {
	OmitBackup bool `json:"omitBackup"`
}
