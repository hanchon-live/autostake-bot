package blockchain

import "fmt"

func GetValidator(address string) (string, error) {
	if address == "evmos1dgpv4leszpeg2jusx2xgyfnhdzghf3rfzw06t3" {
		return "evmosvaloper1dgpv4leszpeg2jusx2xgyfnhdzghf3rf0qq22v", nil
	} else if address == "evmos197ahcv2x9jj0nmvnen4sqqfffhygjga7wc7qkp" {
		return "evmosvaloper197ahcv2x9jj0nmvnen4sqqfffhygjga7rk3shu", nil
	}
	return "", fmt.Errorf("Validator address not registered")
}
