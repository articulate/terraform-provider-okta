package okta

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"log"
)

type policyType struct {
	ID          string
	Description string
	Priority    int
	System      bool
}

func resourcePolicies() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyCreate,
		Read:   resourcePolicyRead,
		Update: resourcePolicyUpdate,
		Delete: resourcePolicyDelete,

		Schema: map[string]*schema.Schema{
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"OKTA_SIGN_ON", "PASSWORD", "MFA_ENROLL", "OAUTH_AUTHORIZATION_POLICY"}, false),
				Description:  "Policy Type: OKTA_SIGN_ON, PASSWORD, MFA_ENROLL, or OAUTH_AUTHORIZATION_POLICY",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Policy Name",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Policy Description",
			},
			"priority": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Policy Priority",
			},
			"status": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ACTIVE",
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "INACTIVE"}, false),
				Description:  "Policy Status: ACTIVE or INACTIVE",
			},
			"conditions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Conditions that must be met during Policy Evaluation",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"groups": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of Group IDs to Include",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"authprovider": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Authentication Provider for the Policy",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"provider": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"OKTA", "ACTIVE_DIRECTORY"}, false),
										Description:  "Authentication Provider: OKTA or ACTIVE_DIRECTORY, Active Directory currently unsupported. Okta default = OKTA",
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
			"settings": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Policy Level Settings for the Particular Policy Type",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"password": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "User Password Policies",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"minlength": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Minimum password length. Okta default = 8",
									},
									"minlowercase": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(0, 1),
										Description:  "If a password must contain at least one lower case letter: 0 = no, 1 = yes. Okta default = 1",
									},
									"minuppercase": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(0, 1),
										Description:  "If a password must contain at least one upper case letter: 0 = no, 1 = yes. Okta default = 1",
									},
									"minnumber": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(0, 1),
										Description:  "If a password must contain at least one number: 0 = no, 1 = yes. Okta default = 1",
									},
									"minsymbol": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(0, 1),
										Description:  "If a password must contain at least one symbol (!@#$%^&*): 0 = no, 1 = yes. Okta default = 1",
									},
									"excludeusername": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If the user name must be excluded from the password. Okta default = true",
									},
									"excludeattributes": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "User profile attributes that must be excluded from the password: allowed values = \"firstname\" and/or \"lastname\"",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"dictionarylookup": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Check Passwords Against Common Password Dictionary. Okta default = false",
									},
									"maxagedays": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Length in days a password is valid before expiry: 0 = no limit. Okta default = 0",
									},
									"expirewarndays": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Length in days a user will be warned before password expiry: 0 = no warning. Okta default = 0",
									},
									"minageminutes": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Minimum time interval in minutes between password changes: 0 = no limit. Okta default = 0",
									},
									"historycount": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Number of distinct passwords that can be created before they can be reused: 0 = none. Okta default = 0",
									},
									"maxlockoutattempts": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Number of unsucessful login attempts allowed before lockout: 0 = no limit. Okta default = 0",
									},
									"autounlockminutes": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Number of minutes before a locked account is unlocked: 0 = no limit. Okta default = 0",
									},
									"showlockoutfailures": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If a user should be informed when their account is locked. Okta default = false",
									},
									"recoveryquestion": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "INACTIVE"}, false),
										Description:  "Enable or Disable the recovery question: ACTIVE or INACTIVE. Okta default = ACTIVE",
									},
									"questionminlength": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Min length of the password recovery question answer. Okta default = 4",
									},
									"recoveryemailtoken": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Lifetime in minutes of the recovery email token. Okta default = 10080",
									},
									"smsrecovery": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "INACTIVE"}, false),
										Description:  "If SMS password recovery is enabled or disabled: ACTIVE or INACTIVE. Okta default = INACTIVE",
									},
									"skipunlock": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "When performing an unlock operation on an Active Directory mastered user who is locked out of Okta, the system should also attempt to unlock the user’s Windows account. Okta default = false",
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

	exists, thisPolicy, err := policyExists(d, m)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Listing Policy in Okta: %v", err)
	}
	if exists == true {
		log.Printf("[INFO] Policy %v already exists in Okta. Adding to Terraform.", d.Get("name").(string))
	} else {
		switch d.Get("type").(string) {
		case "PASSWORD":
			err = policyPassword(thisPolicy, "create", d, m)
			if err != nil {
				return err
			}

		case "OKTA_SIGN_ON":
			err = policySignOn(thisPolicy, "create", d, m)
			if err != nil {
				return err
			}

		case "MFA_ENROLL":
			return fmt.Errorf("[ERROR] MFA Policy not supported in this terraform provider at this time")

		case "OAUTH_AUTHORIZATION_POLICY":
			return fmt.Errorf("[ERROR] Oath Auth Policy not supported in this terraform provider at this time")
		}
	}
	if thisPolicy.System == true {
		log.Printf("[INFO] Policy %v is a System Policy, running Resource Policy Update.", d.Get("name").(string))
		err = resourcePolicyUpdate(d, m)
		if err != nil {
			return err
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
	d.Partial(true)

	exists, thisPolicy, err := policyExists(d, m)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Listing Policy in Okta: %v", err)
	}
	if exists == true {
		switch d.Get("type").(string) {
		case "PASSWORD":
			err = policyPassword(thisPolicy, "update", d, m)
			if err != nil {
				return err
			}

		case "OKTA_SIGN_ON":
			err = policySignOn(thisPolicy, "update", d, m)
			if err != nil {
				return err
			}

		case "MFA_ENROLL":
			return fmt.Errorf("[ERROR] MFA Policy not supported in this terraform provider at this time")

		case "OAUTH_AUTHORIZATION_POLICY":
			return fmt.Errorf("[ERROR] Oath Auth Policy not supported in this terraform provider at this time")
		}
	} else {
		return fmt.Errorf("[ERROR] Error Policy not found in Okta: %v", err)
	}
	d.Partial(false)

	return nil
}

func resourcePolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Delete Policy %v", d.Get("name").(string))
	client := m.(*Config).oktaClient

	exists, thisPolicy, err := policyExists(d, m)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Listing Policy in Okta: %v", err)
	}
	if exists == true {
		if thisPolicy.System == true {
			log.Printf("[INFO] Policy %v is a System Policy, cannot delete from Okta", d.Get("name").(string))
		} else {
			_, err = client.Policies.DeletePolicy(thisPolicy.ID)
			if err != nil {
				return fmt.Errorf("[ERROR] Error Deleting Policy from Okta: %v", err)
			}
		}
	} else {
		return fmt.Errorf("[ERROR] Error Policy not found in Okta: %v", err)
	}
	// remove the policy resource from terraform
	d.SetId("")

	return nil
}

func policyExists(d *schema.ResourceData, m interface{}) (bool, *policyType, error) {
	client := m.(*Config).oktaClient
	var thisPolicy *policyType
	thisPolicy = &policyType{System: false}

	currentPolicies, _, err := client.Policies.GetPoliciesByType(d.Get("type").(string))
	if err != nil {
		return false, thisPolicy, err
	}
	if currentPolicies != nil {
		for _, policy := range currentPolicies.Policies {
			if policy.Name == d.Get("name").(string) {
				thisPolicy = &policyType{
					ID:          policy.ID,
					Description: policy.Description,
					Priority:    policy.Priority,
					System:      policy.System,
				}
				return true, thisPolicy, nil
			}
		}
	}
	return false, thisPolicy, nil
}

func getEveryoneGroup(m interface{}) (string, error) {
	client := m.(*Config).oktaClient
	groups, _, err := client.Groups.ListGroups("q=Everyone")
	if err != nil {
		return "error", fmt.Errorf("[ERROR] ListGroups Error getting Everyone Group ID: %v", err)
	}
	if len(groups.Groups) > 1 {
		return "error", fmt.Errorf("[ERROR] Query for Everyone Default Group resulted in more than one group.")
	}
	return groups.Groups[0].ID, nil
}

func policyConditions(d *schema.ResourceData, m interface{}) ([]string, error) {
	groups := make([]string, 0)
	if len(d.Get("conditions").([]interface{})) > 0 {
		if len(d.Get("conditions.0.groups").([]interface{})) > 0 {
			for _, vals := range d.Get("conditions.0.groups").([]interface{}) {
				groups = append(groups, vals.(string))
			}
		}
		if len(d.Get("conditions.0.authprovider").([]interface{})) > 0 {
			return nil, fmt.Errorf("[ERROR] Active Directory Auth Provider not supported in this terraform provider at this time")
		}
	}
	return groups, nil
}

func policyActivate(id string, d *schema.ResourceData, m interface{}) error {
	client := m.(*Config).oktaClient

	if d.Get("status").(string) == "ACTIVE" {
		_, err := client.Policies.ActivatePolicy(id)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Activating Policy: %v", err)
		}
	}
	if d.Get("status").(string) == "INACTIVE" {
		_, err := client.Policies.DeactivatePolicy(id)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Deactivating Policy: %v", err)
		}
	}
	return nil
}

func policyPassword(thisPolicy *policyType, action string, d *schema.ResourceData, m interface{}) error {
	client := m.(*Config).oktaClient

	template := client.Policies.PasswordPolicy()
	template.Name = d.Get("name").(string)
	template.Type = d.Get("type").(string)
	if thisPolicy.System == true {
		template.Status = "ACTIVE"
		template.Description = thisPolicy.Description
		template.Priority = thisPolicy.Priority

		everyone, err := getEveryoneGroup(m)
		if err != nil {
			return err
		}
		template.Conditions.People.Groups.Include = []string{everyone}

	} else {
		template.Description = d.Get("description").(string)
		template.Priority = d.Get("priority").(int)
		template.Status = d.Get("status").(string)

		groups, err := policyConditions(d, m)
		if err != nil {
			return err
		}
		template.Conditions.People.Groups.Include = groups
	}

	template.Settings.Recovery.Factors.RecoveryQuestion.Status = "ACTIVE" // okta required default
	template.Settings.Recovery.Factors.OktaEmail.Status = "ACTIVE"        // okta required & read-only default

	if len(d.Get("settings.0.password").([]interface{})) > 0 {
		template.Settings.Password.Complexity.MinLength = d.Get("settings.0.password.0.minlength").(int)
		template.Settings.Password.Complexity.MinLowerCase = d.Get("settings.0.password.0.minlowercase").(int)
		template.Settings.Password.Complexity.MinUpperCase = d.Get("settings.0.password.0.minuppercase").(int)
		template.Settings.Password.Complexity.MinNumber = d.Get("settings.0.password.0.minnumber").(int)
		template.Settings.Password.Complexity.MinSymbol = d.Get("settings.0.password.0.minsymbol").(int)
		template.Settings.Password.Complexity.ExcludeUsername = d.Get("settings.0.password.0.excludeusername").(bool)
		if len(d.Get("settings.0.password.0.excludeattributes").([]interface{})) > 0 {
			exclude := make([]string, 0)
			for _, vals := range d.Get("settings.0.password.0.excludeattributes").([]interface{}) {
				exclude = append(exclude, vals.(string))
			}
			template.Settings.Password.Complexity.ExcludeAttributes = exclude
		}
		template.Settings.Password.Complexity.Dictionary.Common.Exclude = d.Get("settings.0.password.0.dictionarylookup").(bool)
		template.Settings.Password.Age.MaxAgeDays = d.Get("settings.0.password.0.maxagedays").(int)
		template.Settings.Password.Age.ExpireWarnDays = d.Get("settings.0.password.0.expirewarndays").(int)
		template.Settings.Password.Age.MinAgeMinutes = d.Get("settings.0.password.0.minageminutes").(int)
		template.Settings.Password.Age.HistoryCount = d.Get("settings.0.password.0.historycount").(int)
		template.Settings.Password.Lockout.MaxAttempts = d.Get("settings.0.password.0.maxlockoutattempts").(int)
		template.Settings.Password.Lockout.AutoUnlockMinutes = d.Get("settings.0.password.0.autounlockminutes").(int)
		template.Settings.Password.Lockout.ShowLockoutFailures = d.Get("settings.0.password.0.showlockoutfailures").(bool)
		if d.Get("settings.0.password.0.recoveryquestion").(string) != "" {
			template.Settings.Recovery.Factors.RecoveryQuestion.Status = d.Get("settings.0.password.0.recoveryquestion").(string)
		}
		template.Settings.Recovery.Factors.RecoveryQuestion.Properties.Complexity.MinLength = d.Get("settings.0.password.0.questionminlength").(int)
		template.Settings.Recovery.Factors.OktaEmail.Properties.RecoveryToken.TokenLifetimeMinutes = d.Get("settings.0.password.0.recoveryemailtoken").(int)
		template.Settings.Recovery.Factors.OktaSms.Status = d.Get("settings.0.password.0.smsrecovery").(string)
		template.Settings.Delegation.Options.SkipUnlock = d.Get("settings.0.password.0.skipunlock").(bool)
	}

	switch action {
	case "create":
		policy, _, err := client.Policies.CreatePolicy(template)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Creating Policy: %v", err)
		}
		log.Printf("[INFO] Okta Policy Created: %+v", policy)

		err = policyActivate(policy.ID, d, m)
		if err != nil {
			return err
		}

	case "update":
		policy, _, err := client.Policies.UpdatePolicy(thisPolicy.ID, template)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Updating Policy: %v", err)
		}
		log.Printf("[INFO] Okta Policy Updated: %+v", policy)

		if thisPolicy.System == false {
			err = policyActivate(policy.ID, d, m)
			if err != nil {
				return err
			}
		}

	default:
		return fmt.Errorf("[ERROR] policyPassword action only supports \"create\" and \"update\"")
	}
	return nil
}

func policySignOn(thisPolicy *policyType, action string, d *schema.ResourceData, m interface{}) error {
	client := m.(*Config).oktaClient

	template := client.Policies.SignOnPolicy()
	template.Name = d.Get("name").(string)
	template.Type = d.Get("type").(string)
	if thisPolicy.System == true {
		template.Status = "ACTIVE"
		template.Description = thisPolicy.Description
		template.Priority = thisPolicy.Priority

		everyone, err := getEveryoneGroup(m)
		if err != nil {
			return err
		}
		template.Conditions.People.Groups.Include = []string{everyone}

	} else {
		template.Description = d.Get("description").(string)
		template.Priority = d.Get("priority").(int)
		template.Status = d.Get("status").(string)

		groups, err := policyConditions(d, m)
		if err != nil {
			return err
		}
		template.Conditions.People.Groups.Include = groups
	}

	switch action {
	case "create":
		policy, _, err := client.Policies.CreatePolicy(template)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Creating Policy: %v", err)
		}
		log.Printf("[INFO] Okta Policy Created: %+v", policy)

		err = policyActivate(policy.ID, d, m)
		if err != nil {
			return err
		}

	case "update":
		policy, _, err := client.Policies.UpdatePolicy(thisPolicy.ID, template)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Updating Policy: %v", err)
		}
		log.Printf("[INFO] Okta Policy Updated: %+v", policy)

		if thisPolicy.System == false {
			err = policyActivate(policy.ID, d, m)
			if err != nil {
				return err
			}
		}

	default:
		return fmt.Errorf("[ERROR] policySignOn action only supports \"create\" and \"update\"")
	}
	return nil
}
