package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-12-01/postgresql"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPostgreSQLConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPostgreSQLConfigurationCreateUpdate,
		Read:   resourceArmPostgreSQLConfigurationRead,
		Delete: resourceArmPostgreSQLConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"server_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"value": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceArmPostgreSQLConfigurationCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).postgresqlConfigurationsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM PostgreSQL Configuration creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	serverName := d.Get("server_name").(string)

	value := d.Get("value").(string)

	properties := postgresql.Configuration{
		ConfigurationProperties: &postgresql.ConfigurationProperties{
			Value: utils.String(value),
		},
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, serverName, name, properties)
	if err != nil {
		return err
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resourceGroup, serverName, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read PostgreSQL Configuration %s (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmPostgreSQLConfigurationRead(d, meta)
}

func resourceArmPostgreSQLConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).postgresqlConfigurationsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	name := id.Path["configurations"]

	resp, err := client.Get(ctx, resourceGroup, serverName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] PostgreSQL Configuration '%s' was not found (resource group '%s')", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure PostgreSQL Configuration %s: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("server_name", serverName)
	d.Set("resource_group_name", resourceGroup)
	d.Set("value", resp.ConfigurationProperties.Value)

	return nil
}

func resourceArmPostgreSQLConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).postgresqlConfigurationsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	name := id.Path["configurations"]

	// "delete" = resetting this to the default value
	resp, err := client.Get(ctx, resourceGroup, serverName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Postgresql Configuration '%s': %+v", name, err)
	}

	properties := postgresql.Configuration{
		ConfigurationProperties: &postgresql.ConfigurationProperties{
			// we can alternatively set `source: "system-default"`
			Value: resp.DefaultValue,
		},
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, serverName, name, properties)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return err
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return err
	}

	return nil
}
