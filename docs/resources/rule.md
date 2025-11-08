# 📜 **Resource Documentation: `wazuh_rule`**

# wazuh_rule

The `wazuh_rule` resource allows you to manage **custom Wazuh ruleset files** via the `/rules/files/{filename}` API.

With this resource you can:

* upload / update a **rules XML file**,
* fetch its current content,
* delete the rules file when no longer needed.

It uses the Wazuh “rules files” API and is conceptually similar to `wazuh_decoder`, but for rules.

---

## Example Usage

### Create / Manage a Local Rules File

```hcl
resource "wazuh_rule" "local_rules" {
  filename  = "local_rules.xml"
  overwrite = true

  content = <<-EOT
    <group name="local,">
      <rule id="100100" level="3">
        <if_sid>100000</if_sid>
        <match>my-custom-rule</match>
        <description>My custom local rule</description>
      </rule>
    </group>
  EOT
}
```

* On `terraform apply` the file `local_rules.xml` will be uploaded (or overwritten) in the Wazuh ruleset.
* On `terraform destroy` the rules file will be deleted from Wazuh.

---

### Import Existing Rules File

If you already have a rules file in Wazuh (e.g. `local_rules.xml`) and want to manage it with Terraform:

```bash
terraform import wazuh_rule.local_rules local_rules.xml
```

> After import, run `terraform plan` to see differences between remote content and your `content` field.

---

## Lifecycle & Behavior

### Create / Update – Upload Rules File

`Create` and `Update` both:

* Call:

  ```http
  PUT /rules/files/{filename}?overwrite=true|false
  Content-Type: application/octet-stream
  ```

* Request body is the raw XML from `content`.

Example behavior:

* If `overwrite = true`:

  * Existing file with the same name will be replaced.
* If `overwrite = false`:

  * If the file already exists, Wazuh returns an error and Terraform will fail the apply.

A successful response looks like:

```json
{
  "data": {
    "affected_items": [],
    "failed_items": [],
    "total_affected_items": 1,
    "total_failed_items": 0
  },
  "message": "Rule was successfully uploaded",
  "error": 0
}
```

`filename` is used as the Terraform resource `id`.

---

### Read – Get Rules File Content

On refresh/plan, Terraform calls:

```http
GET /rules/files/{filename}?raw=true
```

With `raw=true`, Wazuh returns the **plain XML content** of the rules file in the response body.
The provider then updates the `content` attribute so Terraform state stays in sync with what’s actually in Wazuh.

If the file no longer exists, Wazuh returns `404` and the resource is removed from state.

---

### Delete – Remove Rules File

On `terraform destroy`, the resource calls:

```http
DELETE /rules/files/{filename}
```

On success, the rules file is removed from Wazuh and Terraform forgets the resource.

> If Wazuh returns `404` (already deleted), the provider treats it as success and still removes the resource from the state.

---

## Arguments Reference

| Name        | Type   | Required    | Description                                                                                   |
| ----------- | ------ | ----------- | --------------------------------------------------------------------------------------------- |
| `filename`  | string | ✅ **Yes**   | Name of the Wazuh rules file, e.g. `local_rules.xml`. This is also used as the resource `id`.|
| `content`   | string | ✅ **Yes**   | Full XML content of the rules file.                                                          |
| `overwrite` | bool   | 🚫 optional | Whether to overwrite an existing rules file with the same name. Default: `true`.              |
| `relative_dirname` | string | 🚫 optional | Reserved for future use. Currently **ignored by the provider** and **not sent to the Wazuh API**; Wazuh’s default rules location is used. |

> 🔎 **Note:** The Wazuh API supports a `relative_dirname` parameter, but this provider currently relies on Wazuh’s default rules path and does **not** send a custom directory. All managed rules files are uploaded into the default ruleset location as defined by Wazuh.

---

## Attributes Reference

| Name        | Type   | Description                                                                           |
| ----------- | ------ | ------------------------------------------------------------------------------------- |
| `id`        | string | Resource ID (equal to `filename`).                                                    |
| `filename`  | string | Name of the Wazuh rules file.                                                         |
| `content`   | string | Effective XML content currently known in Terraform state (synced from Wazuh on read). |
| `overwrite` | bool   | Whether overwriting is allowed when uploading.                                        |
