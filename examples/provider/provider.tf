provider "vmds" {
  host      = "https://console.mds.vmware.com"
  // Get the authentication token using api_token
  api_token ="MDS_API_TOKEN"

  //Get teh authentication token using client_id, client_secret and org_id associated with the service account
  client_id = "MDS_CLIENT_ID"
  client_secret = "MDS_CLIENT_SECRET"
  org_id = "MDS_ORG_ID"

}