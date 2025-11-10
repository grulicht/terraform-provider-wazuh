# 🛡️ Resource Documentation: `wazuh_security_rule`

# wazuh_security_rule

The `wazuh_security_rule` resource manages **Wazuh security rules** via the Wazuh API:

- `POST /security/rules` – create rule  
- `GET /security/rules` – read rule by ID  
- `PUT /security/rules/{rule_id}` – update rule  
- `DELETE /security/rules` – delete rule by ID  

Security rules are used to define fine-grained access control in Wazuh. They are typically attached to roles.

---

## Example Usage

### Create a Security Rule

```hcl
resource "wazuh_security_rule" "normal_rule" {
  name = "normalRule"

  rule = jsonencode({
    MATCH = {
      definition = "normalRule"
    }
  })
}
````

### Update Rule Body / Name

Just change `name` or `rule` and re-apply:

```hcl
resource "wazuh_security_rule" "normal_rule" {
  name = "normalRuleRenamed"

  rule = jsonencode({
    MATCH = {
      definition = "anotherDefinition"
    }
  })
}
```

### Import Existing Security Rule

If you already have a security rule with ID `5` in Wazuh:

```bash
terraform import wazuh_security_rule.existing 5
```

Then run:

```bash
terraform plan
```

and optionally align `name` and `rule` in your configuration with what exists in Wazuh.

---

## Lifecycle & Behavior

### Create

On `terraform apply`, the provider:

1. Builds a payload:

   ```json
   {
     "name": "<name>",
     "rule": { ...parsed from rule JSON string... }
   }
   ```

2. Sends:

   ```http
   POST /security/rules
   ```

3. Since the create endpoint does **not** return the rule ID, the provider then does:

   ```http
   GET /security/rules?search=<name>&limit=100
   ```

   and finds an **exact name match** to obtain `rule_id`.

> ⚠️ The `name` should be unique, otherwise the provider will error if multiple rules with the same name are found during lookup.

### Read

On refresh / plan, the provider calls:

```http
GET /security/rules?rule_ids=<rule_id>&limit=1
```

and updates:

* `name` – from the API
* `rule` – from the `rule` body, re-encoded as a JSON string
* `rule_id` – from the API

If the rule is not found, the resource is removed from state.

### Update

On change to `name` or `rule`, the provider sends:

```http
PUT /security/rules/{rule_id}
```

with a body containing one or both of:

```json
{
  "name": "<new-name>",
  "rule": { ...parsed from rule JSON string... }
}
```

The API requires at least one field to be present; if nothing changed, no request is sent.

### Delete

On destroy, the provider calls:

```http
DELETE /security/rules?rule_ids=<rule_id>
```

and removes the resource from the Terraform state when successful.

---

## Arguments Reference

| Name   | Type   | Required  | ForceNew | Description                                                                                   |
| ------ | ------ | --------- | -------- | --------------------------------------------------------------------------------------------- |
| `name` | string | ✅ **Yes** | ❌ No     | Security rule name (≤ 64 characters). Must be unique to allow reliable lookup on create.      |
| `rule` | string | ✅ **Yes** | ❌ No     | **JSON-encoded** rule body as expected by Wazuh (e.g. via `jsonencode({ MATCH = { ... } })`). |

> 💡 The `rule` attribute is kept as a **string** so you can use `jsonencode(...)` from HCL and avoid manual JSON formatting.

---

## Attributes Reference

| Name      | Type   | Description                                                      |
| --------- | ------ | ---------------------------------------------------------------- |
| `id`      | string | Terraform resource ID – same as `rule_id`.                       |
| `rule_id` | string | Numeric Wazuh security rule ID as returned by the Wazuh API.     |
| `name`    | string | Rule name as stored in Wazuh.                                    |
| `rule`    | string | JSON representation of the rule body reconstructed from the API. |
