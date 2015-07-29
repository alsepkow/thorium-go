package request

type NewGame struct {
	CharacterToken string `json:"characterToken"`
	Map            string `json:"map"`
	MaxPlayers     int    `json:"maxPlayers"`
}

type RegisterGame struct {
	MachineId int `json:"machineId"`
	Port      int `json:"gameListenPort"`
}

type RegisterMachine struct {
	Port int `json:"serviceListenPort"`
}

type UnregisterMachine struct {
	MachineToken string `json:"machineToken"`
}

type MachineStatus struct {
	MachineToken   string  `json:"machineToken"`
	UsageCPU       float64 `json:"cpuUsagePct"`
	UsageNetwork   float64 `json:"networkUsagePct"`
	PlayerCapacity float64 `json:"playerCapacityPct"`
}

type Authentication struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateCharacter struct {
	AccountToken string `json:"accountToken"`
	Name         string `json:"name"`
}

type SelectCharacter struct {
	AccountToken string `json:"accountToken"`
	ID           int    `json:"characterId"`
}
