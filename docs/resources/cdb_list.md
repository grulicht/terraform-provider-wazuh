# 🧭 **Resource Documentation: `wazuh_cdb_list`**

# wazuh_cdb_list

The `wazuh_cdb_list` resource allows you to **manage Wazuh CDB list files** via the `/lists/files/{filename}` API endpoints.

CDB lists are used by Wazuh for things like IP reputation, malicious domains, AWS event names, and other lookup-based detections.
With this resource you can **create**, **update**, **read**, **import**, and **delete** CDB lists directly from Terraform.

---

## Example Usage

### Create a Custom CDB List

```hcl
resource "wazuh_cdb_list" "example" {
  filename  = "test_list"
  overwrite = true

  content = <<EOF
192.168.1.10:allowed
192.168.1.11:denied
EOF
}
```

After apply, you should see a new list in the Wazuh Dashboard under **CDB Lists**, with `Name: test_list`.

---

### Import an Existing CDB List

If a CDB list already exists in Wazuh (e.g. `malicious-ip`), you can import it:

```bash
terraform import wazuh_cdb_list.malicious_ip malicious-ip
```

And define it in Terraform:

```hcl
resource "wazuh_cdb_list" "malicious_ip" {
  filename  = "malicious-ip"
  overwrite = true

  content = <<EOF
1.2.3.4:bad
5.6.7.8:bad
EOF
}
```

Terraform will:

1. Use `filename` as `malicious-ip`
2. Read the current content via `GET /lists/files/malicious-ip?raw=true`
3. Keep it in state so you can manage it as IaC

---

## Lifecycle & Behavior

* Creating or updating a `wazuh_cdb_list` resource triggers:

  ```http
  PUT /lists/files/{filename}?overwrite=true|false
  Content-Type: application/octet-stream
  ```

  where `{filename}` is the **base name** of the CDB list (e.g. `malicious-ip`).

* Reading a list uses:

  ```http
  GET /lists/files/{filename}?raw=true
  ```

  and returns the **raw file content**, which is stored in the `content` attribute.

* Deleting a list uses:

  ```http
  DELETE /lists/files/{filename}
  ```

Terraform automatically handles:

* Create (`PUT`)
* Read/refresh (`GET`)
* Update (`PUT`)
* Delete (`DELETE`)
* Import by `filename`

> ⚠️ **Important:** `filename` is just the **list name**, e.g. `malicious-ip`, `security-eventchannel`, `aws-sources`.
> Do **not** include paths like `etc/lists` or file extensions like `.cdb`; Wazuh resolves those internally.

---

## Arguments Reference

| Name        | Type   | Required    | Description                                                                                               |
| ----------- | ------ | ----------- | --------------------------------------------------------------------------------------------------------- |
| `filename`  | string | ✅ **Yes**   | Name of the CDB list file (e.g. `malicious-ip`, `security-eventchannel`). Only the filename, no path.     |
| `content`   | string | ✅ **Yes**   | Full content of the CDB list file (plain text). Each line typically contains `key:value` pairs.           |
| `overwrite` | bool   | 🚫 optional | Whether to overwrite the file if it already exists. Maps to `overwrite` query parameter. Default: `true`. |

---

## Attributes Reference

| Name       | Description                                 |
| ---------- | ------------------------------------------- |
| `id`       | Resource ID (same as `filename`).           |
| `filename` | The CDB list name managed by this resource. |
| `content`  | Current file content as stored in Wazuh.    |
