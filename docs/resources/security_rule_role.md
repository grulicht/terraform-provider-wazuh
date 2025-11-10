# 🛡️🧩 **Resource Documentation: `wazuh_security_rule_role`**

# wazuh_security_rule_role

The `wazuh_security_rule_role` resource represents a **relationship** between a Wazuh **role** and one or more **security rules**.

Under the hood it uses:

* `POST /security/roles/{role_id}/rules` – to **link** security rules to a role
* `DELETE /security/roles/{role_id}/rules` – to **unlink** security rules from a role

Use this resource when you want to:

* Attach one or more `security rules` to a `role`
* Manage role–rule relationships declaratively with Terraform/OpenTofu
* Cleanly remove these links on `destroy`

> ℹ️ The security rules themselves are managed using the `wazuh_security_rule` resource.

---

## Example Usage

### Minimal Example – Attach One Security Rule to a Role

```hcl
resource "wazuh_role" "readonly" {
  name = "readonly-role"
}

resource "wazuh_security_rule" "normal_rule" {
  name = "NormalRule"

  rule = jsonencode({
    MATCH = {
      definition = "normalRule"
    }
  })
}

resource "wazuh_security_rule_role" "readonly_security_rule" {
  role_id = wazuh_role.readonly.role_id

  rule_ids = [
    tonumber(wazuh_security_rule.normal_rule.rule_id)
  ]
}
```

### Combined RBAC Example (Role, Policy, User, Security Rule)

```hcl
# Policy
resource "wazuh_policy" "normal_policy" {
  name = var.wazuh_policy_name

  policy = jsonencode({
    actions   = var.wazuh_policy_actions
    resources = var.wazuh_policy_resources
    effect    = var.wazuh_policy_effect
  })
}

# Role
resource "wazuh_role" "readonly" {
  name = var.wazuh_role_name
}

# User
resource "wazuh_user" "user" {
  username = var.wazuh_user_username
  password = var.wazuh_user_password
}

# User ↔ Role
resource "wazuh_role_user" "user_roles" {
  user_id  = wazuh_user.user.user_id
  role_ids = [wazuh_role.readonly.role_id]
  position = var.wazuh_role_user_position
}

# Role ↔ Policy
resource "wazuh_policy_role" "readonly_policy" {
  role_id = wazuh_role.readonly.role_id

  policy_ids = [
    tonumber(wazuh_policy.normal_policy.policy_id)
  ]

  position = var.wazuh_policy_role_position
}

# Security rule
resource "wazuh_security_rule" "normal_rule" {
  name = var.wazuh_security_rule_name

  rule = jsonencode({
    MATCH = {
      definition = var.wazuh_security_rule_definition
    }
  })
}

# Role ↔ Security rule
resource "wazuh_security_rule_role" "readonly_security_rule" {
  role_id = wazuh_role.readonly.role_id

  rule_ids = [
    tonumber(wazuh_security_rule.normal_rule.rule_id)
  ]
}
```

---

## Lifecycle & Behavior

### Create

On `terraform apply`, the provider calls:

```http
POST /security/roles/{role_id}/rules?rule_ids=<comma-separated-ids>
```

Example:

```http
POST /security/roles/100/rules?rule_ids=10,11
```

A typical response:

```json
{
  "data": {
    "affected_items": [],
    "total_affected_items": 1,
    "total_failed_items": 0,
    "failed_items": []
  },
  "message": "All security rules were linked to role 100",
  "error": 0
}
```

The resource:

* Stores `role_id` and `rule_ids` in state
* Generates an `id` for the relationship internally (usually based on role and rule IDs)

### Read

The `Read` operation typically:

* Uses `GET /security/rules?rule_ids=...` and/or provider logic to validate the mapping
* If the role or referenced rules no longer exist, the resource is cleared from state.

(The exact behavior depends on implementation, but from Terraform perspective it’s treated as a “mapping” that can disappear.)

### Delete

On `terraform destroy`, the provider calls:

```http
DELETE /security/roles/{role_id}/rules?rule_ids=<comma-separated-ids>
```

or with `all` if that’s how the resource is modeled.

Example:

```http
DELETE /security/roles/100/rules?rule_ids=10,11
```

Wazuh responds with something like:

```json
{
  "data": {
    "affected_items": [],
    "total_affected_items": 2,
    "total_failed_items": 0,
    "failed_items": []
  },
  "message": "All security rules were unlinked from role 100",
  "error": 0
}
```

---

## Arguments Reference

| Name       | Type         | Required  | ForceNew | Description                                                                                               |
| ---------- | ------------ | --------- | -------- | --------------------------------------------------------------------------------------------------------- |
| `role_id`  | number / int | ✅ **Yes** | ✅ Yes    | Wazuh role ID to which the security rules are attached. Changing this forces a new mapping resource.      |
| `rule_ids` | list(number) | ✅ **Yes** | ✅ Yes    | List of Wazuh security rule IDs to link to the role. Changing this list forces recreation of the mapping. |

> 💡 Typically you’ll obtain `role_id` and `rule_ids` from other resources, e.g. `wazuh_role` and `wazuh_security_rule`, using `tonumber()` on their exported IDs.

---

## Attributes Reference

| Name       | Type         | Description                                                             |
| ---------- | ------------ | ----------------------------------------------------------------------- |
| `id`       | string       | Internal ID of this role–security-rule relationship (provider-managed). |
| `role_id`  | number       | The role ID this mapping is associated with.                            |
| `rule_ids` | list(number) | The list of security rule IDs linked to the role.                       |
