# 🕵️ **Resource Documentation: `wazuh_decoder`**

# wazuh_decoder

The `wazuh_decoder` resource allows you to manage **Wazuh decoder XML files** via the `/decoders/files/{filename}` API endpoints.

Decoders in Wazuh are used to parse, normalize and enrich incoming logs.
With this resource you can **create**, **update**, **read**, **import**, and **delete** decoder files directly from Terraform.

---

## Example Usage

### Create a Custom Decoder

```hcl
resource "wazuh_decoder" "local_example" {
  filename  = "local_decoder.xml"
  overwrite = true

  content = <<EOF
<group name="local,syslog,">
  <decoder name="local-example">
    <prematch>my-app</prematch>
  </decoder>
</group>
EOF
}
```

After `terraform apply`, you should see `local_decoder.xml` among decoders in the Wazuh UI.

---

### Import an Existing Decoder File

If you already have a decoder file (e.g. `local_decoder.xml`), you can import it:

```bash
terraform import wazuh_decoder.local_example local_decoder.xml
```

Terraform will then:

1. Set `id = "local_decoder.xml"`,
2. Load the file content via `GET /decoders/files/local_decoder.xml?raw=true`,
3. Store it into the `content` attribute in state.

You can then copy that XML into your `.tf` file and start managing it as code.

---

## Lifecycle & Behavior

* Creating or updating a `wazuh_decoder` resource triggers:

  ```http
  PUT /decoders/files/{filename}?overwrite=true|false[&relative_dirname=...]
  Content-Type: application/octet-stream
  ```

  where:

  * `{filename}` is the decoder XML filename (e.g. `local_decoder.xml`),
  * `relative_dirname` is optional (see below).

* Reading a decoder uses:

  ```http
  GET /decoders/files/{filename}?raw=true[&relative_dirname=...]
  ```

  and returns the **raw XML content**, which is stored in `content`.

* Deleting a decoder uses:

  ```http
  DELETE /decoders/files/{filename}[?relative_dirname=...]
  ```

Terraform automatically handles:

* Create/Update (`PUT`),
* Read/refresh (`GET`),
* Delete (`DELETE`),
* Import by filename.

> ⚠️ **Note:** If you change the `filename`, Terraform will destroy the old decoder file and create a new one.

---

## Arguments Reference

| Name               | Type   | Required    | Description                                                                                                                                |
| ------------------ | ------ | ----------- | ------------------------------------------------------------------------------------------------------------------------------------------ |
| `filename`         | string | ✅ **Yes**   | Decoder XML filename (e.g. `local_decoder.xml`). Only the filename, not an absolute path.                                                  |
| `content`          | string | ✅ **Yes**   | Full XML content of the decoder file.                                                                                                      |
| `overwrite`        | bool   | 🚫 optional | Whether to overwrite the file if it already exists. Maps to `overwrite` query parameter. Default: `true`.                                  |
| `relative_dirname` | string | 🚫 optional | Optional relative directory name where the decoder resides, e.g. `ruleset/decoders`. If omitted, Wazuh uses its default decoder directory. |

---

## Attributes Reference

| Name               | Description                                                     |
| ------------------ | --------------------------------------------------------------- |
| `id`               | Resource ID (same as `filename`).                               |
| `filename`         | Decoder filename managed by this resource.                      |
| `content`          | Current XML content as stored in Wazuh.                         |
| `relative_dirname` | Directory context used when interacting with the file (if set). |
