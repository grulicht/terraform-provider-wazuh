# 🛡️ Wazuh Terraform/OpenTofu Provider
A Terraform/OpenTofu provider to manage [Wazuh](https://wazuh.com/) resources via its API using Terraform/OpenTofu.

It supports provisioning and configuration of Wazuh users and will be extended to support other objects such as groups, hosts, templates, triggers, users etc.

## 🏷️ Provider Support
| Provider       | Provider Support Status              |
|----------------|--------------------------------------|
| [Terraform](https://registry.terraform.io/providers/grulicht/wazuh/latest)      | ![Done](https://img.shields.io/badge/status-done-brightgreen)           |
| [OpenTofu](https://search.opentofu.org/provider/grulicht/wazuh/latest)          | ![Done](https://img.shields.io/badge/status-done-brightgreen) |

## ⚙️ **Example Provider Configuration**

```hcl
provider "wazuh" {
  endpoint        = "https://wazuh.example.com:55000"
  user            = "wazuh-wui"
  password        = "MyS3cr37P450r.*-"
  skip_ssl_verify = true    # optional (default value is `false`)
}
```

### 💡 Notes:

* The **default Wazuh API port** is `55000`.
* Authentication uses **JWT tokens**, automatically obtained by the provider via `/security/user/authenticate`.
* Token expiration defaults to 900 seconds (15 minutes). The provider will refresh it automatically in a future release.
* Use `skip_ssl_verify = true` only for local testing with self-signed certificates.
* [Docs for change password of default user (wazuh-wui) for Wazuh API.](https://documentation.wazuh.com/current/deployment-options/docker/changing-default-password.html)

---

## 🔐 **Authentication**

The Wazuh provider supports **basic authentication** (username/password), which internally retrieves a **JWT bearer token** from the API:

```
POST /security/user/authenticate
→ Authorization: Basic base64(username:password)
← Response: { "data": { "token": "<JWT_TOKEN>" } }
```

This token is then attached to every request as:

```
Authorization: Bearer <JWT_TOKEN>
```

### Example with Environment Variables

```bash
export WAZUH_ENDPOINT="https://localhost:55000"
export WAZUH_USER="wazuh-wui"
export WAZUH_PASSWORD="MyS3cr37P450r.*-"
export WAZUH_SKIP_SSL_VERIFY=true
```

---

## 🧩 **Arguments Reference**

| Name              | Type    | Required | Description                                                                         |
| ----------------- | ------- | -------- | ----------------------------------------------------------------------------------- |
| `endpoint`        | string  | ✅ Yes    | Full URL of the Wazuh API endpoint (e.g. `https://localhost:55000`).                |
| `user`            | string  | ✅ Yes    | Username for Wazuh API authentication (e.g. `wazuh-wui`).                           |
| `password`        | string  | ✅ Yes    | Password for the API user.                                                          |
| `skip_ssl_verify` | boolean | ❌ No     | Skip TLS certificate verification (useful for self-signed certs). Default: `false`. |

---

## 🧱 Supported Resources
| Resource                                       | Status                                                                |
|------------------------------------------------|-----------------------------------------------------------------------|
| `wazuh_group`                                  | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_configuration`                          | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_active_response`                        | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_event`                                  | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_cdb_list`                               | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_decoder`                                | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_node_configuration`                     | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_node_restart`                           | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_node_analysisd_reload`                  | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_logtest`                                | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_rootcheck`                              | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_syscheck`                               | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_rule`                                   | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_manager_configuration`                  | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_manager_restart`                        | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_agent`                                  | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_agent_group`                            | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_agent_node_restart`                     | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_agent_reconnect`                        | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_agent_restart`                          | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_agent_restart_group`                    | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_agent_upgrade`                          | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_agent_upgrade_custom`                   | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_policy`                                 | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_policy_role`                            | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_role`                                   | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_role_user`                              | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_security_config`                        | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_security_rule`                          | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_security_rule_role`                     | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |
| `wazuh_user`                                   | ![Done](https://img.shields.io/badge/status-done-brightgreen)        |

---

### 💡 Missing a resource?
Is there a Wazuh resource you'd like to see supported?

👉 [Open an issue](https://github.com/grulicht/terraform-provider-wazuh/issues/new?template=feature_request.md) and we’ll consider it for implementation — or even better, submit a [Pull Request](https://github.com/grulicht/terraform-provider-wazuh/pulls) to contribute directly!

📘 See [CONTRIBUTING.md](https://github.com/grulicht/terraform-provider-wazuh/blob/main/.github/CONTRIBUTING.md) for guidelines.

## 💬 Community & Feedback
Have questions, suggestions or want to contribute ideas?  
Want to report issues, submit pull requests or browse the source code?  
Check out the [GitHub Repository](https://github.com/grulicht/terraform-provider-wazuh) for this provider.

## ✅ Daily End-to-End Testing
To ensure maximum reliability and functionality of this provider, **automated end-to-end tests are executed every day** via GitHub Actions.

These tests run against a real Wazuh instance (started using docker compose) and validate the majority of supported resources using real Terraform/OpenTofu plans and applies.

> 💡 This helps catch regressions early and ensures the provider remains fully operational and compatible with the Wazuh API.

## License
This module is 100% Open Source and is distributed under the MIT License.  
See the [LICENSE](https://github.com/grulicht/terraform-provider-wazuh/blob/main/LICENSE) file for more information.


## Acknowledgements
- HashiCorp Terraform
- [Wazuh](https://www.wazuh.com/)
- [OpenTofu](https://opentofu.org/)
- [Docker](https://www.docker.com/)
