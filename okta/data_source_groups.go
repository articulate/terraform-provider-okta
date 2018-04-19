// this data source is a temporary "fix" for the policy tests to lookup the Everyone group ID
// this data source needs to be deleted after the groups resource is added to the provider

package okta

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Policy name",
				Required:    true,
			},
		},
	}
}

func dataSourceGroupRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Data Source Group Read %v", d.Get("name").(string))
	client := m.(*Config).oktaClient

	query := fmt.Sprintf("q=%v", d.Get("name").(string))
	groups, _, err := client.Groups.ListGroups(query)
	if err != nil {
		return fmt.Errorf("[ERROR] ListGroups query error: %v", err)
	}
	if groups != nil {
		if len(groups.Groups) > 1 {
			return fmt.Errorf("[ERROR] Group query resulted in more than one group.")
		}
		d.SetId(groups.Groups[0].ID)
	} else {
		return fmt.Errorf("[ERROR] Group query resulted in no groups.")
	}
	return nil
}
