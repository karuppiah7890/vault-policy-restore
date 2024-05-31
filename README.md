# vault-policy-restore

Using this CLI tool, you can restore Vault Policies to a Vault instance! :D

`vault-policy-restore` is a CLI tool that can be used to restore a previously backed up set of Vault policies. The backup could be done using a tool like [`vault-policy-backup`](https://github.com/karuppiah7890/vault-policy-backup)

Note: The tool is written in Golang and uses Vault Official Golang API. The official Vault API documentation is here - https://pkg.go.dev/github.com/hashicorp/vault/api

Note: The tool needs Vault credentials of a user/account that has access to Vault, to create/write the Vault Policies

Note: We have tested this only with some versions of Vault (like v1.15.x). So beware to test this in a testing environment with whatever version of Vault you are using, before using this in critical environments like production! Also, ensure that the testing environment is as close to your production environment as possible so that your testing makes sense

Note ⚠️‼️🚨: If the Vault instance has some policies already defined with the same name as the Policies present in the Vault backup JSON file, when restoring to the Vault instance using the Vault backup JSON file, the Policies in the Vault instance will be overwritten! All the Vault Policies in Vault backup JSON file will be present in the Vault instance. If the Vault instance has some extra Vault Policies configured, it might have those untouched and intact

Note: This does NOT restore the `root` Vault Policy in case one is present in the Vault Policies backup JSON file - this is because Vault does not support updating / changing the root policy - which would happen during the restore process. The `root` policy can neither be deleted, created or updated (or changed). Also, `root` Vault Policy is just an empty policy, with no content. I believe it's just a placeholder policy which is assumed to have all the access to Vault

# Building

```bash
CGO_ENABLED=0 go build -v
```

or

```bash
make
```

# Usage

```bash
$ ./vault-policy-restore --help
usage: vault-policy-restore [-quiet|--quiet] [-file|-file <vault-policy-backup-json-file-path>]

Usage of ./vault-policy-restore:

  -file / --file string
      vault policy backup json file path (default "vault_policy_backup.json")

  -quiet / --quiet
      quiet progress (default false).
      By default vault-policy-restore CLI will show a lot of details
      about the restore process and detailed progress during the
      restore process

  -h / -help / --help
      show help

examples:

# show help
vault-policy-restore -h

# restores all vault policies from the JSON file
# except the root policy if it's present in the JSON file.
# also, any existing vault policies with the same name as
# the policy name in the JSON file will be overwritten.
vault-policy-restore -file <path-to-vault-policy-backup-json-file>

# OR you can use --file too instead of -file

vault-policy-restore --file <path-to-vault-policy-backup-json-file>

# quietly restore all vault policies.
# this will just show dots (.) for progress
vault-policy-restore -quiet -file <path-to-vault-policy-backup-json-file>

# OR you can use --quiet too instead of -quiet

vault-policy-restore --quiet --file <path-to-vault-policy-backup-json-file>
```

# How `vault-policy-restore` works

`vault-policy-restore` expects the policies backup file to be a JSON file with content structure being similar to the ones in this example structure -

```json
{
  "policies": [
    {
      "name": "allow_secrets",
      "policy": "path \"secret/*\" {\n  capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"]\n}\n"
    }
  ]
}
```

So, the JSON structure looks something like -

```json
{
  "policies": [
    {
      "name": "<policy1-name>",
      "policy": "<policy1-content-in-hcl-format>"
    },
    {
      "name": "<policy2-name>",
      "policy": "<policy2-content-in-hcl-format>"
    },
    {
      "name": "<policy3-name>",
      "policy": "<policy3-content-in-hcl-format>"
    }
  ]
}
```

That's just three policies, but the policies backup JSON file can contain any number of Vault Policies :)

Once you have Vault Policies Backup JSON file, you just need to pass it to `vault-policy-restore` tool and it will restore it for you given a Vault instance and it's details like connectivity details etc - Hostname/IP, Port, and Vault Token with just enough access to write the Vault Policies to create and/ update Vault Policies in the Vault instance. So, that's all you need to do!

# Demo

## Demo 1

In my local machine, I have a sample `vault_policy_backup.json` file like this -

```bash
$ cat vault_policy_backup.json
{"policies":[{"name":"allow_secrets","policy":"path \"secret/*\" {\n  capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"]\n}\n"},{"name":"allow_stage_kv_secrets","policy":"# KV v2 secrets engine mount path is \"stage-kv\"\npath \"stage-kv/*\" {\n  capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"]\n}\n"},{"name":"allow_test_kv_secrets","policy":"# KV v2 secrets engine mount path is \"test-kv\"\npath \"test-kv/*\" {\n  capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"]\n}\n"},{"name":"default","policy":"\n# Allow tokens to look up their own properties\npath \"auth/token/lookup-self\" {\n    capabilities = [\"read\"]\n}\n\n# Allow tokens to renew themselves\npath \"auth/token/renew-self\" {\n    capabilities = [\"update\"]\n}\n\n# Allow tokens to revoke themselves\npath \"auth/token/revoke-self\" {\n    capabilities = [\"update\"]\n}\n\n# Allow a token to look up its own capabilities on a path\npath \"sys/capabilities-self\" {\n    capabilities = [\"update\"]\n}\n\n# Allow a token to look up its own entity by id or name\npath \"identity/entity/id/{{identity.entity.id}}\" {\n  capabilities = [\"read\"]\n}\npath \"identity/entity/name/{{identity.entity.name}}\" {\n  capabilities = [\"read\"]\n}\n\n\n# Allow a token to look up its resultant ACL from all policies. This is useful\n# for UIs. It is an internal path because the format may change at any time\n# based on how the internal ACL features and capabilities change.\npath \"sys/internal/ui/resultant-acl\" {\n    capabilities = [\"read\"]\n}\n\n# Allow a token to renew a lease via lease_id in the request body; old path for\n# old clients, new path for newer\npath \"sys/renew\" {\n    capabilities = [\"update\"]\n}\npath \"sys/leases/renew\" {\n    capabilities = [\"update\"]\n}\n\n# Allow looking up lease properties. This requires knowing the lease ID ahead\n# of time and does not divulge any sensitive information.\npath \"sys/leases/lookup\" {\n    capabilities = [\"update\"]\n}\n\n# Allow a token to manage its own cubbyhole\npath \"cubbyhole/*\" {\n    capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"]\n}\n\n# Allow a token to wrap arbitrary values in a response-wrapping token\npath \"sys/wrapping/wrap\" {\n    capabilities = [\"update\"]\n}\n\n# Allow a token to look up the creation time and TTL of a given\n# response-wrapping token\npath \"sys/wrapping/lookup\" {\n    capabilities = [\"update\"]\n}\n\n# Allow a token to unwrap a response-wrapping token. This is a convenience to\n# avoid client token swapping since this is also part of the response wrapping\n# policy.\npath \"sys/wrapping/unwrap\" {\n    capabilities = [\"update\"]\n}\n\n# Allow general purpose tools\npath \"sys/tools/hash\" {\n    capabilities = [\"update\"]\n}\npath \"sys/tools/hash/*\" {\n    capabilities = [\"update\"]\n}\n\n# Allow checking the status of a Control Group request if the user has the\n# accessor\npath \"sys/control-group/request\" {\n    capabilities = [\"update\"]\n}\n\n# Allow a token to make requests to the Authorization Endpoint for OIDC providers.\npath \"identity/oidc/provider/+/authorize\" {\n    capabilities = [\"read\", \"update\"]\n}\n"}]}

$ cat vault_policy_backup.json | jq

{
  "policies": [
    {
      "name": "allow_secrets",
      "policy": "path \"secret/*\" {\n  capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"]\n}\n"
    },
    {
      "name": "allow_stage_kv_secrets",
      "policy": "# KV v2 secrets engine mount path is \"stage-kv\"\npath \"stage-kv/*\" {\n  capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"]\n}\n"
    },
    {
      "name": "allow_test_kv_secrets",
      "policy": "# KV v2 secrets engine mount path is \"test-kv\"\npath \"test-kv/*\" {\n  capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"]\n}\n"
    },
    {
      "name": "default",
      "policy": "\n# Allow tokens to look up their own properties\npath \"auth/token/lookup-self\" {\n    capabilities = [\"read\"]\n}\n\n# Allow tokens to renew themselves\npath \"auth/token/renew-self\" {\n    capabilities = [\"update\"]\n}\n\n# Allow tokens to revoke themselves\npath \"auth/token/revoke-self\" {\n    capabilities = [\"update\"]\n}\n\n# Allow a token to look up its own capabilities on a path\npath \"sys/capabilities-self\" {\n    capabilities = [\"update\"]\n}\n\n# Allow a token to look up its own entity by id or name\npath \"identity/entity/id/{{identity.entity.id}}\" {\n  capabilities = [\"read\"]\n}\npath \"identity/entity/name/{{identity.entity.name}}\" {\n  capabilities = [\"read\"]\n}\n\n\n# Allow a token to look up its resultant ACL from all policies. This is useful\n# for UIs. It is an internal path because the format may change at any time\n# based on how the internal ACL features and capabilities change.\npath \"sys/internal/ui/resultant-acl\" {\n    capabilities = [\"read\"]\n}\n\n# Allow a token to renew a lease via lease_id in the request body; old path for\n# old clients, new path for newer\npath \"sys/renew\" {\n    capabilities = [\"update\"]\n}\npath \"sys/leases/renew\" {\n    capabilities = [\"update\"]\n}\n\n# Allow looking up lease properties. This requires knowing the lease ID ahead\n# of time and does not divulge any sensitive information.\npath \"sys/leases/lookup\" {\n    capabilities = [\"update\"]\n}\n\n# Allow a token to manage its own cubbyhole\npath \"cubbyhole/*\" {\n    capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"]\n}\n\n# Allow a token to wrap arbitrary values in a response-wrapping token\npath \"sys/wrapping/wrap\" {\n    capabilities = [\"update\"]\n}\n\n# Allow a token to look up the creation time and TTL of a given\n# response-wrapping token\npath \"sys/wrapping/lookup\" {\n    capabilities = [\"update\"]\n}\n\n# Allow a token to unwrap a response-wrapping token. This is a convenience to\n# avoid client token swapping since this is also part of the response wrapping\n# policy.\npath \"sys/wrapping/unwrap\" {\n    capabilities = [\"update\"]\n}\n\n# Allow general purpose tools\npath \"sys/tools/hash\" {\n    capabilities = [\"update\"]\n}\npath \"sys/tools/hash/*\" {\n    capabilities = [\"update\"]\n}\n\n# Allow checking the status of a Control Group request if the user has the\n# accessor\npath \"sys/control-group/request\" {\n    capabilities = [\"update\"]\n}\n\n# Allow a token to make requests to the Authorization Endpoint for OIDC providers.\npath \"identity/oidc/provider/+/authorize\" {\n    capabilities = [\"read\", \"update\"]\n}\n"
    }
  ]
}
```

Now, let's run a simple brand new Vault instance in my local - say, a Dev Server of Vault and try out the `vault-policy-restore` command :) :D

Let's run a simple Dev Server of Vault with this command in my local -

```bash
vault server -dev -dev-root-token-id root -dev-listen-address 127.0.0.1:8200
```

Let's check if the Vault Dev Server is working and is all good

```bash
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ export VAULT_TOKEN="root"

$ vault status
Key             Value
---             -----
Seal Type       shamir
Initialized     true
Sealed          false
Total Shares    1
Threshold       1
Version         1.15.4
Build Date      2023-12-04T17:45:28Z
Storage Type    inmem
Cluster Name    vault-cluster-ab30bce7
Cluster ID      20464c56-def0-e2df-537b-13036cef40f0
HA Enabled      false
```

Awesome :D So, it's up and running now :D And it has some Vault Policies already present in it. Let's delete the ones we can delete to keep it as a clean slate

```bash

$ vault status
Key             Value
---             -----
Seal Type       shamir
Initialized     true
Sealed          false
Total Shares    1
Threshold       1
Version         1.15.4
Build Date      2023-12-04T17:45:28Z
Storage Type    inmem
Cluster Name    vault-cluster-dde6e6b4
Cluster ID      6f0a7185-3327-6abb-3ef3-5965929d7299
HA Enabled      false

$ vault policy list
default
root

$ vault policy read default
# Allow tokens to look up their own properties
path "auth/token/lookup-self" {
    capabilities = ["read"]
}

# Allow tokens to renew themselves
path "auth/token/renew-self" {
    capabilities = ["update"]
}

# Allow tokens to revoke themselves
path "auth/token/revoke-self" {
    capabilities = ["update"]
}

# Allow a token to look up its own capabilities on a path
path "sys/capabilities-self" {
    capabilities = ["update"]
}

# Allow a token to look up its own entity by id or name
path "identity/entity/id/{{identity.entity.id}}" {
  capabilities = ["read"]
}
path "identity/entity/name/{{identity.entity.name}}" {
  capabilities = ["read"]
}


# Allow a token to look up its resultant ACL from all policies. This is useful
# for UIs. It is an internal path because the format may change at any time
# based on how the internal ACL features and capabilities change.
path "sys/internal/ui/resultant-acl" {
    capabilities = ["read"]
}

# Allow a token to renew a lease via lease_id in the request body; old path for
# old clients, new path for newer
path "sys/renew" {
    capabilities = ["update"]
}
path "sys/leases/renew" {
    capabilities = ["update"]
}

# Allow looking up lease properties. This requires knowing the lease ID ahead
# of time and does not divulge any sensitive information.
path "sys/leases/lookup" {
    capabilities = ["update"]
}

# Allow a token to manage its own cubbyhole
path "cubbyhole/*" {
    capabilities = ["create", "read", "update", "delete", "list"]
}

# Allow a token to wrap arbitrary values in a response-wrapping token
path "sys/wrapping/wrap" {
    capabilities = ["update"]
}

# Allow a token to look up the creation time and TTL of a given
# response-wrapping token
path "sys/wrapping/lookup" {
    capabilities = ["update"]
}

# Allow a token to unwrap a response-wrapping token. This is a convenience to
# avoid client token swapping since this is also part of the response wrapping
# policy.
path "sys/wrapping/unwrap" {
    capabilities = ["update"]
}

# Allow general purpose tools
path "sys/tools/hash" {
    capabilities = ["update"]
}
path "sys/tools/hash/*" {
    capabilities = ["update"]
}

# Allow checking the status of a Control Group request if the user has the
# accessor
path "sys/control-group/request" {
    capabilities = ["update"]
}

# Allow a token to make requests to the Authorization Endpoint for OIDC providers.
path "identity/oidc/provider/+/authorize" {
    capabilities = ["read", "update"]
}

$ vault policy delete default
Error deleting default: Error making API request.

URL: DELETE http://127.0.0.1:8200/v1/sys/policies/acl/default
Code: 400. Errors:

* cannot delete default policy

$ vault policy read root
No policy named: root

$ vault read sys/policies/acl/root
Key       Value
---       -----
name      root
policy    n/a

$ vault policy delete root
Error deleting root: Error making API request.

URL: DELETE http://127.0.0.1:8200/v1/sys/policies/acl/root
Code: 400. Errors:

* cannot delete "root" policy
```

We just have `root` and `default` policies and they can't be deleted. So, let's move on

Note: By the way, as mentioned before `root` policy cannot be updated or changed too, and as we can see from the output above, `root` policy cannot be deleted too. But - `default` policy can be updated or changed. We'll see that in a moment as the Vault Policy Backup JSON file has `default` policy defined in it ;) :)

Let's run the `vault-policy-restore` command now :D

```bash
$ ./vault-policy-restore --file vault_policy_backup.json 

restoring 4 vault policies in vault

restoring the following vault policies in vault: [allow_secrets allow_stage_kv_secrets allow_test_kv_secrets default]

restoring `allow_secrets` vault policy

`allow_secrets` vault policy rules: path "secret/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}


restoring `allow_stage_kv_secrets` vault policy

`allow_stage_kv_secrets` vault policy rules: # KV v2 secrets engine mount path is "stage-kv"
path "stage-kv/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}


restoring `allow_test_kv_secrets` vault policy

`allow_test_kv_secrets` vault policy rules: # KV v2 secrets engine mount path is "test-kv"
path "test-kv/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}


restoring `default` vault policy

`default` vault policy rules: 
# Allow tokens to look up their own properties
path "auth/token/lookup-self" {
    capabilities = ["read"]
}

# Allow tokens to renew themselves
path "auth/token/renew-self" {
    capabilities = ["update"]
}

# Allow tokens to revoke themselves
path "auth/token/revoke-self" {
    capabilities = ["update"]
}

# Allow a token to look up its own capabilities on a path
path "sys/capabilities-self" {
    capabilities = ["update"]
}

# Allow a token to look up its own entity by id or name
path "identity/entity/id/{{identity.entity.id}}" {
  capabilities = ["read"]
}
path "identity/entity/name/{{identity.entity.name}}" {
  capabilities = ["read"]
}


# Allow a token to look up its resultant ACL from all policies. This is useful
# for UIs. It is an internal path because the format may change at any time
# based on how the internal ACL features and capabilities change.
path "sys/internal/ui/resultant-acl" {
    capabilities = ["read"]
}

# Allow a token to renew a lease via lease_id in the request body; old path for
# old clients, new path for newer
path "sys/renew" {
    capabilities = ["update"]
}
path "sys/leases/renew" {
    capabilities = ["update"]
}

# Allow looking up lease properties. This requires knowing the lease ID ahead
# of time and does not divulge any sensitive information.
path "sys/leases/lookup" {
    capabilities = ["update"]
}

# Allow a token to manage its own cubbyhole
path "cubbyhole/*" {
    capabilities = ["create", "read", "update", "delete", "list"]
}

# Allow a token to wrap arbitrary values in a response-wrapping token
path "sys/wrapping/wrap" {
    capabilities = ["update"]
}

# Allow a token to look up the creation time and TTL of a given
# response-wrapping token
path "sys/wrapping/lookup" {
    capabilities = ["update"]
}

# Allow a token to unwrap a response-wrapping token. This is a convenience to
# avoid client token swapping since this is also part of the response wrapping
# policy.
path "sys/wrapping/unwrap" {
    capabilities = ["update"]
}

# Allow general purpose tools
path "sys/tools/hash" {
    capabilities = ["update"]
}
path "sys/tools/hash/*" {
    capabilities = ["update"]
}

# Allow checking the status of a Control Group request if the user has the
# accessor
path "sys/control-group/request" {
    capabilities = ["update"]
}

# Allow a token to make requests to the Authorization Endpoint for OIDC providers.
path "identity/oidc/provider/+/authorize" {
    capabilities = ["read", "update"]
}
```

Now, let's list and read the policies :)

```bash
$ vault policy list
allow_secrets
allow_stage_kv_secrets
allow_test_kv_secrets
default
root

$ vault policy read allow_secrets
path "secret/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}

$ vault policy read allow_stage_kv_secrets
# KV v2 secrets engine mount path is "stage-kv"
path "stage-kv/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}

$ vault policy read allow_test_kv_secrets
# KV v2 secrets engine mount path is "test-kv"
path "test-kv/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}

$ vault policy read default
# Allow tokens to look up their own properties
path "auth/token/lookup-self" {
    capabilities = ["read"]
}

# Allow tokens to renew themselves
path "auth/token/renew-self" {
    capabilities = ["update"]
}

# Allow tokens to revoke themselves
path "auth/token/revoke-self" {
    capabilities = ["update"]
}

# Allow a token to look up its own capabilities on a path
path "sys/capabilities-self" {
    capabilities = ["update"]
}

# Allow a token to look up its own entity by id or name
path "identity/entity/id/{{identity.entity.id}}" {
  capabilities = ["read"]
}
path "identity/entity/name/{{identity.entity.name}}" {
  capabilities = ["read"]
}


# Allow a token to look up its resultant ACL from all policies. This is useful
# for UIs. It is an internal path because the format may change at any time
# based on how the internal ACL features and capabilities change.
path "sys/internal/ui/resultant-acl" {
    capabilities = ["read"]
}

# Allow a token to renew a lease via lease_id in the request body; old path for
# old clients, new path for newer
path "sys/renew" {
    capabilities = ["update"]
}
path "sys/leases/renew" {
    capabilities = ["update"]
}

# Allow looking up lease properties. This requires knowing the lease ID ahead
# of time and does not divulge any sensitive information.
path "sys/leases/lookup" {
    capabilities = ["update"]
}

# Allow a token to manage its own cubbyhole
path "cubbyhole/*" {
    capabilities = ["create", "read", "update", "delete", "list"]
}

# Allow a token to wrap arbitrary values in a response-wrapping token
path "sys/wrapping/wrap" {
    capabilities = ["update"]
}

# Allow a token to look up the creation time and TTL of a given
# response-wrapping token
path "sys/wrapping/lookup" {
    capabilities = ["update"]
}

# Allow a token to unwrap a response-wrapping token. This is a convenience to
# avoid client token swapping since this is also part of the response wrapping
# policy.
path "sys/wrapping/unwrap" {
    capabilities = ["update"]
}

# Allow general purpose tools
path "sys/tools/hash" {
    capabilities = ["update"]
}
path "sys/tools/hash/*" {
    capabilities = ["update"]
}

# Allow checking the status of a Control Group request if the user has the
# accessor
path "sys/control-group/request" {
    capabilities = ["update"]
}

# Allow a token to make requests to the Authorization Endpoint for OIDC providers.
path "identity/oidc/provider/+/authorize" {
    capabilities = ["read", "update"]
}
```

As you can see, all the Vault Policies that were backed up are restored :D `allow_secrets`, `allow_stage_kv_secrets`, `allow_test_kv_secrets` and `default` policies were restored.

Yes, `default` policy was also restored - more like - overwritten :D

You can also run the `vault-policy-restore` command with the `-quiet` / `--quiet` flag to hide too much details and just show quiet progress

```bash
$ ./vault-policy-restore --quiet --file vault_policy_backup.json 

restoring 4 vault policies in vault

restoring the following vault policies in vault: [allow_secrets allow_stage_kv_secrets allow_test_kv_secrets default]
....

```

Now, let's list and read the policies again! :) To see if all is still good, after a the second restore with the same `vault_policy_backup.json` file

```bash
$ vault policy list
allow_secrets
allow_stage_kv_secrets
allow_test_kv_secrets
default
root

$ vault policy read allow_secrets
path "secret/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}

$ vault policy read allow_stage_kv_secrets
# KV v2 secrets engine mount path is "stage-kv"
path "stage-kv/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}

$ vault policy read allow_test_kv_secrets
# KV v2 secrets engine mount path is "test-kv"
path "test-kv/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}

$ vault policy read default
# Allow tokens to look up their own properties
path "auth/token/lookup-self" {
    capabilities = ["read"]
}

# Allow tokens to renew themselves
path "auth/token/renew-self" {
    capabilities = ["update"]
}

# Allow tokens to revoke themselves
path "auth/token/revoke-self" {
    capabilities = ["update"]
}

# Allow a token to look up its own capabilities on a path
path "sys/capabilities-self" {
    capabilities = ["update"]
}

# Allow a token to look up its own entity by id or name
path "identity/entity/id/{{identity.entity.id}}" {
  capabilities = ["read"]
}
path "identity/entity/name/{{identity.entity.name}}" {
  capabilities = ["read"]
}


# Allow a token to look up its resultant ACL from all policies. This is useful
# for UIs. It is an internal path because the format may change at any time
# based on how the internal ACL features and capabilities change.
path "sys/internal/ui/resultant-acl" {
    capabilities = ["read"]
}

# Allow a token to renew a lease via lease_id in the request body; old path for
# old clients, new path for newer
path "sys/renew" {
    capabilities = ["update"]
}
path "sys/leases/renew" {
    capabilities = ["update"]
}

# Allow looking up lease properties. This requires knowing the lease ID ahead
# of time and does not divulge any sensitive information.
path "sys/leases/lookup" {
    capabilities = ["update"]
}

# Allow a token to manage its own cubbyhole
path "cubbyhole/*" {
    capabilities = ["create", "read", "update", "delete", "list"]
}

# Allow a token to wrap arbitrary values in a response-wrapping token
path "sys/wrapping/wrap" {
    capabilities = ["update"]
}

# Allow a token to look up the creation time and TTL of a given
# response-wrapping token
path "sys/wrapping/lookup" {
    capabilities = ["update"]
}

# Allow a token to unwrap a response-wrapping token. This is a convenience to
# avoid client token swapping since this is also part of the response wrapping
# policy.
path "sys/wrapping/unwrap" {
    capabilities = ["update"]
}

# Allow general purpose tools
path "sys/tools/hash" {
    capabilities = ["update"]
}
path "sys/tools/hash/*" {
    capabilities = ["update"]
}

# Allow checking the status of a Control Group request if the user has the
# accessor
path "sys/control-group/request" {
    capabilities = ["update"]
}

# Allow a token to make requests to the Authorization Endpoint for OIDC providers.
path "identity/oidc/provider/+/authorize" {
    capabilities = ["read", "update"]
}
```

## Demo 2

Let's see how overwriting of Vault Policies happen. Let me get a brand new Vault server first :) after killing the old Vault server

```bash
vault server -dev -dev-root-token-id root -dev-listen-address 127.0.0.1:8200
```

Next, let me create some policies in Vault. This time, I'm creating a policy that allows me to read and list secrets in secrets engine "secret"

```bash
$ cat /Users/karuppiah.n/every-day-log/allow_read_secrets.hcl 
# KV v2 secrets engine mount path is "secret"
path "secret/*" {
  capabilities = ["read", "list"]
}

$ vault policy write allow_secrets /Users/karuppiah.n/every-day-log/allow_read_secrets.hcl 
Success! Uploaded policy: allow_secrets

$ vault policy read allow_secrets
# KV v2 secrets engine mount path is "secret"
path "secret/*" {
  capabilities = ["read", "list"]
}
```

Notice how the policy name is still `allow_secrets` and not `allow_read_secrets`. This will give us an opportunity to now overrwrite the `allow_secrets` policy with a new policy from our backed up Vault Policies JSON file that we have already used once in our first demo

Let's look the Vault Policies Backup JSON file that we have -

```bash
$ cat vault_policy_backup.json | jq
{"policies":[{"name":"allow_secrets","policy":"path \"secret/*\" {\n  capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"]\n}\n"},{"name":"allow_stage_kv_secrets","policy":"# KV v2 secrets engine mount path is \"stage-kv\"\npath \"stage-kv/*\" {\n  capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"]\n}\n"},{"name":"allow_test_kv_secrets","policy":"# KV v2 secrets engine mount path is \"test-kv\"\npath \"test-kv/*\" {\n  capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"]\n}\n"},{"name":"default","policy":"\n# Allow tokens to look up their own properties\npath \"auth/token/lookup-self\" {\n    capabilities = [\"read\"]\n}\n\n# Allow tokens to renew themselves\npath \"auth/token/renew-self\" {\n    capabilities = [\"update\"]\n}\n\n# Allow tokens to revoke themselves\npath \"auth/token/revoke-self\" {\n    capabilities = [\"update\"]\n}\n\n# Allow a token to look up its own capabilities on a path\npath \"sys/capabilities-self\" {\n    capabilities = [\"update\"]\n}\n\n# Allow a token to look up its own entity by id or name\npath \"identity/entity/id/{{identity.entity.id}}\" {\n  capabilities = [\"read\"]\n}\npath \"identity/entity/name/{{identity.entity.name}}\" {\n  capabilities = [\"read\"]\n}\n\n\n# Allow a token to look up its resultant ACL from all policies. This is useful\n# for UIs. It is an internal path because the format may change at any time\n# based on how the internal ACL features and capabilities change.\npath \"sys/internal/ui/resultant-acl\" {\n    capabilities = [\"read\"]\n}\n\n# Allow a token to renew a lease via lease_id in the request body; old path for\n# old clients, new path for newer\npath \"sys/renew\" {\n    capabilities = [\"update\"]\n}\npath \"sys/leases/renew\" {\n    capabilities = [\"update\"]\n}\n\n# Allow looking up lease properties. This requires knowing the lease ID ahead\n# of time and does not divulge any sensitive information.\npath \"sys/leases/lookup\" {\n    capabilities = [\"update\"]\n}\n\n# Allow a token to manage its own cubbyhole\npath \"cubbyhole/*\" {\n    capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"]\n}\n\n# Allow a token to wrap arbitrary values in a response-wrapping token\npath \"sys/wrapping/wrap\" {\n    capabilities = [\"update\"]\n}\n\n# Allow a token to look up the creation time and TTL of a given\n# response-wrapping token\npath \"sys/wrapping/lookup\" {\n    capabilities = [\"update\"]\n}\n\n# Allow a token to unwrap a response-wrapping token. This is a convenience to\n# avoid client token swapping since this is also part of the response wrapping\n# policy.\npath \"sys/wrapping/unwrap\" {\n    capabilities = [\"update\"]\n}\n\n# Allow general purpose tools\npath \"sys/tools/hash\" {\n    capabilities = [\"update\"]\n}\npath \"sys/tools/hash/*\" {\n    capabilities = [\"update\"]\n}\n\n# Allow checking the status of a Control Group request if the user has the\n# accessor\npath \"sys/control-group/request\" {\n    capabilities = [\"update\"]\n}\n\n# Allow a token to make requests to the Authorization Endpoint for OIDC providers.\npath \"identity/oidc/provider/+/authorize\" {\n    capabilities = [\"read\", \"update\"]\n}\n"}]}

$ cat vault_policy_backup.json | jq

{
  "policies": [
    {
      "name": "allow_secrets",
      "policy": "path \"secret/*\" {\n  capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"]\n}\n"
    },
    {
      "name": "allow_stage_kv_secrets",
      "policy": "# KV v2 secrets engine mount path is \"stage-kv\"\npath \"stage-kv/*\" {\n  capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"]\n}\n"
    },
    {
      "name": "allow_test_kv_secrets",
      "policy": "# KV v2 secrets engine mount path is \"test-kv\"\npath \"test-kv/*\" {\n  capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"]\n}\n"
    },
    {
      "name": "default",
      "policy": "\n# Allow tokens to look up their own properties\npath \"auth/token/lookup-self\" {\n    capabilities = [\"read\"]\n}\n\n# Allow tokens to renew themselves\npath \"auth/token/renew-self\" {\n    capabilities = [\"update\"]\n}\n\n# Allow tokens to revoke themselves\npath \"auth/token/revoke-self\" {\n    capabilities = [\"update\"]\n}\n\n# Allow a token to look up its own capabilities on a path\npath \"sys/capabilities-self\" {\n    capabilities = [\"update\"]\n}\n\n# Allow a token to look up its own entity by id or name\npath \"identity/entity/id/{{identity.entity.id}}\" {\n  capabilities = [\"read\"]\n}\npath \"identity/entity/name/{{identity.entity.name}}\" {\n  capabilities = [\"read\"]\n}\n\n\n# Allow a token to look up its resultant ACL from all policies. This is useful\n# for UIs. It is an internal path because the format may change at any time\n# based on how the internal ACL features and capabilities change.\npath \"sys/internal/ui/resultant-acl\" {\n    capabilities = [\"read\"]\n}\n\n# Allow a token to renew a lease via lease_id in the request body; old path for\n# old clients, new path for newer\npath \"sys/renew\" {\n    capabilities = [\"update\"]\n}\npath \"sys/leases/renew\" {\n    capabilities = [\"update\"]\n}\n\n# Allow looking up lease properties. This requires knowing the lease ID ahead\n# of time and does not divulge any sensitive information.\npath \"sys/leases/lookup\" {\n    capabilities = [\"update\"]\n}\n\n# Allow a token to manage its own cubbyhole\npath \"cubbyhole/*\" {\n    capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"]\n}\n\n# Allow a token to wrap arbitrary values in a response-wrapping token\npath \"sys/wrapping/wrap\" {\n    capabilities = [\"update\"]\n}\n\n# Allow a token to look up the creation time and TTL of a given\n# response-wrapping token\npath \"sys/wrapping/lookup\" {\n    capabilities = [\"update\"]\n}\n\n# Allow a token to unwrap a response-wrapping token. This is a convenience to\n# avoid client token swapping since this is also part of the response wrapping\n# policy.\npath \"sys/wrapping/unwrap\" {\n    capabilities = [\"update\"]\n}\n\n# Allow general purpose tools\npath \"sys/tools/hash\" {\n    capabilities = [\"update\"]\n}\npath \"sys/tools/hash/*\" {\n    capabilities = [\"update\"]\n}\n\n# Allow checking the status of a Control Group request if the user has the\n# accessor\npath \"sys/control-group/request\" {\n    capabilities = [\"update\"]\n}\n\n# Allow a token to make requests to the Authorization Endpoint for OIDC providers.\npath \"identity/oidc/provider/+/authorize\" {\n    capabilities = [\"read\", \"update\"]\n}\n"
    }
  ]
}
```

Let me check the existing Vault Policies in this brand new Vault instance first

```bash
$ vault policy list
default
root
```

Now, let me restore the policies from the Vault Policies Backup JSON file, and then check if overwriting of the policy happens as expected. So, let's run the `vault-policy-restore` command, but this time with the `-quiet` / `--quiet` flag to hide too much details and just show quiet progress. `vault-policy-restore` will overwrite the Vault Policies in the Vault instance if the policies already exist in the Vault instance

Let's see what actually happens, and notice that I check the `allow_secrets` policy once before doing the restore -

```bash
$ vault policy read allow_secrets
# KV v2 secrets engine mount path is "secret"
path "secret/*" {
  capabilities = ["read", "list"]
}

$ ./vault-policy-restore --quiet --file vault_policy_backup.json

restoring 4 vault policies in vault

restoring the following vault policies in vault: [allow_secrets allow_stage_kv_secrets allow_test_kv_secrets default]
....

$ vault policy list
allow_secrets
allow_stage_kv_secrets
allow_test_kv_secrets
default
root

$ vault policy read allow_secrets
path "secret/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}
```

As you can see, `allow_secrets` policy's content was overwritten. It was overwritten with the `allow_secrets` policy content that was passed through the backup policies JSON file.

So, this is how `vault-policy-restore` works :D

## Demo 3

In this demo, let show how `default` Vault policy can also be overwritten. We had been overwriting the `default` policy in the previous demo, but with the same policy content as what it already had, so it wasn't obvious that it got overwritten. This time, let's properly overwrite it with some different policy content to see the overwriting actually happen in a clear way

Let me run a brand new Vault server first

```bash
vault server -dev -dev-root-token-id root -dev-listen-address 127.0.0.1:8200 
```

Now, let me check the existing Vault Policies in this brand new Vault instance first

```bash
$ vault policy list
default
root

$ vault policy read default
# Allow tokens to look up their own properties
path "auth/token/lookup-self" {
    capabilities = ["read"]
}

# Allow tokens to renew themselves
path "auth/token/renew-self" {
    capabilities = ["update"]
}

# Allow tokens to revoke themselves
path "auth/token/revoke-self" {
    capabilities = ["update"]
}

# Allow a token to look up its own capabilities on a path
path "sys/capabilities-self" {
    capabilities = ["update"]
}

# Allow a token to look up its own entity by id or name
path "identity/entity/id/{{identity.entity.id}}" {
  capabilities = ["read"]
}
path "identity/entity/name/{{identity.entity.name}}" {
  capabilities = ["read"]
}


# Allow a token to look up its resultant ACL from all policies. This is useful
# for UIs. It is an internal path because the format may change at any time
# based on how the internal ACL features and capabilities change.
path "sys/internal/ui/resultant-acl" {
    capabilities = ["read"]
}

# Allow a token to renew a lease via lease_id in the request body; old path for
# old clients, new path for newer
path "sys/renew" {
    capabilities = ["update"]
}
path "sys/leases/renew" {
    capabilities = ["update"]
}

# Allow looking up lease properties. This requires knowing the lease ID ahead
# of time and does not divulge any sensitive information.
path "sys/leases/lookup" {
    capabilities = ["update"]
}

# Allow a token to manage its own cubbyhole
path "cubbyhole/*" {
    capabilities = ["create", "read", "update", "delete", "list"]
}

# Allow a token to wrap arbitrary values in a response-wrapping token
path "sys/wrapping/wrap" {
    capabilities = ["update"]
}

# Allow a token to look up the creation time and TTL of a given
# response-wrapping token
path "sys/wrapping/lookup" {
    capabilities = ["update"]
}

# Allow a token to unwrap a response-wrapping token. This is a convenience to
# avoid client token swapping since this is also part of the response wrapping
# policy.
path "sys/wrapping/unwrap" {
    capabilities = ["update"]
}

# Allow general purpose tools
path "sys/tools/hash" {
    capabilities = ["update"]
}
path "sys/tools/hash/*" {
    capabilities = ["update"]
}

# Allow checking the status of a Control Group request if the user has the
# accessor
path "sys/control-group/request" {
    capabilities = ["update"]
}

# Allow a token to make requests to the Authorization Endpoint for OIDC providers.
path "identity/oidc/provider/+/authorize" {
    capabilities = ["read", "update"]
}

$ vault policy read root
No policy named: root
```

Now, let's create a sample `vault_policy_backup.json` file like this -

```json
{
  "policies": [
    {
      "name": "default",
      "policy": "path \"secret/*\" {\n  capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"]\n}\n"
    }
  ]
}
```

Now, let's run `vault-policy-restore` command with this `vault_policy_backup.json` file to overwrite the `default` policy in the Vault instance

```bash
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ export VAULT_TOKEN="root"

$ ./vault-policy-restore --quiet --file vault_policy_backup.json

restoring 1 vault policies in vault

restoring the following vault policies in vault: [default]
.

$ vault policy list
default
root

$ vault policy read default
path "secret/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}
```

As you can see, the big `default` Vault policy content got overwritten / replaced by the `vault_policy_backup.json` file content :D

# Caveats

- If for some reason / somehow you have a `root` policy in the JSON file, it will be ignored.

- As mentioned before, `vault-policy-restore` command will overwrite the Vault Policies in the Vault instance if the policies already exist in the Vault instance. So, it is recommended to make sure that you have a backup of the Vault Policies in the Vault instance before doing the restore. This includes overwriting the `default` policy too.

# Error Scenarios

There are quite some different kinds of errors that can occur while running `vault-policy-restore` command. Let's see some of them -

Access error like the following - if the Vault Token / Credentials used for the Vault is not valid / wrong / does not have enough access, then the tool throws errors similar to this -

```bash
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ export VAULT_TOKEN="some-wrong-token"

$ ./vault-policy-restore --quiet --file vault_policy_backup.json

restoring 4 vault policies in vault

restoring the following vault policies in vault: [allow_secrets allow_stage_kv_secrets allow_test_kv_secrets default]
.error writing 'allow_secrets' vault policy to vault: Error making API request.

URL: PUT http://127.0.0.1:8200/v1/sys/policies/acl/allow_secrets
Code: 403. Errors:

* permission denied
```

File reading errors like -

```bash
$ ./vault-policy-restore --quiet --file vault_policy_backup_does_NOT_exist.json
error reading vault policy backup json file: open vault_policy_backup_does_NOT_exist.json: no such file or directory
```

# Future Ideas

Talking about future ideas, here are some of the ideas for the future -
- Give warning/information to user about how restoring of root policy is not done, so that they don't have to read the docs so much to understand this information. Reasoning - This is because - empty `root` Vault Policy - nothing to restore to Vault. Also, Vault does NOT support updating it / changing it - which would happen during the restore processThe `root` policy can neither be deleted, created or updated (or changed).
