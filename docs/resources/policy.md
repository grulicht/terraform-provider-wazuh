# 📝 Resource Documentation: `wazuh_policy`

# wazuh_policy

The `wazuh_policy` resource manages **Wazuh security policies** through the Wazuh API:

- `POST /security/policies` – create policy  
- `GET /security/policies` – read policy by ID  
- `PUT /security/policies/{policy_id}` – update policy  
- `DELETE /security/policies` – delete policy by ID  

Policies in Wazuh define **which actions are allowed/denied on which resources** and are typically attached to roles.

---

## Example Usage

### Create a Policy

```hcl
resource "wazuh_policy" "normal_policy" {
  name = "normal_policy"

  policy = jsonencode({
    actions   = []
    resources = []
    effect    = "allow"
  })
}
````

### Update Policy Definition

To change the policy actions or resources, just edit `policy` and re-apply:

```hcl
resource "wazuh_policy" "normal_policy" {
  name = "normal_policy"

  policy = jsonencode({
    actions   = ["agent:read", "agent:modify"]
    resources = ["agent:id:*"]
    effect    = "allow"
  })
}
```

### Import Existing Policy

If you already have a policy with ID `7` in Wazuh:

```bash
terraform import wazuh_policy.existing 7
```

Then run:

```bash
terraform plan
```

and align `name` / `policy` in your configuration with what exists in Wazuh.

---

## Lifecycle & Behavior

### Create

On `terraform apply`, the provider:

1. Validates that `policy` is valid JSON.

2. Builds a payload like:

   ```json
   {
     "name": "normal_policy",
     "policy": {
       "actions": [],
       "resources": [],
       "effect": "allow"
     }
   }
   ```

3. Sends:

   ```http
   POST /security/policies
   ```

4. Because the create endpoint does **not** return the policy ID, the provider then issues:

   ```http
   GET /security/policies?search=normal_policy&limit=100
   ```

   and finds an **exact name match** to obtain `policy_id`.

> ⚠️ To avoid ambiguity, `name` should be **unique**. If multiple policies with the same name exist, the provider will return an error during lookup.

### Read

On refresh / plan, the provider calls:

```http
GET /security/policies?policy_ids=<policy_id>&limit=1
```

and updates:

* `policy_id` – numeric ID from Wazuh
* `name` – policy name
* `policy` – JSON representation of the policy body, normalized/compacted

If the policy is not found, the resource is removed from the Terraform state.

### Update

On changes to `name` and/or `policy`, the provider sends:

```http
PUT /security/policies/{policy_id}
```

with a body containing at least one of:

```json
{
  "name": "new_name",
  "policy": {
    "actions": [],
    "resources": [],
    "effect": "allow"
  }
}
```

If nothing has actually changed, no API call is made and the provider just re-reads state.

### Delete

On destroy, the provider calls:

```http
DELETE /security/policies?policy_ids=<policy_id>
```

and removes the resource from the Terraform state when successful.

---

## Arguments Reference

| Name     | Type   | Required  | ForceNew | Description                                                                                            |
| -------- | ------ | --------- | -------- | ------------------------------------------------------------------------------------------------------ |
| `name`   | string | ✅ **Yes** | ❌ No     | Policy name (≤ 64 characters). Should be unique to allow reliable lookup on creation.                  |
| `policy` | string | ✅ **Yes** | ❌ No     | **JSON-encoded** policy definition (`actions`, `resources`, `effect`, etc.), usually via `jsonencode`. |

> 💡 Keeping `policy` as a string lets you compose complex JSON using HCL’s `jsonencode()` without hand-writing JSON.

---

## Attributes Reference

| Name        | Type   | Description                                                        |
| ----------- | ------ | ------------------------------------------------------------------ |
| `id`        | string | Terraform resource ID – same as `policy_id`.                       |
| `policy_id` | string | Numeric Wazuh policy ID as returned by the Wazuh API.              |
| `name`      | string | Policy name as stored in Wazuh.                                    |
| `policy`    | string | JSON representation of the policy body reconstructed from the API. |
