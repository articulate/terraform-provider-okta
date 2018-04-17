package okta

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourcePolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePolicyRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Policy name",
				Required:    true,
			},
			"type": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"OKTA_SIGN_ON", "PASSWORD", "MFA_ENROLL", "OAUTH_AUTHORIZATION_POLICY"}, false),
				Description:  "Policy type: OKTA_SIGN_ON, PASSWORD, MFA_ENROLL, or OAUTH_AUTHORIZATION_POLICY",
				Required:     true,
			},
		},
	}
}

func dataSourcePolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Data Source Policy Read %v", d.Get("name").(string))
	client := m.(*Config).oktaClient

	set := false
	currentPolicies, _, err := client.Policies.GetPoliciesByType(d.Get("type").(string))
	if err != nil {
		return fmt.Errorf("[ERROR] Error Listing Policies in Okta: %v", err)
	}
	if currentPolicies != nil {
		for _, policy := range currentPolicies.Policies {
			if policy.Name == d.Get("name").(string) {
				d.SetId(policy.ID)
				set = true
			}
		}
		if set == false {
			return fmt.Errorf("[ERROR] Policy does not exist for type %v", d.Get("type").(string))
		}
	} else {
		return fmt.Errorf("[ERROR] No policies retrieved for policy type %v", d.Get("type").(string))
	}
	return nil
}
