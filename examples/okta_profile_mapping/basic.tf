resource okta_profile_mapping test {
  source_id          = "0oaes3ebcorItLikJ0h7"
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
