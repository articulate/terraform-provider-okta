package okta

import (
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUserSchemas() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserSchemaCreate,
		Read:   resourceUserSchemaRead,
		Update: resourceUserSchemaUpdate,
		Delete: resourceUserSchemaDelete,

		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
			if d.Get("subschema").(string) == "base" {
				return fmt.Errorf("editing a base user subschema not supported in this terraform provider at this time")
				// TODO: for the base subschema, description, format, enum, & oneof are not supported
			}
			// for an existing subschema, the subschema, index, or type fields cannot change
			prev, _ := d.GetChange("subschema")
			if prev.(string) != "" && d.HasChange("subschema") {
				return fmt.Errorf("You cannot change the subschema field for an existing User SubSchema")
			}
			prev, _ := d.GetChange("index")
			if prev.(string) != "" && d.HasChange("index") {
				return fmt.Errorf("You cannot change the index field for an existing User SubSchema")
			}
			prev, _ := d.GetChange("type")
			if prev.(string) != "" && d.HasChange("type") {
				return fmt.Errorf("You cannot change the type field for an existing User SubSchema")
			}

			return nil
		},

		Schema: map[string]*schema.Schema{
			"subschema": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"base", "custom"}, false),
				Description:  "SubSchema Type: base or custom",
			},
			"index": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subschema unique string identifier",
			},
			"title": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subschema title (display name)",
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"string", "boolean", "number", "integer", "array"}, false),
				Description:  "Subschema type: string, boolean, number, integer, or array",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom Subschema description",
			},
			"format": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom Subschema format (i.e. email)",
			},
			"required": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "whether the Subschema is required, true or false. Default = false",
			},
			"minlength": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subschema of type string minlength",
			},
			"maxlength": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subschema of type string maxlength",
			},
			"enum": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Custom Subschema enumerated value of the property. see: developer.okta.com/docs/api/resources/schemas#user-profile-schema-property-object",
			},
			"oneof": &schema.Schema{
				Type:        schema.TypeList, // this is a slice of maps - maptype?
				Optional:    true,
				Description: "Custom Subschema json schemas. see: developer.okta.com/docs/api/resources/schemas#user-profile-schema-property-object",
			},
			"permissions": &schema.Schema{
				Type:         schema.TypeList,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"HIDE", "READ_ONLY", "READ_WRITE"}, false),
				Description:  "SubSchema permissions: HIDE, READ_ONLY, or READ_WRITE. Default = READ_ONLY",
			},
			"master": &schema.Schema{
				Type:         schema.TypeList,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PROFILE_MASTER", "OKTA"}, false),
				Description:  "SubSchema profile manager: PROFILE_MASTER or OKTA. Default = PROFILE_MASTER",
			},
		},
	}
}

func resourceUserSchemaCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Creating User Schema %v", d.Get("index").(string))
	client := m.(*Config).oktaClient

	exists, err := userSchemaExists(d.Get("index").(string), d, m)
	if err != nil {
		return err
	}
	if exists == true {
		log.Printf("[INFO] User Schema %v already exists in Okta. Adding to Terraform.", d.Get("index").(string))
	} else {
		switch d.Get("subschema").(string) {
		case "base":

		case "custom":
			err = userCustomSchemaTemplate(d, m)
			if err != nil {
				return err
			}
		}
	}
	d.SetId(d.Get("index").(string))

	return nil
}

func resourceUserSchemaRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] List User Schema %v", d.Get("index").(string))
	client := m.(*Config).oktaClient

	exists, err := userSchemaExists(d.Get("index").(string), d, m)
	if err != nil {
		return err
	}
	if exists == false {
		// if the schhema does not exist in Okta, delete it from the terraform state
		d.SetId("")
		return nil
	}

	return nil
}

func resourceUserSchemaUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Update User Schema %v", d.Get("index").(string))
	client := m.(*Config).oktaClient

	exists, err := userSchemaExists(d.Get("index").(string), d, m)
	if err != nil {
		return err
	}

	d.Partial(true)
	if exists == true {
		switch d.Get("subschema").(string) {
		case "base":

		case "custom":
			err = userCustomSchemaTemplate(d, m)
			if err != nil {
				return err
			}
		}
	}
	d.Partial(false)

	return nil
}

func resourceUserSchemaDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Delete User Schema %v", d.Get("index").(string))
	client := m.(*Config).oktaClient

	exists, err := userSchemaExists(d.Get("index").(string), d, m)
	if err != nil {
		return err
	}
	if exists == true {
		switch d.Get("subschema").(string) {
		case "base":
			return fmt.Errorf("[ERROR] Error you cannot delete a base subschema")

		case "custom":
			_, _, err := client.Schemas.DeleteUserCustomSubSchema(d.Get("index").(string))
			if err != nil {
				return err
			}
		}
	}
	// remove the schema resource from terraform
	d.SetId("")

	return nil
}

func userSchemaExists(index string, d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*Config).oktaClient

	exists := false
	subschemas, _, err := client.Schemas.GetUserSubSchemaIndex(d.Get("subschema").(string))
	if err != nil {
		return fmt.Errorf("[ERROR] Error Listing User Subschemas in Okta: %v", err)
	}
	for _, key := range subschemas {
		if key == d.Get("index").(string) {
			exists = true
			break
		}
	}
	return exists, err
}

func userCustomSchemaTemplate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Config).oktaClient

	template := client.Schemas.CustomSubSchema()
	template.Index = d.Get("index").(string)
	template.Title = d.Get("title").(string)
	template.Type = d.Get("type").(string)
	if _, ok := d.GetOk("description"); ok {
		template.Descrption = d.Get("description").(string)
	}
	if _, ok := d.GetOk("format"); ok {
		template.Format = d.Get("format").(string)
	}
	if _, ok := d.GetOk("required"); ok {
		template.Required = d.Get("required").(string)
	}
	if _, ok := d.GetOk("required"); ok {
		template.Required = d.Get("required").(string)
	}
	if _, ok := d.GetOk("minlength"); ok {
		template.MinLength = d.Get("minlength").(string)
	}
	if _, ok := d.GetOk("maxlength"); ok {
		template.MaxLength = d.Get("maxlength").(string)
	}
	if _, ok := d.GetOk("master"); ok {
		template.Master = d.Get("master").(string)
	}

	// enum oneof, permissions

	_, _, err := client.Schemas.UpdateUserCustomSubSchema(template)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Creating/Updating Custom User Subschema in Okta: %v", err)
	}

	return nil
}
