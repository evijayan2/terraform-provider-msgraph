---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "msgraph_application_web Resource - terraform-provider-msgraph"
subcategory: ""
description: |-
  Application Web RedirectURI config resource
---

# msgraph_application_web (Resource)

Application Web RedirectURI config resource



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `app_id` (String) Application Client ID
- `redirect_uri` (String) Redirect URI attribute

### Read-Only

- `application` (Object) Application data (see [below for nested schema](#nestedatt--application))
- `id` (String) identifier

<a id="nestedatt--application"></a>
### Nested Schema for `application`

Read-Only:

- `appId` (String)
- `displayName` (String)
- `id` (String)
- `web` (Object) (see [below for nested schema](#nestedobjatt--application--web))

<a id="nestedobjatt--application--web"></a>
### Nested Schema for `application.web`

Read-Only:

- `implicitGrantSettings` (Object) (see [below for nested schema](#nestedobjatt--application--web--implicitGrantSettings))
- `redirectUris` (List of String)

<a id="nestedobjatt--application--web--implicitGrantSettings"></a>
### Nested Schema for `application.web.implicitGrantSettings`

Read-Only:

- `enableAccessTokenIssuance` (Boolean)
- `enableIdTokenIssuance` (Boolean)


