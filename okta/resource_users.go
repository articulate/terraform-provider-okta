package okta

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUsers() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,

		Schema: map[string]*schema.Schema{
			"firstname": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "User first name",
			},
			"lastname": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "User last name",
			},
			"email": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "User email address",
			},
			"login": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User Okta login",
			},
			"role": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User Okta role",
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Creating User %v", d.Get("email").(string))
	client := m.(*Config).oktaClient

	filter := client.Users.UserListFilterOptions()
	filter.EmailEqualTo = d.Get("email").(string)
	newUser, _, err := client.Users.ListWithFilter(&filter)
	if len(newUser) == 0 {
		userTemplate("create", d, m)
		if err != nil {
			return err
		}
		return nil
	}
	if len(newUser) > 1 {
		return fmt.Errorf("[ERROR] Retrieved more than one Okta user for the email %v", d.Get("email").(string))
	}
	log.Printf("[INFO] User already exists in Okta. Adding to Terraform")
	// add the user resource to terraform
	d.SetId(newUser[0].ID)

	return nil
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] List User %v", d.Get("email").(string))
	client := m.(*Config).oktaClient

	_, _, err := client.Users.GetByID(d.Id())
	if err != nil {
		// if the user does not exist in okta, delete from terraform state
		if client.OktaErrorCode == "E0000007" {
			d.SetId("")
			return nil
		} else {
			return fmt.Errorf("[ERROR] Error GetByID: %v", err)
		}
	} else {
		userRoles, _, err := client.Users.ListRoles(d.Id())
		if err != nil {
			return fmt.Errorf("[ERROR] Error listing user role: %v", err)
		}
		if userRoles != nil {
			if len(userRoles.Role) > 1 {
				return fmt.Errorf("[ERROR] User has more than one role. This terraform provider presently only supports a single role per user. Please review the user's role assignments in Okta.")
			}
		}
	}

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Update User %v", d.Get("email").(string))
	client := m.(*Config).oktaClient
	d.Partial(true)

	_, _, err := client.Users.GetByID(d.Id())
	if err != nil {
		return fmt.Errorf("[ERROR] Error GetByID: %v", err)
	}
	userTemplate("update", d, m)
	if err != nil {
		return err
	}
	d.Partial(false)

	return nil
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Delete User %v", d.Get("email").(string))
	client := m.(*Config).oktaClient

	userList, _, err := client.Users.GetByID(d.Id())
	if err != nil {
		return fmt.Errorf("[ERROR] Error GetByID: %v", err)
	}
	// must deactivate the user before deletion
	if userList.Status != "DEPROVISIONED" {
		_, err := client.Users.Deactivate(d.Id())
		if err != nil {
			return fmt.Errorf("[ERROR] Error Deactivating user: %v", err)
		}
	}
	// delete the user
	_, err = client.Users.Delete(d.Id())
	if err != nil {
		return fmt.Errorf("[ERROR] Error Deleting user: %v", err)
	}
	// delete the user resource from terraform
	d.SetId("")

	return nil
}

func userTemplate(action string, d *schema.ResourceData, m interface{}) error {
	client := m.(*Config).oktaClient

	template := client.Users.NewUser()
	template.Profile.FirstName = d.Get("firstname").(string)
	template.Profile.LastName = d.Get("lastname").(string)
	template.Profile.Email = d.Get("email").(string)
	if _, ok := d.GetOk("login"); ok {
		template.Profile.Login = d.Get("login").(string)
	} else {
		template.Profile.Login = d.Get("email").(string)
	}

	switch action {
	case "create":
		// activate user but send an email to set their password
		// okta user status will be "Password Reset" until they complete
		// the okta signup process
		createNewUserAsActive := true

		newUser, _, err := client.Users.Create(template, createNewUserAsActive)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Creating User: %v", err)
		}
		log.Printf("[INFO] Okta User Created: %+v", newUser)

		// assign the user a role, if specified
		if _, ok := d.GetOk("role"); ok {
			log.Printf("[INFO] Assigning role: " + d.Get("role").(string))
			_, err := client.Users.AssignRole(newUser.ID, d.Get("role").(string))
			if err != nil {
				return fmt.Errorf("[ERROR] Error assigning role to user: %v", err)
			}
		}
		// add the user resource to terraform
		d.SetId(newUser.ID)

	case "update":
		updateUser, _, err := client.Users.Update(template, d.Id())
		if err != nil {
			return fmt.Errorf("[ERROR] Error Updating User: %v", err)
		}
		log.Printf("[INFO] Okta User Updated: %+v", updateUser)

		userRoles, _, err := client.Users.ListRoles(d.Id())
		if err != nil {
			return fmt.Errorf("[ERROR] Error listing user role: %v", err)
		}

		if d.HasChange("role") {
			if userRoles != nil {
				log.Printf("[INFO] Removing role: " + userRoles.Role[0].Type)
				_, err = client.Users.UnAssignRole(d.Id(), userRoles.Role[0].ID)
				if err != nil {
					return fmt.Errorf("[ERROR] Error removing role from user: %v", err)
				}
			}
			if _, ok := d.GetOk("role"); ok {
				log.Printf("[INFO] Assigning role: " + d.Get("role").(string))
				_, err = client.Users.AssignRole(d.Id(), d.Get("role").(string))
				if err != nil {
					return fmt.Errorf("[ERROR] Error assigning role to user: %v", err)
				}
			}
		}

	default:
		return fmt.Errorf("[ERROR] userTemplate action only supports \"create\" and \"update\"")
	}

	return nil
}
