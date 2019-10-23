resource okta_app_oauth test {
  label         = "testAcc_replace_with_uuid"
  type          = "service"
  redirect_uris = []

  response_types = [
    "token",
  ]

  grant_types = [
    "client_credentials",
  ]

  token_endpoint_auth_method = "client_secret_basic"
}

resource okta_auth_server test {
  name      = "${okta_app_oauth.test.label}"
  audiences = ["api://${okta_app_oauth.test.label}"]
}

resource "okta_auth_server_policy" test {
  auth_server_id   = "${okta_auth_server.test.id}"
  status           = "ACTIVE"
  name             = "Allow Client: ${okta_app_oauth.test.id}"
  description      = "Allow Client: ${okta_app_oauth.test.label}"
  priority         = 1
  client_whitelist = ["${okta_app_oauth.test.id}"]
}

resource okta_auth_server_scope test_read {
  auth_server_id = "${okta_auth_server.test.id}"
  consent        = "REQUIRED"
  description    = "Scope: read"
  name           = "${okta_auth_server.test.name}.read"
}

resource okta_auth_server_scope test_write {
  auth_server_id = "${okta_auth_server.test.id}"
  consent        = "REQUIRED"
  description    = "Scope: write"
  name           = "${okta_auth_server.test.name}.write"
}
