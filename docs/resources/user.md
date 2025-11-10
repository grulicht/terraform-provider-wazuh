# 👤 **Resource Documentation: `wazuh_user`**

# wazuh_user

The `wazuh_user` resource manages **Wazuh API users** via the `/security/users` endpoints.

It allows you to:

* **Create** API users (`POST /security/users`)
* **Read** user details by `user_id` (`GET /security/users?user_ids=...`)
* **Update** user passwords (`PUT /security/users/{user_id}`)
* **Delete** users (`DELETE /security/users?user_ids=...`)
* **Import** existing users into Terraform state

This is ideal for managing Wazuh API access in a GitOps-friendly and auditable way.

---

## Example Usage

### Create a Wazuh API User

```hcl
resource "wazuh_user" "api_user" {
  username = "tf_api_user"
  password = "pa$$w0rd123!"
}
```

### Rotate User Password

You can rotate a user's password simply by changing `password` and re-applying:

```hcl
resource "wazuh_user" "api_user" {
  username = "tf_api_user"
  password = "newSuperSecurePa$$w0rd!"
}
```

Terraform will call:

```http
PUT /security/users/{user_id}
```

with the new password.

### Import Existing User

If you already created a user in Wazuh (e.g. via UI or API), you can import it into Terraform by `user_id`:

```bash
terraform import wazuh_user.api_user 3
```

Then you can set/update the `password` and let Terraform manage it going forward (note: Terraform will **not** read the existing password from Wazuh – you must define it in your config if you want it managed).

---

## Lifecycle & Behavior

### Create – `POST /security/users`

On `terraform apply`, resource creation:

1. Sends:

   ```http
   POST /security/users
   Content-Type: application/json

   {
     "username": "tf_api_user",
     "password": "pa$$w0rd123!"
   }
   ```

2. The Wazuh API responds with a success message, but does **not** directly return `user_id`.

3. The provider then calls:

   ```http
   GET /security/users?search=tf_api_user&limit=100
   ```

   and finds the exact match by `username` to determine the numeric `user_id`.

4. The resource `id` is set to this `user_id`, and `user_id` is exposed as an attribute.

> ⚠️ If multiple users exist with the same `username`, creation will fail with an error

---

### Read – `GET /security/users?user_ids=<id>`

On refresh/plan, the provider:

```http
GET /security/users?user_ids=<user_id>&limit=1
```

* If the user exists:

  * updates `user_id` and `username`,
  * leaves `password` as-is in Terraform state (Wazuh never returns it).
* If the user no longer exists:

  * clears `id` → Terraform will treat it as tainted/removed.

---

### Update – `PUT /security/users/{user_id}`

Only the **password** is updatable.

If you change `password` in your Terraform configuration:

```hcl
resource "wazuh_user" "api_user" {
  username = "tf_api_user"
  password = "newSuperSecurePa$$w0rd!"
}
```

then the provider calls:

```http
PUT /security/users/{user_id}
Content-Type: application/json

{
  "password": "newSuperSecurePa$$w0rd!"
}
```

* HTTP status code is used to determine success.
* No username changes are allowed – `username` is **ForceNew** (change => recreate).

---

### Delete – `DELETE /security/users?user_ids=<id>`

On `terraform destroy`, the provider calls:

```http
DELETE /security/users?user_ids=<user_id>
```

If successful, the user is removed from Wazuh and from Terraform state.

To delete multiple users at once, just define more `wazuh_user` resources.

---

## Arguments Reference

| Name       | Type   | Required    | ForceNew | Description                                                                                                             |
| ---------- | ------ | ----------- | -------- | ----------------------------------------------------------------------------------------------------------------------- |
| `username` | string | ✅ **Yes**  | ✅ Yes    | Wazuh API username (4–64 characters). Changing this forces a new resource (user recreation).                           |
| `password` | string | ✅ **Yes**  | ❌ No     | Password for the Wazuh API user. Required on create. Changing this updates the existing user's password. **Sensitive**.|

---

## Attributes Reference

| Name       | Type   | Description                                                                  |
| ---------- | ------ | ---------------------------------------------------------------------------- |
| `id`       | string | Terraform resource ID. Equals the Wazuh numeric `user_id`.                   |
| `user_id`  | string | Numeric Wazuh user ID (e.g. `"3"`). Same value as `id`.                      |
| `username` | string | Username of the Wazuh API user.                                              |
| `password` | string | Password known to Terraform (never read back from the Wazuh API). Sensitive. |
