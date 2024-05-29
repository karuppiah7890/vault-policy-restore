package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

type VaultPolicyBackup struct {
	Policies VaultPolicies `json:"policies"`
}

type VaultPolicies []VaultPolicy

type VaultPolicy struct {
	Name   string `json:"name"`
	Policy string `json:"policy"`
}

func convertJSONToVaultPolicyBackup(JSONData []byte) (*VaultPolicyBackup, error) {
	vaultPolicyBackup, err := fromJSON(JSONData)
	if err != nil {
		return nil, err
	}
	return vaultPolicyBackup, nil
}

func restoreVaultPolicies(client *api.Client, vaultPolicyBackup *VaultPolicyBackup, quietProgress bool) error {

	policies := vaultPolicyBackup.Policies

	policyNames := make([]string, 0, len(policies))
	for _, policy := range policies {
		// Skip root policy
		if policy.Name == "root" {
			continue
		}
		policyNames = append(policyNames, policy.Name)
	}

	fmt.Fprintf(os.Stdout, "\nrestoring %d vault policies in vault\n", len(policies))
	fmt.Fprintf(os.Stdout, "\nrestoring the following vault policies in vault: %+v\n", policyNames)

	// Restore all Vault policies
	for _, policy := range policies {
		// Skip root policy
		if policy.Name == "root" {
			continue
		}
		if quietProgress {
			fmt.Fprintf(os.Stdout, ".")
		} else {
			fmt.Fprintf(os.Stdout, "\nrestoring `%s` vault policy\n", policy.Name)
			fmt.Fprintf(os.Stdout, "\n`%s` vault policy rules: %+v\n", policy.Name, policy.Policy)
		}
		err := client.Sys().PutPolicy(policy.Name, policy.Policy)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error writing '%s' vault policy to vault: %s\n", policy.Name, err)
			os.Exit(1)
		}
	}

	fmt.Fprintf(os.Stdout, "\n")

	return nil
}
