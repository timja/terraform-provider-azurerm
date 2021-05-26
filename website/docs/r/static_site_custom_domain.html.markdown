---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_static_site_custom_domain"
description: |-
  Manages a Static Site Custom Domain.
---

# azurerm_static_site_custom_domain

Manages a Static Site Custom Domain.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_static_site" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

// CNAME validation
resource "azurerm_dns_cname_record" "example" {
  name                = "my-domain"
  zone_name           = "contoso.com"
  resource_group_name = azurerm_resource_group.example.name
  ttl                 = 300
  record              = azurerm_static_site.example.default_host_name
}

resource "azurerm_static_site_custom_domain" "example" {
  static_site_id  = azurerm_static_site.example.id
  name            = "${azurerm_dns_cname_record.example.name}.${azurerm_dns_cname_record.example.zone_name}"
  validation_type = "cname-delegation"
}

// TXT Validation

resource "azurerm_static_site_custom_domain" "example2" {
  static_site_id  = azurerm_static_site.example.id
  name            = "my-domain-2.${azurerm_dns_cname_record.example.zone_name}"
  validation_type = "dns-txt-token"
}

resource "azurerm_dns_txt_record" "example" {
  name                = "my-domain-2"
  zone_name           = "contoso.com"
  resource_group_name = azurerm_resource_group.example.name
  ttl                 = 300
  record {
    value = azurerm_static_site_custom_domain.example2.validation_token
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Static Site Custom Domain. Changing this forces a new Static Site Custom Domain to be created.

* `static_site_id` - (Required) The ID of the Static Site. Changing this forces a new Static Site Custom Domain to be created.

* `validation_type` - (Required) One of `cname-delegation` or `dns-txt-token`. Changing this forces a new Static Site Custom Domain to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Static Site Custom Domain.

* `validation_token` - Token to be used with `dns-txt-token` validation.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Static Site Custom Domain.
* `read` - (Defaults to 5 minutes) Used when retrieving the Static Site Custom Domain.
* `update` - (Defaults to 30 minutes) Used when updating the Static Site Custom Domain.
* `delete` - (Defaults to 30 minutes) Used when deleting the Static Site Custom Domain.

## Import

Static Site Custom Domains can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_static_site_custom_domain.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Web/staticSites/my-static-site1/customDomains/name.contoso.com
```
