resource "okta_user_base_schema" "login" {
  index       = "login"
  title       = "Username"
  type        = "string"
  master      = "PROFILE_MASTER"
  permissions = "READ_ONLY"
  required    = true
  min_length  = 4
  max_length  = 70
}
