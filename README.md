<p align="center">
  <a href="https://registry.terraform.io/providers/grulicht/wazuh/latest/docs">
    <img src="https://camo.githubusercontent.com/cdda8928975712cecce7be8b6a1506e3b327b1643cd3391dcf40515e25b54f73/68747470733a2f2f7777772e6461746f636d732d6173736574732e636f6d2f323838352f313733313337333331302d7465727261666f726d5f77686974652e737667" alt="Terraform Logo" width="200">
  </a>
  &nbsp;&nbsp;&nbsp;
  <a href="https://github.com/grulicht/terraform-provider-wazuh">
    <img src="https://wazuh.com/uploads/2022/05/Logo-blogpost.png" alt="terraform-provider-wazuh" width="200">
  </a>
  &nbsp;&nbsp;&nbsp;
  <a href="https://search.opentofu.org/provider/grulicht/wazuh/latest">
    <img src="https://raw.githubusercontent.com/opentofu/brand-artifacts/main/full/transparent/SVG/on-dark.svg#gh-dark-mode-only" alt="wazuh-provider-opentofu" width="200">
  </a>
  <h3 align="center" style="font-weight: bold">Terraform Provider for Wazuh</h3>
  <p align="center">
    <a href="https://github.com/grulicht/terraform-provider-wazuh/graphs/contributors">
      <img alt="Contributors" src="https://img.shields.io/github/contributors/grulicht/terraform-provider-wazuh">
    </a>
    <a href="https://golang.org/doc/devel/release.html">
      <img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/grulicht/terraform-provider-wazuh">
    </a>
    <a href="https://github.com/grulicht/terraform-provider-wazuh/actions?query=workflow%3Arelease">
      <img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/grulicht/terraform-provider-wazuh/release.yml?tag=latest&label=release">
    </a>
    <a href="https://github.com/grulicht/terraform-provider-wazuh/releases">
      <img alt="GitHub release (latest by date including pre-releases)" src="https://img.shields.io/github/v/release/grulicht/terraform-provider-wazuh?include_prereleases">
    </a>
  </p>
  <p align="center">
    <a href="https://github.com/grulicht/terraform-provider-wazuh/tree/main/docs"><strong>Explore the docs »</strong></a>
  </p>
</p>

# 🛡️ Wazuh Terraform/OpenTofu Provider
A Terraform/OpenTofu provider to manage [Wazuh](https://wazuh.com/) resources via its API using Terraform/OpenTofu.

It supports provisioning and configuration of Wazuh users and will be extended to support other objects such as hosts, templates, triggers, users etc.

## Requirements
- Go 1.21+ (if building from source)

## Building and Installing
```hcl
make build
```

## Provider Support
| Provider                                                                                   | Provider Support Status   |
|--------------------------------------------------------------------------------------------|---------------------------|
| [Terraform](https://registry.terraform.io/providers/grulicht/wazuh/latest)                 | ✅                        |
| [OpenTofu](https://search.opentofu.org/provider/grulicht/wazuh/latest)                     | ✅                        |


## ⚙️ **Example Provider Configuration**

```hcl
provider "wazuh" {
  endpoint        = "https://wazuh.example.com:55000"
  user            = "wazuh-wui"
  password        = "MyS3cr37P450r.*-"
  skip_ssl_verify = true
}
```

### 💡 Notes:

* The **default Wazuh API port** is `55000`.
* Authentication uses **JWT tokens**, automatically obtained by the provider via `/security/user/authenticate`.
* Token expiration defaults to 900 seconds (15 minutes). The provider will refresh it automatically in a future release.
* Use `skip_ssl_verify = true` only for local testing with self-signed certificates.

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

## Usage
See our [examples](./docs/resources/) per resources in docs.

| Resource                                   | Documentation                                                                                  | Example                                              | Status | Terraform Import / Create ⇒ Update | E2E Tests |
|--------------------------------------------|------------------------------------------------------------------------------------------------|------------------------------------------------------|--------|-------------------------------------|-----------|
| `wazuh_group`                              | [group.md](docs/resources/group.md)                                                            | [example](examples/group/)                           | ✅     | ✅ / ❌                             | ✅        |
| `wazuh_group_configuration`                | [group_configuration.md](docs/resources/group_configuration.md)                                | [example](examples/group_configuration/)             | ✅     | ✅ / ✅                             | ✅        |
| `wazuh_active_response`                    | [active_response.md](docs/resources/active_response.md)                                        | [example](examples/active_response/)                 | ✅     | ❌ / ❌                             | ✅        |
| `wazuh_event`                              | [event.md](docs/resources/event.md)                                                            | [example](examples/event/)                           | ✅     | ❌ / ❌                             | ✅        |
| `wazuh_cdb_list`                           | [cdb_list.md](docs/resources/cdb_list.md)                                                      | [example](examples/cdb_list/)                        | ✅     | ✅ / ✅                             | ✅        |
| `wazuh_decoder`                            | [decoder.md](docs/resources/decoder.md)                                                        | [example](examples/decoder/)                         | ✅     | ✅ / ✅                             | ✅        |
| `wazuh_node_configuration`                 | [node_configuration.md](docs/resources/node_configuration.md)                                  | [example](examples/node_configuration/)              | ✅     | ✅ / ✅                             | ✅        |
| `wazuh_node_restart`                       | [node_restart.md](docs/resources/node_restart.md)                                              | [example](examples/node_restart/)                    | ✅     | ❌ / ❌                             | ❌        |
| `wazuh_node_analysisd_reload`              | [node_analysisd_reload.md](docs/resources/node_analysisd_reload.md)                            | [example](examples/node_analysisd_reload/)           | ✅     | ❌ / ❌                             | ❌        |
| `wazuh_logtest`                            | [logtest.md](docs/resources/logtest.md)                                                        | [example](examples/logtest/)                         | ✅     | ❌ / ❌                             | ✅        |
| `wazuh_rootcheck`                          | [rootcheck.md](docs/resources/rootcheck.md)                                                    | [example](examples/rootcheck/)                       | ✅     | ✅ / ❌                             | ✅        |
| `wazuh_syscheck`                           | [syscheck.md](docs/resources/syscheck.md)                                                      | [example](examples/syscheck/)                        | ✅     | ✅ / ❌                             | ✅        |
| `wazuh_rule`                               | [rule.md](docs/resources/rule.md)                                                              | [example](examples/rule/)                            | ✅     | ✅ / ✅                             | ✅        |
| `wazuh_manager_configuration`              | [manager_configuration.md](docs/resources/manager_configuration.md)                            | [example](examples/manager_configuration/)           | ✅     | ✅ / ✅                             | ✅        |
| `wazuh_manager_restart`                    | [manager_restart.md](docs/resources/manager_restart.md)                                        | [example](examples/manager_restart/)                 | ✅     | ❌ / ❌                             | ✅        |
| `wazuh_agent`                              | [agent.md](docs/resources/agent.md)                                                            | [example](examples/agent/)                           | 📝     | ❌ / ❌                             | ❌        |
| `wazuh_agent_group`                        | [agent_group.md](docs/resources/agent_group.md)                                                | [example](examples/agent_group/)                     | 📝     | ❌ / ❌                             | ❌        |
| `wazuh_agent_restart`                      | [agent_restart.md](docs/resources/agent_restart.md)                                            | [example](examples/agent_restart/)                   | 📝     | ❌ / ❌                             | ❌        |
| `wazuh_agent_upgrade`                      | [agent_upgrade.md](docs/resources/agent_upgrade.md)                                            | [example](examples/agent_upgrade/)                   | 📝     | ❌ / ❌                             | ❌        |
| `wazuh_agent_upgrade_custom`               | [agent_upgrade_custom.md](docs/resources/agent_upgrade_custom.md)                              | [example](examples/agent_upgrade_custom/)            | 📝     | ❌ / ❌                             | ❌        |
| `wazuh_user`                               | [user.md](docs/resources/user.md)                                                              | [example](examples/user/)                            | 📝     | ❌ / ❌                             | ❌        |
| `wazuh_role`                               | [role.md](docs/resources/role.md)                                                              | [example](examples/role/)                            | 📝     | ❌ / ❌                             | ❌        |
| `wazuh_security_rule`                      | [security_rule.md](docs/resources/security_rule.md)                                            | [example](examples/security_rule/)                   | 📝     | ❌ / ❌                             | ❌        |
| `wazuh_policy`                             | [policy.md](docs/resources/policy.md)                                                          | [example](examples/policy/)                          | 📝     | ❌ / ❌                             | ❌        |
| `wazuh_role_user`                          | [role_user.md](docs/resources/role_user.md)                                                    | [example](examples/role_user/)                       | 📝     | ❌ / ❌                             | ❌        |
| `wazuh_security_rule_user`                 | [security_rule_user.md](docs/resources/security_rule_user.md)                                  | [example](examples/security_rule_user/)              | 📝     | ❌ / ❌                             | ❌        |
| `wazuh_policy_user`                        | [policy_user.md](docs/resources/policy_user.md)                                                | [example](examples/policy_user/)                     | 📝     | ❌ / ❌                             | ❌        |
| `wazuh_security_config`                    | [security_config.md](docs/resources/security_config.md)                                        | [example](examples/security_config/)                 | 📝     | ❌ / ❌                             | ❌        |


#### ℹ️ Note on Create ⇒ Update Behavior

Some resources support a "Create-or-Update" mechanism, when this behavior is implemented, it means:
> During the initial terraform apply, if an entity with the given name already exists, the resource will detect it and perform an update instead of attempting to create a duplicate => this is achieved by filtering existing entities by name before creation.
- This avoids the need for manual terraform import without having to have a terraform tfstate file or cleanup of existing resources in Wazuh.
- It's especially useful during migrations, initial setup, or when applying configuration into environments with pre-existing state.

---

### 💡 Missing a resource?
Is there a Wazuh resource you'd like to see supported?

👉 [Open an issue](https://github.com/grulicht/terraform-provider-wazuh/issues/new?template=feature_request.md) and we’ll consider it for implementation — or even better, submit a [Pull Request](https://github.com/grulicht/terraform-provider-wazuh/pulls) to contribute directly!

📘 See [CONTRIBUTING.md](./.github/CONTRIBUTING.md) for guidelines.

## 💬 Community & Feedback
Have questions, suggestions or want to contribute ideas?  
Want to report issues, submit pull requests or browse the source code?  
Check out the [GitHub Repository](https://github.com/grulicht/terraform-provider-wazuh) for this provider.

## ✅ Daily End-to-End Testing
To ensure maximum reliability and functionality of this provider, **automated end-to-end tests are executed every day** via GitHub Actions.

These tests run against a real Wazuh instance (started using docker compose) and validate the majority of supported resources using real Terraform/OpenTofu plans and applies.

> 💡 This helps catch regressions early and ensures the provider remains fully operational and compatible with the Wazuh API.

### 🔄 Workflows
The project uses GitHub Actions to automate validation and testing of the provider.

- Validate and lint documentation files (`README.md` and `docs/`)
- Initialize, test and check the Wazuh provider with **Terraform** and **OpenTofu**
- Publish the new version of the Wazuh Terraform provider to Terraform Registry
- Run daily **E2E Terraform tests** against a live Wazuh instance spun up via Docker Compose (`make up`) at **07:00 UTC**

### 🧪 Localy Testing
To test the provider locally, start the Wazuh Web UI using Docker Compose:
```sh
make up
```
Then open `https://localhost:443` in your browser manually or by command:
```sh
make launch
```

### 🔐 Predefined Test Credentials for Login (use also E2E tests)

#### GUI
| **Field**    | **Value**                                                                  |
|--------------|----------------------------------------------------------------------------|
| Username     | `admin`                                                                    |
| Password     | `SecretPassword`                                                           |
| URL          | `https://localhost:443`                                                    |

#### API
| **Field**    | **Value**                                                                  |
|--------------|----------------------------------------------------------------------------|
| Username     | `wazuh-wui`                                                                |
| Password     | `MyS3cr37P450r.*-`                                                         |
| URL          | `https://localhost:55000`                                                  |

> [Docs for change password of default users for Wazuh.](https://documentation.wazuh.com/current/deployment-options/docker/changing-default-password.html)
>
> You can now apply your Terraform/OpenTofu templates and observe changes live in the UI.

### Testing a new version of the Wazuh provider
After making changes to the provider source code, follow these steps:
Build the provider binary:
```sh
make build
```
Install the binary into the local Terraform/OpenTofu plugin directory:
```sh
make install-plugin
```
Update your main.tf to use the local provider source
Add the following to your Terraform/OpenTofu configuration:
```sh
terraform {
  required_providers {
    wazuh = {
      source  = "localdomain/local/wazuh"
    }
  }
}
```
Now you're ready to test your provider against the local Wazuh instance.

## Roadmap
See the [open issues](https://github.com/grulicht/terraform-provider-wazuh/issues) for a list of proposed features (and known issues). See [CONTRIBUTING](./.github/CONTRIBUTING.md) for more information.

## License
This module is 100% Open Source and is distributed under the MIT License.  
See the [LICENSE](https://github.com/grulicht/terraform-provider-wazuh/blob/main/LICENSE) file for more information.


## Acknowledgements
- HashiCorp Terrafor
- [Wazuh](https://wazuh.com/)
- [OpenTofu](https://opentofu.org/)
- [Docker](https://www.docker.com/)
