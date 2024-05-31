package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

var usage = `usage: vault-policy-restore [-quiet|--quiet] [-file|-file <vault-policy-backup-json-file-path>]

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
`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "%s", usage)
	}
	quietProgress := flag.Bool("quiet", false, "quiet progress")
	vaultPolicyBackupJsonFileName := flag.String("file", "vault_policy_backup.json", "vault policy backup json file path")
	flag.Parse()

	// TODO: Take a backup of the destination policy/policies in case it is already present,
	// regardless of if they have the same name as source policy.
	// So that we have a backup just in case, especially before overwriting.

	// Question: Should we do a complete backup before doing policy copy one by one?
	// Or should we do a backup of each policy one by one? When we copy them one by one
	// that is. Basically - take a backup one by one, at vaultPolicyCopy() function level or
	// take a complete backup of all destination vault policies at vaultPolicyCopyAll()
	// function level.

	config := api.DefaultConfig()
	client, err := api.NewClient(config)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating vault client: %s\n", err)
		os.Exit(1)
	}

	vaultPolicyBackupJsonFileContent, err := readFile(*vaultPolicyBackupJsonFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading vault policy backup json file: %s\n", err)
		os.Exit(1)
	}

	vaultPolicyBackup, err := convertJSONToVaultPolicyBackup(vaultPolicyBackupJsonFileContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing vault policy backup json file content: %s\n", err)
		os.Exit(1)
	}

	err = restoreVaultPolicies(client, vaultPolicyBackup, *quietProgress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error restoring vault policies: %s\n", err)
		os.Exit(1)
	}
}
