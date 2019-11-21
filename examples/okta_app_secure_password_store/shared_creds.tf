resource "okta_app_secure_password_store" "test" {
  label              = "testAcc_replace_with_uuid"
  url                = "https://example.com/users/sign_in"
  password_field     = "password"
  username_field     = "user"
  credentials_scheme = "SHARED_USERNAME_AND_PASSWORD"
  shared_username    = "testAcc_replace_with_uuid"
  shared_password    = "secret stuff"
}
