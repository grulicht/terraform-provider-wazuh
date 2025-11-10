# 🔐 Resource Documentation: `wazuh_role_user`

# wazuh_role_user

The `wazuh_role_user` resource manages **role assignments for Wazuh API users** via:

- `POST /security/users/{user_id}/roles` – assign roles to a user  
- `DELETE /security/users/{user_id}/roles` – remove roles from a user  

It represents the *relationship* between:

- one **user** (`user_id`) and  
- one or more **roles** (`role_ids`).

On `apply`, the resource **links** the given roles to the user.  
On `destroy`, it **unlinks** those roles from the user.

---

## Example Usage

### 1) Assign roles to an existing user

```hcl
resource "wazuh_role_user" "user_roles" {
  user_id  = 5              # existing Wazuh user ID
  role_ids = [2, 3]         # existing Wazuh role IDs
  position = 0              # optional; security position
}
````

### 2) Create role, user & assign the role (typical pattern)

```hcl
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
  position = 0
}

output "wazuh_readonly_role_id" {
  value = wazuh_role.readonly.role_id
}

output "wazuh_user_id" {
  value = wazuh_user.user.user_id
}

output "wazuh_role_user_message" {
  value = wazuh_role_user.user_roles.message
}
```

---

## Lifecycle & Behavior

### Create

On `terraform apply`, the provider calls:

```http
POST /security/users/{user_id}/roles?role_ids=1,2,3&position=0
```

A typical response looks like:

```json
{
  "data": {
    "affected_items": [],
    "total_affected_items": 1,
    "total_failed_items": 0,
    "failed_items": []
  },
  "message": "All roles were linked to TestWazuhUsername",
  "error": 0
}
```

The resource then stores:

* `message`
* `total_affected`
* `total_failed`
* `error_code`
* `timestamp` (UTC)
* and a synthetic `id` (e.g. `user-5-roles-1,2,3`)

### Read

The `Read` operation is effectively a **no-op**:

* It does **not** re-query Wazuh for current role mappings.
* It keeps the information stored from the last `Create`.

> If roles are modified outside Terraform, state may become out of sync.
> In that case, re-apply your desired configuration via Terraform.

### Delete

On `terraform destroy`, the provider calls:

```http
DELETE /security/users/{user_id}/roles?role_ids=1,2,3
```

which removes the **specified roles** from that user (or all roles if you used `"all"` internally).

After a successful delete, the resource is removed from Terraform state.

---

## Arguments Reference

| Name       | Type         | Required    | ForceNew | Description                                                                                                                     |
| ---------- | ------------ | ----------- | -------- | ------------------------------------------------------------------------------------------------------------------------------- |
| `user_id`  | string / int | ✅ **Yes**   | ✅ Yes    | Wazuh user ID to which the roles will be assigned. Changing this forces a new resource.                                         |
| `role_ids` | list(number) | ✅ **Yes**   | ✅ Yes    | List of Wazuh role IDs to link to the user. Changing the set of roles forces a new resource (removal of old + creation of new). |
| `position` | number       | 🚫 optional | ✅ Yes    | Optional *security position* for the roles/policies. Changing this forces a new resource.                                       |

> ℹ️ All arguments are **ForceNew** – changing them will destroy and recreate the mapping, which means **roles will be re-linked according to the new config**.

---

## Attributes Reference

| Name             | Type   | Description                                                                                         |
| ---------------- | ------ | --------------------------------------------------------------------------------------------------- |
| `id`             | string | Synthetic ID representing the mapping (e.g. `user-5-roles-2,3`).                                    |
| `user_id`        | string | The Wazuh user ID this resource is managing.                                                        |
| `role_ids`       | list   | The list of Wazuh role IDs assigned to the user via this resource.                                  |
| `position`       | number | The position value used when creating the link.                                                     |
| `message`        | string | Human-readable message from the Wazuh API (e.g. *"All roles were linked to TestWazuhUsername"*).    |
| `total_affected` | int    | Number of affected items reported by the API.                                                       |
| `total_failed`   | int    | Number of failed items reported by the API.                                                         |
| `error_code`     | int    | Raw error code from the Wazuh API response (`0` = success, `>0` indicates partial or full failure). |
| `timestamp`      | string | UTC timestamp (RFC3339) when the mapping was last applied via Terraform.                            |
