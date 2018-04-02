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
				Description: "Policy Type: OKTA_SIGN_ON, PASSWORD, MFA_ENROLL, or OAUTH_AUTHORIZATION_POLICY",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Policy Name",
			},
			"system": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Set to true if Policy is a System Policy",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Policy Name",
			},
			"priority": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Policy Priority",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ACTIVE",
				Description: "Policy Status",
			},
			"users_condition": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of Users to be Included or Excluded for the Policy",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"include": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"exclude": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"groups_condition": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of Groups to be Included or Excluded for the Policy",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"include": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"exclude": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"authProvider_conditon": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Authentication Provider for the Policy",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "OKTA",
							Description: "Policy Status: OKTA or ACTIVE_DIRECTORY",
						},
						"include": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of Active Directory Integrations",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
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
