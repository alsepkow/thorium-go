package request

type LoginResponse struct {
	UserToken    string `json:"userToken"`
	CharacterIDs []int  `json:"characters"`
}
