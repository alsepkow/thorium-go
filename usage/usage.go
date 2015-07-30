package usage

import "math/rand"

func GetCPU() (float64, error) {
	usageCpu := rand.NormFloat64()*10 + 50
	return usageCpu, nil
}

func GetNetworkUtilization() (float64, error) {
	usageNetwork := rand.NormFloat64()*10 + 50
	return usageNetwork, nil
}
