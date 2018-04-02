package okta

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePolicies() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyCreate,
		Read:   resourcePolicyRead,
		Update: resourcePolicyUpdate,
		Delete: resourcePolicyDelete,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy ID (generated)",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Policy Type",
			},

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Policy Name",
			},
		},
	}
}

func resourcePolicyCreate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourcePolicyRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourcePolicyUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourcePolicyDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
