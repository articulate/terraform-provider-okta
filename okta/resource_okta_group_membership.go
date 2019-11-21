package okta

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGroupMembership() *schema.Resource {
	return &schema.Resource{

		Create: resourceGroupMembershipCreate,
		Read:   resourceGroupMembershipRead,
		Delete: resourceGroupMembershipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"user_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Use to associate group with",
			},
			"group_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Group associated with the user",
			},
			"count": &schema.Schema{
				Type:        schema.TypeInt,
				Default:     1,
				Optional:    true,
				Description: "Allow looping through a list of users to assign groups",
			},
		},
	}
}

func resourceGroupMembershipCreate(d *schema.ResourceData, m interface{}) error {

	client := m.(*Config).oktaClient

	if d.Get("group_id").(string) != "" {

		if err := assignGroupsToUser(d.Get("user_id").(string), d.Get("group_id").([]string), client); err != nil {

			return err

		}
	}

	d.SetId(d.Get("group_id").(string))

	return resourceGroupMembershipRead(d, m)
}

func resourceGroupMembershipRead(d *schema.ResourceData, m interface{}) error {

	userIdList, err := listGroupUserIds(m, d.Id())
	if err != nil {
		return err
	}

	return d.Set("users", convertStringSetToInterface(userIdList))
}

func resourceGroupMembershipDelete(d *schema.ResourceData, m interface{}) error {

	userIdList, err := listGroupUserIds(m, d.Id())
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
