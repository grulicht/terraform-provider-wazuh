# 📝🧩 **Resource Documentation: `wazuh_policy_role`**

# wazuh_policy_role

The `wazuh_policy_role` resource manages the **relationship between Wazuh roles and policies** via the Wazuh Security API:

* `POST /security/roles/{role_id}/policies`
* `DELETE /security/roles/{role_id}/policies`

Use it to:

* assign one or more **policies** to a **role**,
* keep role–policy mappings under Terraform control,
* cleanly remove those mappings on `terraform destroy`.

> One **role** can have **multiple policies**, and this resource is exactly for modelling that link.

---

## Example Usage

### Basic: Attach a Policy to a Role

```hcl
resource "wazuh_policy" "normal_policy" {
  name = "normal_policy"

  policy = jsonencode({
    actions   = ["agent:read"]
    resources = ["agent:id:*"]
    effect    = "allow"
  })
}

resource "wazuh_role" "readonly" {
  name = "readonly"
}

resource "wazuh_policy_role" "readonly_policy" {
  role_id = wazuh_role.readonly.role_id

  # policy_ids from the API are numeric
  policy_ids = [
    tonumber(wazuh_policy.normal_policy.policy_id)
  ]

  # optional ordering/priority
  position = 0
}
```

---

### Full Example: User + Role + Policy

```hcl
variable "wazuh_policy_name"          { default = "normal_policy" }
variable "wazuh_policy_actions"       { default = ["agent:read"] }
variable "wazuh_policy_resources"     { default = ["agent:id:*"] }
variable "wazuh_policy_effect"        { default = "allow" }

variable "wazuh_role_name"            { default = "readonly" }
variable "wazuh_user_username"        { default = "terraform-user" }
variable "wazuh_user_password"        { default = "SuperSecret123!" }

variable "wazuh_role_user_position"   { default = 0 }
variable "wazuh_policy_role_position" { default = 0 }

resource "wazuh_policy" "normal_policy" {
  name = var.wazuh_policy_name

  policy = jsonencode({
    actions   = var.wazuh_policy_actions
    resources = var.wazuh_policy_resources
    effect    = var.wazuh_policy_effect
  })
}

resource "wazuh_role" "readonly" {
  name = var.wazuh_role_name
}

resource "wazuh_user" "user" {
  username = var.wazuh_user_username
  password = var.wazuh_user_password
}

resource "wazuh_role_user" "user_roles" {
  user_id  = wazuh_user.user.user_id
  role_ids = [wazuh_role.readonly.role_id]
  position = var.wazuh_role_user_position
}

resource "wazuh_policy_role" "readonly_policy" {
  role_id = wazuh_role.readonly.role_id

  policy_ids = [
    tonumber(wazuh_policy.normal_policy.policy_id)
  ]

  position = var.wazuh_policy_role_position
}
```

---

## Lifecycle & Behavior

### Create – Add Policies to Role

On `terraform apply`, the resource calls:

```http
POST /security/roles/{role_id}/policies?policy_ids=10,11&position=0
```

with:

* `role_id` from the resource,
* `policy_ids` joined as a comma-separated list,
* optional `position` if provided.

Typical response:

```json
{
  "data": {
    "affected_items": [],
    "total_affected_items": 2,
    "total_failed_items": 0,
    "failed_items": []
  },
  "message": "All policies were linked to role 100",
  "error": 0
}
```

The provider then:

* stores `role_id` and `policy_ids` in state,
* sets `id` (usually derived from `role_id` + hash of `policy_ids`).

### Read – No-op

This resource represents a **relationship / side-effect**, not a full object with a dedicated `GET` endpoint.

So the `Read` step is effectively a **no-op** and just keeps current state (no extra API call).

### Delete – Remove Policies from Role

On `terraform destroy`, the resource calls:

```http
DELETE /security/roles/{role_id}/policies?policy_ids=10,11
```

The mapping between that role and those policies is removed.
The **role** and the **policy objects themselves** are **not** deleted – only the relation.

---

## Arguments Reference

| Name         | Type         | Required    | ForceNew | Description                                                                                                       |
| ------------ | ------------ | ----------- | -------- | ----------------------------------------------------------------------------------------------------------------- |
| `role_id`    | string       | ✅ **Yes**   | ✅ Yes    | ID of the Wazuh role. Usually obtained from `wazuh_role.<name>.role_id`. Changing this forces a new relationship. |
| `policy_ids` | list(number) | ✅ **Yes**   | ✅ Yes    | List of Wazuh policy IDs to link to this role. Changing this list forces recreation (adds/removes mappings).      |
| `position`   | number       | 🚫 optional | ✅ Yes    | Optional security position / order for the role–policy relation. Changing this forces recreation of the resource. |

> ℹ️ `policy_id` values from `wazuh_policy` are strings, while `policy_ids` here are `list(number)`,

---

## Attributes Reference

| Name         | Type         | Description                                                     |
| ------------ | ------------ | --------------------------------------------------------------- |
| `id`         | string       | Internal Terraform ID for this mapping (role + policies combo). |
| `role_id`    | string       | Role ID used in the relation.                                   |
| `policy_ids` | list(number) | List of policy IDs linked to the role.                          |
| `position`   | number       | Position value (if set).                                        |
