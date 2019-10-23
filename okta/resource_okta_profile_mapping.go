package okta

import (
	"github.com/articulate/terraform-provider-okta/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var (
	sourceSchema = &schema.Schema{
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
				"type": &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
				"name": &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}

	mappingSchema = &schema.Schema{
		Type: schema.TypeSet,
		Elem: schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Description: "The mapping property key.",
				},
				"expression": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"push_status": &schema.Schema{
					Type:         schema.TypeString,
					Optional:     true,
					Default:      dontPush,
					ValidateFunc: validation.StringInSlice([]string{push, dontPush}, false),
				},
			},
		},
	}
)

const (
	push     = "PUSH"
	dontPush = "DONT_PUSH"
)

func resourceOktaProfileMapping() *schema.Resource {
	return &schema.Resource{
		Create: resourceProfileMappingCreate,
		Read:   resourceProfileMappingRead,
		Update: resourceProfileMappingUpdate,
		Delete: resourceProfileMappingDelete,
		Exists: resourceProfileMappingExists,

		Schema: map[string]*schema.Schema{
			"source_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The source id of the mapping to manage.",
			},
			"delete_when_absent": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When turned on this flag will trigger the provider to delete mapping properties that are not defined in config. By default, we do not delete missing properties.",
			},
			"source":   sourceSchema,
			"target":   sourceSchema,
			"mappings": mappingSchema,
		},
	}
}

func getProfileMapping() (*sdk.ProfileMapping, error) {
	return nil, nil
}

func resourceProfileMappingCreate(d *schema.ResourceData, m interface{}) error {
	return resourceProfileMappingRead(d, m)
}

func resourceProfileMappingDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceProfileMappingExists(d *schema.ResourceData, m interface{}) (bool, error) {
	return false, nil
}

func resourceProfileMappingRead(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	sourceId := d.Get("source_id").(string)
	mapping, _, err := client.GetProfileBySourceId(sourceId)

	if err != nil {
		return err
	}

	return nil
}

func resourceProfileMappingUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceProfileMappingRead(d, m)
}
