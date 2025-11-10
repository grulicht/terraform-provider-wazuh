# 🧩 **Resource Documentation: `wazuh_role`**

# wazuh_role

The `wazuh_role` resource manages **Wazuh security roles** via the `/security/roles` API endpoints.

It allows you to:

* **Create** roles (`POST /security/roles`)
* **Read** role details by `role_id` (`GET /security/roles?role_ids=...`)
* **Update** role names (`PUT /security/roles/{role_id}`)
* **Delete** roles (`DELETE /security/roles?role_ids=...`)
* **Import** existing roles into Terraform state

> ⚠️ This resource **does not** manage policies attached to roles.
> Policies are handled via their own endpoints and should be modelled with separate resources in the future.

---

## Example Usage

### Create a Role

```hcl
resource "wazuh_role" "readonly" {
  name = "tf_readonly_role"
}
```

### Rename an Existing Role

Renaming is supported via `PUT /security/roles/{role_id}` – just change the `name` and re-apply:

```hcl
resource "wazuh_role" "readonly" {
  name = "tf_readonly_role_v2"
}
```

Terraform will update the existing role in Wazuh rather than recreating it.

### Import Existing Role

If you already have a role in Wazuh (created through the UI or API), you can import it by its numeric `role_id`:

```bash
terraform import wazuh_role.readonly 3
```

Then you can set `name` in the configuration to match the existing one or update it as needed.

---

## Lifecycle & Behavior

### Create – `POST /security/roles`

On `terraform apply`, when creating a role, the provider sends:

```http
POST /security/roles
Content-Type: application/json

{
  "name": "tf_readonly_role"
}
```

The Wazuh API returns a generic success payload (without directly giving `role_id`), so the provider:

1. Calls:

   ```http
   GET /security/roles?search=tf_readonly_role&limit=100
   ```

2. Filters results for an **exact match** on `name`.

3. Expects exactly one matching role; otherwise it fails:

   * No match → error: could not determine `role_id`.
   * Multiple matches → error: ambiguous role name.

Once the numeric `role_id` has been found:

* `id` is set to that value,
* `role_id` attribute is populated,
* `Read` is executed to sync the latest state.

---

### Read – `GET /security/roles?role_ids=<id>`

On refresh/plan, the provider calls:

```http
GET /security/roles?role_ids=<role_id>&limit=1
```

* If the role exists:

  * `role_id` remains set,
  * `name` is updated from the API payload.
* If the role no longer exists:

  * `id` is cleared → Terraform will treat the resource as removed.

A typical response looks like:

```json
{
  "data": {
    "affected_items": [
      {
        "id": 3,
        "name": "tf_readonly_role"
      }
    ],
    "total_affected_items": 1,
    "total_failed_items": 0,
    "failed_items": []
  },
  "message": "All specified roles were returned",
  "error": 0
}
```

---

### Update – `PUT /security/roles/{role_id}`

Only the **role name** is updated by this resource.

If `name` changes, Terraform calls:

```http
PUT /security/roles/{role_id}
Content-Type: application/json

{
  "name": "tf_readonly_role_v2"
}
```

* HTTP status codes are used to determine success.
* `role_id` / `id` never change during update.

> ℹ️ Associated policies are **not** modified by this endpoint or resource — only the role’s name is updated.

---

### Delete – `DELETE /security/roles?role_ids=<id>`

On `terraform destroy`, the provider calls:

```http
DELETE /security/roles?role_ids=<role_id>
```

If successful:

* the role is removed from Wazuh,
* the resource is removed from Terraform state.

> ⚠️ As per Wazuh API docs, **policies linked to roles are not removed** when a role is deleted.

---

## Arguments Reference

| Name   | Type   | Required  | Description                                                                     |
| ------ | ------ | --------- | ------------------------------------------------------------------------------- |
| `name` | string | ✅ **Yes** | Role name (1–64 characters). This is the **display name** of the role in Wazuh. |

> ✏️ `name` is **updatable** – changing it will trigger a `PUT /security/roles/{role_id}` request, not a recreate.

---

## Attributes Reference

| Name      | Type   | Description                                                |
| --------- | ------ | ---------------------------------------------------------- |
| `id`      | string | Terraform resource ID. Equals the Wazuh numeric `role_id`. |
| `role_id` | string | Numeric Wazuh role ID (e.g. `"3"`). Same as `id`.          |
| `name`    | string | Current role name as stored in Wazuh.                      |
