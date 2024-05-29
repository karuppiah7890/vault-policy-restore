package main

import "encoding/json"

func fromJSON(JSONData []byte) (*VaultPolicyBackup, error) {
	vaultPolicyBackup := &VaultPolicyBackup{}
	err := json.Unmarshal(JSONData, vaultPolicyBackup)
	if err != nil {
		return nil, err
	}
	return vaultPolicyBackup, nil
}
