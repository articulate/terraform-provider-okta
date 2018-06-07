package okta

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceUserSchemas() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserSchemaCreate,
		Read:   resourceUserSchemaRead,
		Update: resourceUserSchemaUpdate,
		Delete: resourceUserSchemaDelete,

		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {

			// for an existing subschema, the subschema, index, or type fields cannot change
			prev, _ := d.GetChange("subschema")
			if prev.(string) != "" && d.HasChange("subschema") {
				return fmt.Errorf("You cannot change the subschema field for an existing User SubSchema")
			}
			prev, _ = d.GetChange("index")
			if prev.(string) != "" && d.HasChange("index") {
				return fmt.Errorf("You cannot change the index field for an existing User SubSchema")
			}
			prev, _ = d.GetChange("type")
			if prev.(string) != "" && d.HasChange("type") {
				return fmt.Errorf("You cannot change the type field for an existing User SubSchema")
			}

			// arraytype field only required if type field is array
			if _, ok := d.GetOk("arraytype"); ok {
				if d.Get("type").(string) != "array" {
					return fmt.Errorf("arraytype field not required if type field is not array")
				}
			} else {
				if d.Get("type").(string) == "array" {
					return fmt.Errorf("arraytype field required if type field is array")
				}
			}

			// error out in the terraform plan stage if user adds to config options not supported yet in this provider
			if d.Get("subschema").(string) == "base" {
				return fmt.Errorf("Editing a base user SubSchema not supported in this terraform provider at this time")
				// todo: for the base subschema, description, enum, & oneof are not supported
			}
			switch d.Get("type").(string) {
			case "boolean":
				return fmt.Errorf("Editing a custom SubSchema of type boolean not supported in this terraform provider at this time")

			case "number":
				return fmt.Errorf("Editing a custom SubSchema of type number not supported in this terraform provider at this time")

			case "interger":
				return fmt.Errorf("Editing a custom SubSchema of type interger not supported in this terraform provider at this time")

			case "array":
				if d.Get("arraytype").(string) != "string" {
					return fmt.Errorf("Editing a custom SubSchema of type array (number, interger, or reference) not supported in this terraform provider at this time")
				}
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
			"arraytype": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"string", "number", "interger", "reference"}, false),
				Description:  "Subschema array type: string, number, interger, reference. Type field must be an array.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom Subschema description",
			},
			"required": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "whether the Subschema is required, true or false. Default = false",
			},
			"minlength": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Subschema of type string minlength",
			},
			"maxlength": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Subschema of type string maxlength",
			},
			"enum": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Custom Subschema enumerated value of the property. see: developer.okta.com/docs/api/resources/schemas#user-profile-schema-property-object",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"oneof": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Custom Subschema json schemas. see: developer.okta.com/docs/api/resources/schemas#user-profile-schema-property-object",
				//Elem:        schema.TypeString,
			},
			"permissions": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"HIDE", "READ_ONLY", "READ_WRITE"}, false),
				Description:  "SubSchema permissions: HIDE, READ_ONLY, or READ_WRITE. Default = READ_ONLY",
			},
			"master": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PROFILE_MASTER", "OKTA"}, false),
				Description:  "SubSchema profile manager: PROFILE_MASTER or OKTA. Default = PROFILE_MASTER",
			},
		},
	}
}

func resourceUserSchemaCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Creating User Schema %v", d.Get("index").(string))

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
		return exists, fmt.Errorf("[ERROR] Error Listing User Subschemas in Okta: %v", err)
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

	perms := client.Schemas.Permissions()
	perms.Principal = "SELF"
	perms.Action = "READ_ONLY"

	template := client.Schemas.CustomSubSchema()
	template.Index = d.Get("index").(string)
	template.Title = d.Get("title").(string)
	template.Type = d.Get("type").(string)
	template.Master.Type = "PROFILE_MASTER"
	if _, ok := d.GetOk("arraytype"); ok {
		template.Items.Type = d.Get("arraytype").(string)
	}
	if _, ok := d.GetOk("description"); ok {
		template.Description = d.Get("description").(string)
	}
	if _, ok := d.GetOk("required"); ok {
		template.Required = d.Get("required").(bool)
	}
	if _, ok := d.GetOk("minlength"); ok {
		template.MinLength = d.Get("minlength").(int)
	}
	if _, ok := d.GetOk("maxlength"); ok {
		template.MaxLength = d.Get("maxlength").(int)
	}
	if _, ok := d.GetOk("master"); ok {
		template.Master.Type = d.Get("master").(string)
	}
	if _, ok := d.GetOk("permissions"); ok {
		perms.Action = d.Get("permissions").(string)
	}
	template.Permissions = append(template.Permissions, perms)
	if _, ok := d.GetOk("enum"); ok {
		enum := userEnumSchema(d)
		template.Enum = enum
	}

	return fmt.Errorf("%+v", d)

	// oneof

	_, _, err := client.Schemas.UpdateUserCustomSubSchema(template)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Creating/Updating Custom User Subschema in Okta: %v", err)
	}

	return nil
}

func userEnumSchema(d *schema.ResourceData) []string {
	enum := make([]string, 0)
	if len(d.Get("enum").([]interface{})) > 0 {
		for _, vals := range d.Get("enum").([]interface{}) {
			enum = append(enum, vals.(string))
		}
	}
	return enum
}
