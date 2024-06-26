---
page_title: "snowflake_scim_integration Resource - terraform-provider-snowflake"
subcategory: ""
description: |-
  
---

# snowflake_scim_integration (Resource)



## Example Usage

```terraform
# basic resource
resource "snowflake_scim_integration" "test" {
  name          = "test"
  enabled       = true
  scim_client   = "GENERIC"
  sync_password = true
}
# resource with all fields set
resource "snowflake_scim_integration" "test" {
  name           = "test"
  enabled        = true
  scim_client    = "GENERIC"
  sync_password  = true
  network_policy = "network_policy_test"
  run_as_role    = "GENERIC_SCIM_PROVISIONER"
  comment        = "foo"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `enabled` (Boolean) Specify whether the security integration is enabled.
- `name` (String) String that specifies the identifier (i.e. name) for the integration; must be unique in your account.
- `run_as_role` (String) Specify the SCIM role in Snowflake that owns any users and roles that are imported from the identity provider into Snowflake using SCIM. Provider assumes that the specified role is already provided. Valid options are: [OKTA_PROVISIONER AAD_PROVISIONER GENERIC_SCIM_PROVISIONER].
- `scim_client` (String) Specifies the client type for the scim integration. Valid options are: [OKTA AZURE GENERIC].

### Optional

- `comment` (String) Specifies a comment for the integration.
- `network_policy` (String) Specifies an existing network policy that controls SCIM network traffic.
- `sync_password` (Boolean) Specifies whether to enable or disable the synchronization of a user password from an Okta SCIM client as part of the API request to Snowflake.

### Read-Only

- `created_on` (String) Date and time when the SCIM integration was created.
- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import snowflake_scim_integration.example name
```
