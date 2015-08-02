package request

type LoginResponse struct {
	UserToken    string `json:"userToken"`
	CharacterIDs []int  `json:"characters"`
}

type CharacterSessionResponse struct {
	CharacterToken string `json:"characterToken"`
	GameId         int    `json:"gameId"`
}

type MachineRegisterResponse struct {
	MachineId    int    `json:"machineId"`
	MachineToken string `json:"machineToken"`
}
