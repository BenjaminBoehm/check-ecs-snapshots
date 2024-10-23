# Description
Check for Elastic Cloud Server (ECS) snapshot age on Open Telekom Cloud (OTC)

# Environment

* Option 1: clouds.yaml
  * create config file (~/.config/openstack/clouds.yaml or /etc/openstack/clouds.yaml)
    ```yaml
    ---
    clouds:
      myCloud:
        auth:
          username: '<USERNAME>'
          password: '<PASSWORD>'
          project_name: 'eu-de'
          user_domain_name: '<DOMAIN_NAME>'
          auth_url: 'https://iam.eu-de.otc.t-systems.com/v3'
        interface: 'public'
        identity_api_version: 3 # !Important
    ```
  * set env which cloud credentials should be used
    ```bash
    export OS_CLOUD=myCloud
    ```

* Option 2: Environment Variables
    ```bash
    export OS_PROJECT_NANME=eu-de
    export OS_AUTH_URL="https://iam.eu-de.otc.t-systems.com:443/v3"
    export OS_USER_DOMAIN_NAME="<DOMAIN_NAME>"
    export OS_USERNAME="<USERNAME>"
    export OS_PASSWORD="<PASSWORD>"
    ```

# Run
```bash
go run cmd/main.go
```

# Build
```bash
make build
```

# Dependencies

* OTC Cloud SDK: https://github.com/opentelekomcloud/gophertelekomcloud/
* Netways monitoring plugins library: https://github.com/NETWAYS/go-check
