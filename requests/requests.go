package request

type NewGame struct {
	Map        string
	MaxPlayers int
}

type RegisterGame struct {
	MachineId int
	Port      int
}

type RegisterMachine struct {
	Port int
}
