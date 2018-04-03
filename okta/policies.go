package okta

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourcePolicies() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyCreate,
		Read:   resourcePolicyRead,
		Update: resourcePolicyUpdate,
		Delete: resourcePolicyDelete,

		Schema: map[string]*schema.Schema{
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
			"conditions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Conditions that must be met during Policy Evaluation",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"users": {
							Type:          schema.TypeList,
							Optional:      true,
							Description:   "List of Users to be Included or Excluded in the Policy",
							ConflictsWith: []string{"conditions.groups"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"include": {
										Type:          schema.TypeList,
										Optional:      true,
										Description:   "List of User IDs to Include",
										ConflictsWith: []string{"conditions.users.exclude"},
										Elem:          &schema.Schema{Type: schema.TypeString},
									},
									"exclude": {
										Type:          schema.TypeList,
										Optional:      true,
										Description:   "List of User IDs to Exclude",
										ConflictsWith: []string{"conditions.users.include"},
										Elem:          &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"groups": {
							Type:          schema.TypeList,
							Optional:      true,
							Description:   "List of Groups to be Included or Excluded in the Policy",
							ConflictsWith: []string{"conditions.users"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"include": {
										Type:          schema.TypeList,
										Optional:      true,
										Description:   "List of Group IDs to Include",
										ConflictsWith: []string{"conditions.groups.exclude"},
										Elem:          &schema.Schema{Type: schema.TypeString},
									},
									"exclude": {
										Type:          schema.TypeList,
										Optional:      true,
										Description:   "List of Group IDs to Exclude",
										ConflictsWith: []string{"conditions.groups.include"},
										Elem:          &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"auth_provider": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Authentication Provider for the Policy",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"provider": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "OKTA",
										Description: "Authentication Provider: OKTA or ACTIVE_DIRECTORY",
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
				},
			},
		},
	}
}

func resourcePolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Creating Policy %v", d.Get("name").(string))
	client := m.(*Config).oktaClient

	exists, _, err := policyExists(d, m)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Listing Policy in Okta: %v", err)
	}
	if exists == true {
		log.Printf("[INFO] Policy %v already exists in Okta. Adding to Terraform.", d.Get("name").(string))
	} else {
		switch d.Get("type").(string) {
		case "PASSWORD":
			template := client.Policies.PasswordPolicy()
			template.Name = d.Get("name").(string)
			template.Description = d.Get("description").(string)
			template.Type = d.Get("type").(string)
			template.Priority = d.Get("priority").(int)
			template.System = d.Get("system").(bool)
			template.Conditions.AuthProvider.Provider = "OKTA"                    // default
			template.Settings.Recovery.Factors.OktaEmail.Status = "ACTIVE"        // default
			template.Settings.Recovery.Factors.RecoveryQuestion.Status = "ACTIVE" // default

			newPolicy, _, err := client.Policies.CreatePolicy(template)
			if err != nil {
				return fmt.Errorf("[ERROR] Error Creating Policy: %v", err)
			}
			log.Printf("[INFO] Okta Policy Created: %+v", newPolicy)

			if d.Get("status").(string) == "INACTIVE" {
				_, err = client.Policies.DeactivatePolicy(newPolicy.ID)
				if err != nil {
					return fmt.Errorf("[ERROR] Error Deactivating Policy: %v", err)
				}
			}

		case "OKTA_SIGN_ON":
			template := client.Policies.SignOnPolicy()
			template.Name = d.Get("name").(string)
			template.Description = d.Get("description").(string)
			template.Type = d.Get("type").(string)
			template.Priority = d.Get("priority").(int)
			template.System = d.Get("system").(bool)

			newPolicy, _, err := client.Policies.CreatePolicy(template)
			if err != nil {
				return fmt.Errorf("[ERROR] Error Creating Policy: %v", err)
			}
			log.Printf("[INFO] Okta Policy Created: %+v", newPolicy)

			if d.Get("status").(string) == "INACTIVE" {
				_, err = client.Policies.DeactivatePolicy(newPolicy.ID)
				if err != nil {
					return fmt.Errorf("[ERROR] Error Deactivating Policy: %v", err)
				}
			}

		case "MFA_ENROLL":
			return fmt.Errorf("[ERROR] MFA Policy not supported in this terraform provider at this time")

		case "OAUTH_AUTHORIZATION_POLICY":
			return fmt.Errorf("[ERROR] Oath Auth Policy not supported in this terraform provider at this time")
		}
	}
	// add the policy resource to terraform
	d.SetId(d.Get("name").(string))

	return nil
}

func resourcePolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] List Policy %v", d.Get("name").(string))

	exists, _, err := policyExists(d, m)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Listing Policy in Okta: %v", err)
	}
	if exists == false {
		// if the policy does not exist in okta, delete from terraform state
		d.SetId("")
		return nil
	}

	return nil
}

func resourcePolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Update Policy %v", d.Get("name").(string))
	client := m.(*Config).oktaClient

	exists, id, err := policyExists(d, m)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Listing Policy in Okta: %v", err)
	}
	if exists == true {
		switch d.Get("type").(string) {
		case "PASSWORD":
			template := client.Policies.PasswordPolicy()
			template.Name = d.Get("name").(string)
			template.Description = d.Get("description").(string)
			template.Type = d.Get("type").(string)
			template.Status = d.Get("status").(string)
			template.Priority = d.Get("priority").(int)
			template.System = d.Get("system").(bool)
			template.Conditions.AuthProvider.Provider = "OKTA"                    // default
			template.Settings.Recovery.Factors.OktaEmail.Status = "ACTIVE"        // default
			template.Settings.Recovery.Factors.RecoveryQuestion.Status = "ACTIVE" // default

			updatePolicy, _, err := client.Policies.UpdatePolicy(id, template)
			if err != nil {
				return fmt.Errorf("[ERROR] Error Updating Policy: %v", err)
			}
			log.Printf("[INFO] Okta Policy Updated: %+v", updatePolicy)

			if d.Get("status").(string) == "ACTIVE" {
				_, err = client.Policies.ActivatePolicy(updatePolicy.ID)
				if err != nil {
					return fmt.Errorf("[ERROR] Error Activating Policy: %v", err)
				}
			}
			if d.Get("status").(string) == "INACTIVE" {
				_, err = client.Policies.DeactivatePolicy(updatePolicy.ID)
				if err != nil {
					return fmt.Errorf("[ERROR] Error Deactivating Policy: %v", err)
				}
			}

		case "OKTA_SIGN_ON":
			template := client.Policies.SignOnPolicy()
			template.Name = d.Get("name").(string)
			template.Description = d.Get("description").(string)
			template.Type = d.Get("type").(string)
			template.Status = d.Get("status").(string)
			template.Priority = d.Get("priority").(int)
			template.System = d.Get("system").(bool)

			updatePolicy, _, err := client.Policies.UpdatePolicy(id, template)
			if err != nil {
				return fmt.Errorf("[ERROR] Error Updating Policy: %v", err)
			}
			log.Printf("[INFO] Okta Policy Updated: %+v", updatePolicy)

			if d.Get("status").(string) == "ACTIVE" {
				_, err = client.Policies.ActivatePolicy(updatePolicy.ID)
				if err != nil {
					return fmt.Errorf("[ERROR] Error Activating Policy: %v", err)
				}
			}
			if d.Get("status").(string) == "INACTIVE" {
				_, err = client.Policies.DeactivatePolicy(updatePolicy.ID)
				if err != nil {
					return fmt.Errorf("[ERROR] Error Deactivating Policy: %v", err)
				}
			}

		case "MFA_ENROLL":
			return fmt.Errorf("[ERROR] MFA Policy not supported in this terraform provider at this time")

		case "OAUTH_AUTHORIZATION_POLICY":
			return fmt.Errorf("[ERROR] Oath Auth Policy not supported in this terraform provider at this time")
		}
	} else {
		return fmt.Errorf("[ERROR] Error Policy not found in Okta: %v", err)
	}

	return nil
}

func resourcePolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Delete Policy %v", d.Get("name").(string))
	client := m.(*Config).oktaClient

	exists, id, err := policyExists(d, m)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Listing Policy in Okta: %v", err)
	}
	if exists == true {
		_, err = client.Policies.DeletePolicy(id)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Deleting Policy from Okta: %v", err)
		}
		// delete the policy resource from terraform
		d.SetId("")
	} else {
		return fmt.Errorf("[ERROR] Error Policy not found in Okta: %v", err)
	}

	return nil
}

func policyExists(d *schema.ResourceData, m interface{}) (bool, string, error) {
	client := m.(*Config).oktaClient

	currentPolicies, _, err := client.Policies.GetPoliciesByType(d.Get("type").(string))
	if err != nil {
		return false, "", err
	}
	for _, policy := range currentPolicies.Policies {
		if policy.Name == d.Get("name").(string) {
			return true, policy.ID, nil
		}
	}
	return false, "", nil
}
