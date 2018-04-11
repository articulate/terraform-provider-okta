package okta

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"log"
)

// fields retrieved from the policy in Okta & referenced in our resource functions
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

		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
			// user cannot change name or type for an existing policy
			prev, _ := d.GetChange("name")
			if prev.(string) != "" {
				if d.HasChange("type") || d.HasChange("name") {
					return fmt.Errorf("You cannot change the name field or type field of an existing Policy")
				}
			}

			// add custom error messages if user supplies options not supported by the policy
			switch d.Get("type").(string) {
			case "PASSWORD":

			case "OKTA_SIGN_ON":
				if len(d.Get("conditions.0.authprovider").([]interface{})) > 0 {
					return fmt.Errorf("authprovider condition options not supported in the Okta SignOn Policy")
				}
				if len(d.Get("settings.0.password").([]interface{})) > 0 {
					return fmt.Errorf("password settings options not supported in the Okta SignOn Policy")
				}

			case "MFA_ENROLL":

			case "OAUTH_AUTHORIZATION_POLICY":

			}

			// password settings option excludeattributes only supports "firstName" and/or "lastName"
			// ValidateFunc currently not supported in terraform for list types so we'll add our check here
			if d.HasChange("settings.0.password.0.excludeattributes") {
				for _, vals := range d.Get("settings.0.password.0.excludeattributes").([]interface{}) {
					if vals.(string) != "firstName" || vals.(string) != "lastName" {
						return fmt.Errorf("accepted values for excludeattributes password settings are \"firstName\" and/or \"lastName\"")
					}
				}
			}

			return nil
		},

		Schema: map[string]*schema.Schema{
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"OKTA_SIGN_ON", "PASSWORD", "MFA_ENROLL", "OAUTH_AUTHORIZATION_POLICY"}, false),
				Description:  "Policy Type: OKTA_SIGN_ON, PASSWORD, MFA_ENROLL, or OAUTH_AUTHORIZATION_POLICY",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
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
										Description:  "Authentication Provider: OKTA or ACTIVE_DIRECTORY. Default = OKTA",
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
										Description: "Minimum password length. Default = 8",
									},
									"minlowercase": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(0, 1),
										Description:  "If a password must contain at least one lower case letter: 0 = no, 1 = yes. Default = 1",
									},
									"minuppercase": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(0, 1),
										Description:  "If a password must contain at least one upper case letter: 0 = no, 1 = yes. Default = 1",
									},
									"minnumber": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(0, 1),
										Description:  "If a password must contain at least one number: 0 = no, 1 = yes. Default = 1",
									},
									"minsymbol": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(0, 1),
										Description:  "If a password must contain at least one symbol (!@#$%^&*): 0 = no, 1 = yes. Default = 1",
									},
									"excludeusername": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If the user name must be excluded from the password. Default = true",
									},
									"excludeattributes": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "User profile attributes that must be excluded from the password: allowed values = \"firstName\" and/or \"lastName\"",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"dictionarylookup": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Check Passwords Against Common Password Dictionary. Default = false",
									},
									"maxagedays": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Length in days a password is valid before expiry: 0 = no limit. Default = 0",
									},
									"expirewarndays": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Length in days a user will be warned before password expiry: 0 = no warning. Default = 0",
									},
									"minageminutes": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Minimum time interval in minutes between password changes: 0 = no limit. Default = 0",
									},
									"historycount": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Number of distinct passwords that can be created before they can be reused: 0 = none. Default = 0",
									},
									"maxlockoutattempts": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Number of unsucessful login attempts allowed before lockout: 0 = no limit. Default = 0",
									},
									"autounlockminutes": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Number of minutes before a locked account is unlocked: 0 = no limit. Default = 0",
									},
									"showlockoutfailures": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If a user should be informed when their account is locked. Default = false",
									},
									"recoveryquestion": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "INACTIVE"}, false),
										Description:  "Enable or Disable the recovery question: ACTIVE or INACTIVE. Default = ACTIVE",
									},
									"questionminlength": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Min length of the password recovery question answer. Default = 4",
									},
									"recoveryemailtoken": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Lifetime in minutes of the recovery email token. Default = 10080",
									},
									"smsrecovery": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "INACTIVE"}, false),
										Description:  "If SMS password recovery is enabled or disabled: ACTIVE or INACTIVE. Default = INACTIVE",
									},
									"skipunlock": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "When an Active Directory user is locked out of Okta, the Okta unlock operation should also attempt to unlock the user's Windows account. Default = false",
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
		return err
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
		return err
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
		return err
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
		return err
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

// check if policy exists in Okta & return a struct of policy fields we'll reference during apply
func policyExists(d *schema.ResourceData, m interface{}) (bool, *policyType, error) {
	client := m.(*Config).oktaClient
	var thisPolicy *policyType
	thisPolicy = &policyType{System: false}

	currentPolicies, _, err := client.Policies.GetPoliciesByType(d.Get("type").(string))
	if err != nil {
		return false, thisPolicy, fmt.Errorf("[ERROR] Error Listing Policy in Okta: %v", err)
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

// system default policies use the Everyone group in the groups condition include, get that group ID
func getEveryoneGroup(m interface{}) (string, error) {
	client := m.(*Config).oktaClient
	groups, _, err := client.Groups.ListGroups("q=Everyone")
	if err != nil {
		return "error", fmt.Errorf("[ERROR] ListGroups Error querying Everyone Group ID: %v", err)
	}
	if len(groups.Groups) > 1 {
		return "error", fmt.Errorf("[ERROR] Query for Everyone Default Group resulted in more than one group.")
	}
	return groups.Groups[0].ID, nil
}

// populate policy conditions with the terraform schema conditions fields
func policyConditions(d *schema.ResourceData) ([]string, error) {
	groups := make([]string, 0)
	if len(d.Get("conditions").([]interface{})) > 0 {
		if len(d.Get("conditions.0.groups").([]interface{})) > 0 {
			for _, vals := range d.Get("conditions.0.groups").([]interface{}) {
				groups = append(groups, vals.(string))
			}
		}
		if len(d.Get("conditions.0.authprovider").([]interface{})) > 0 {
			return groups, fmt.Errorf("[ERROR] Active Directory Auth Provider not supported in this terraform provider at this time")
		}
	}
	return groups, nil
}

// activate or deactivate a policy according to the terraform schema status field
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

// create or update a password policy
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

		groups, err := policyConditions(d)
		if err != nil {
			return err
		}
		template.Conditions.People.Groups.Include = groups
	}

	// if our password settings schema fields are undefined, use the Okta defaults
	// we add the defaults here & not in the schema map to avoid defaults appearing in the terraform plan diff
	template.Settings.Recovery.Factors.OktaEmail.Status = "ACTIVE" // okta required & read-only default
	if len(d.Get("settings.0.password").([]interface{})) > 0 {
		if minlength, ok := d.GetOk("settings.0.password.0.minlength"); ok {
			template.Settings.Password.Complexity.MinLength = minlength.(int)
		} else {
			template.Settings.Password.Complexity.MinLength = 8
		}
		if minlowercase, ok := d.GetOk("settings.0.password.0.minlowercase"); ok {
			template.Settings.Password.Complexity.MinLowerCase = minlowercase.(int)
		} else {
			template.Settings.Password.Complexity.MinLowerCase = 1
		}
		if minuppercase, ok := d.GetOk("settings.0.password.0.minuppercase"); ok {
			template.Settings.Password.Complexity.MinUpperCase = minuppercase.(int)
		} else {
			template.Settings.Password.Complexity.MinUpperCase = 1
		}
		if minnumber, ok := d.GetOk("settings.0.password.0.minnumber"); ok {
			template.Settings.Password.Complexity.MinNumber = minnumber.(int)
		} else {
			template.Settings.Password.Complexity.MinNumber = 1
		}
		if minsymbol, ok := d.GetOk("settings.0.password.0.minsymbol"); ok {
			template.Settings.Password.Complexity.MinSymbol = minsymbol.(int)
		} else {
			template.Settings.Password.Complexity.MinSymbol = 1
		}
		if excludeusername, ok := d.GetOk("settings.0.password.0.excludeusername"); ok {
			template.Settings.Password.Complexity.ExcludeUsername = excludeusername.(bool)
		} else {
			template.Settings.Password.Complexity.ExcludeUsername = true
		}
		if len(d.Get("settings.0.password.0.excludeattributes").([]interface{})) > 0 {
			exclude := make([]string, 0)
			for _, vals := range d.Get("settings.0.password.0.excludeattributes").([]interface{}) {
				exclude = append(exclude, vals.(string))
			}
			template.Settings.Password.Complexity.ExcludeAttributes = exclude
		}
		if dictionarylookup, ok := d.GetOk("settings.0.password.0.dictionarylookup"); ok {
			template.Settings.Password.Complexity.Dictionary.Common.Exclude = dictionarylookup.(bool)
		} else {
			template.Settings.Password.Complexity.Dictionary.Common.Exclude = false
		}
		if maxagedays, ok := d.GetOk("settings.0.password.0.maxagedays"); ok {
			template.Settings.Password.Age.MaxAgeDays = maxagedays.(int)
		} else {
			template.Settings.Password.Age.MaxAgeDays = 0
		}
		if expirewarndays, ok := d.GetOk("settings.0.password.0.expirewarndays"); ok {
			template.Settings.Password.Age.ExpireWarnDays = expirewarndays.(int)
		} else {
			template.Settings.Password.Age.ExpireWarnDays = 0
		}
		if minageminutes, ok := d.GetOk("settings.0.password.0.minageminutes"); ok {
			template.Settings.Password.Age.MinAgeMinutes = minageminutes.(int)
		} else {
			template.Settings.Password.Age.MinAgeMinutes = 0
		}
		if historycount, ok := d.GetOk("settings.0.password.0.historycount"); ok {
			template.Settings.Password.Age.HistoryCount = historycount.(int)
		} else {
			template.Settings.Password.Age.HistoryCount = 0
		}
		if maxlockoutattempts, ok := d.GetOk("settings.0.password.0.maxlockoutattempts"); ok {
			template.Settings.Password.Lockout.MaxAttempts = maxlockoutattempts.(int)
		} else {
			template.Settings.Password.Lockout.MaxAttempts = 0
		}
		if autounlockminutes, ok := d.GetOk("settings.0.password.0.autounlockminutes"); ok {
			template.Settings.Password.Lockout.AutoUnlockMinutes = autounlockminutes.(int)
		} else {
			template.Settings.Password.Lockout.AutoUnlockMinutes = 0
		}
		if showlockoutfailures, ok := d.GetOk("settings.0.password.0.showlockoutfailures"); ok {
			template.Settings.Password.Lockout.ShowLockoutFailures = showlockoutfailures.(bool)
		} else {
			template.Settings.Password.Lockout.ShowLockoutFailures = false
		}
		if recoveryquestion, ok := d.GetOk("settings.0.password.0.recoveryquestion"); ok {
			template.Settings.Recovery.Factors.RecoveryQuestion.Status = recoveryquestion.(string)
		} else {
			template.Settings.Recovery.Factors.RecoveryQuestion.Status = "ACTIVE"
		}
		if questionminlength, ok := d.GetOk("settings.0.password.0.questionminlength"); ok {
			template.Settings.Recovery.Factors.RecoveryQuestion.Properties.Complexity.MinLength = questionminlength.(int)
		} else {
			template.Settings.Recovery.Factors.RecoveryQuestion.Properties.Complexity.MinLength = 4
		}
		if recoveryemailtoken, ok := d.GetOk("settings.0.password.0.recoveryemailtoken"); ok {
			template.Settings.Recovery.Factors.OktaEmail.Properties.RecoveryToken.TokenLifetimeMinutes = recoveryemailtoken.(int)
		} else {
			template.Settings.Recovery.Factors.OktaEmail.Properties.RecoveryToken.TokenLifetimeMinutes = 10080
		}
		if smsrecovery, ok := d.GetOk("settings.0.password.0.smsrecovery"); ok {
			template.Settings.Recovery.Factors.OktaSms.Status = smsrecovery.(string)
		} else {
			template.Settings.Recovery.Factors.OktaSms.Status = "INACTIVE"
		}
		if skipunlock, ok := d.GetOk("settings.0.password.0.skipunlock"); ok {
			template.Settings.Delegation.Options.SkipUnlock = skipunlock.(bool)
		} else {
			template.Settings.Delegation.Options.SkipUnlock = false
		}
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

// create or update a signon policy
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

		groups, err := policyConditions(d)
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
