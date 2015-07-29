package request

type LoginResponse struct {
	UserToken    string `json:"userToken"`
	CharacterIDs []int  `json:"characters"`
}

type MachineRegisterResponse struct {
	MachineId    int    `json:"machineId"`
	MachineToken string `json:"machineToken"`
}
