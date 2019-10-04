resource "okta_user_base_schema" "firstName" {
  index       = "firstName"
  master      = "PROFILE_MASTER"
  permissions = "READ_ONLY"
  title       = "First name"
  type        = "string"
  min_length  = 1
  max_length  = 50
}
