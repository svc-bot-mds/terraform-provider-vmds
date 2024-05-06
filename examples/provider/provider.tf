provider "vmds" {
  host      = "https://console.mds.vmware.com"

  //Get the authentication with "username and password"
  username = "< Username >"
  password = " < Password > "

  type = "user_creds"
}