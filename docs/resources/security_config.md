# 🛡️ Resource Documentation: `wazuh_security_config`

# wazuh_security_config

The `wazuh_security_config` resource manages the **global Wazuh security configuration** via the `/security/config` API.

It allows you to:

- read the current **RBAC mode** and **JWT token lifetime**,
- update security settings in a **declarative** way,
- optionally restore the **default security configuration**.

> This is a **singleton** resource – there is only **one** security configuration per Wazuh manager.  
> The resource always uses a fixed ID (`security_config`).

---

## Example Usage

### Set custom token timeout & RBAC mode

```hcl
resource "wazuh_security_config" "config" {
  auth_token_exp_timeout = 1800  # 30 minutes
  rbac_mode              = "white"
}

output "security_config_timeout" {
  value = wazuh_security_config.config.auth_token_exp_timeout
}

output "security_config_rbac_mode" {
  value = wazuh_security_config.config.rbac_mode
}
````

### Read-only (just sync current config into state)

If you don’t want to change anything and only want to **read** what’s configured now:

```hcl
resource "wazuh_security_config" "config" {}
```

When no arguments are set, the provider **does not call** `PUT /security/config` and simply calls `GET /security/config` to populate state.

---

## Lifecycle & Behavior

### Create / Update – `PUT /security/config`

On `terraform apply` / `tofu apply`:

1. The provider builds a JSON payload **only from the fields you set**, for example:

   ```json
   {
     "auth_token_exp_timeout": 1800,
     "rbac_mode": "white"
   }
   ```

2. It sends:

   ```http
   PUT /security/config
   ```

3. After a successful response, the provider:

   * sets the resource ID to `security_config`,
   * sets `auth_token_exp_timeout` and `rbac_mode` in state based on what was sent or returned.

> ⚠️ Wazuh behavior:
> Every change to `/security/config` **revokes all existing JWT tokens**.

#### Special case: no arguments set

If the resource has no arguments set (`auth_token_exp_timeout` nor `rbac_mode`):

* the provider **does not perform a PUT**,
* instead it calls:

  ```http
  GET /security/config
  ```

  and populates the state with the current configuration.

This is useful for a “read-only” mode.

---

### Read – `GET /security/config`

During a regular `terraform plan` / `refresh`:

```http
GET /security/config
```

is called and the Terraform state is updated to match the current Wazuh configuration:

* `auth_token_exp_timeout`
* `rbac_mode`

This works in new Terraform runs (with a fresh JWT token obtained by the provider).

---

### Delete – Restore Defaults (if implemented)

If you run `terraform destroy` on this resource and the provider implements reset behavior, the typical flow is:

* call:

  ```http
  DELETE /security/config
  ```

  (the “Restore default security config” endpoint),
* Wazuh resets the configuration to **default values**,
* the resource is removed from the Terraform state.

> Note: deleting this resource **does not disable security** – it only resets configuration to defaults and removes the Terraform representation.

---

## Arguments Reference

All input arguments are **optional**. If they are not set, the resource behaves effectively as “read-only”.

| Name                     | Type   | Required | ForceNew | Description                                                                                  |
| ------------------------ | ------ | -------- | -------- | -------------------------------------------------------------------------------------------- |
| `auth_token_exp_timeout` | number | No       | ❌ No     | JWT token expiration time in **seconds**. Minimum `30`. Overrides the default (e.g., `900`). |
| `rbac_mode`              | string | No       | ❌ No     | RBAC mode, either `"white"` or `"black"`, as defined in Wazuh security configuration.        |

---

## Attributes Reference

| Name                     | Type   | Description                                                                                 |
| ------------------------ | ------ | ------------------------------------------------------------------------------------------- |
| `id`                     | string | Always `"security_config"`. Singleton resource ID.                                          |
| `auth_token_exp_timeout` | number | Effective configured token expiration time in seconds (from the API or your configuration). |
| `rbac_mode`              | string | Effective RBAC mode (`white` / `black`).                                                    |
