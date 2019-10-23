resource okta_profile_mapping test {
  source_id = "${okta_idp_saml.test.id}"
  delete_when_absent = true

  mappings {
    id          = "nickName"
    expression  = "user.nickName"
    push_status = "PUSH"
  }

  mappings {
    id         = "fullName"
    expression = "user.firstName + user.lastName"
  }
}

resource okta_idp_saml test {
  name                     = "testAcc_replace_with_uuid"
  acs_binding              = "HTTP-POST"
  acs_type                 = "INSTANCE"
  sso_url                  = "https://idp.example.com"
  sso_destination          = "https://idp.example.com"
  sso_binding              = "HTTP-POST"
  username_template        = "idpuser.email"
  kid                      = "${okta_idp_saml_key.test.id}"
  issuer                   = "https://idp.example.com"
  request_signature_scope  = "REQUEST"
  response_signature_scope = "ANY"
}
