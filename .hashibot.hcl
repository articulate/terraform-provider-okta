poll "label_issue_migrater" "provider_migrater" {
    schedule = "0 20 * * * *"
    new_owner = "terraform-providers" 
    repo_prefix = "terraform-provider-"
    label_prefix = "provider/"
    issue_header = <<-EOF
    _This issue was originally opened by @${var.user} as ${var.repository}#${var.issue_number}. The original body of the issue is below._
    
    <hr>
    
    EOF
    migrated_comment = "This repository has changed owners! This issue has been automatically migrated to ${var.repository}#${var.issue_number}."
}
