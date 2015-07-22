package thordb

import "time"

type AccountSession struct {
	UserID int    `json:"uid"`
	Token  string `json:"token"`
	Account
}

type Account struct {
	UserID         int       `json:"uid"`
	CharacterIDs   []int     `json:"characters"`
	HashedPassword []byte    `json:"hashedPassword"`
	Salt           []byte    `json:"salt"`
	Algorithm      string    `json:"hashAlgorithm"`
	CreatedOn      time.Time `json:"createdOn"`
	LastLogin      time.Time `json:"lastLogin"`
}

type AccountPublicView struct {
	UserID       int       `json:"uid"`
	CharacterIDs []int     `json:"characters"`
	CreatedOn    time.Time `json:"createdOn"`
	LastLogin    time.Time `json:"lastLogin"`
}

func (account *Account) NewPublicView() AccountPublicView {

	var publicView AccountPublicView
	publicView.UserID = account.UserID
	publicView.CharacterIDs = account.CharacterIDs
	publicView.CreatedOn = account.CreatedOn
	publicView.LastLogin = account.LastLogin
}

func (account *Account) Validate() (bool, error) {
	// helper function to keep basic validation here?
	return true, nil
}

type Character struct {
	ID     int
	UserID int
	Name   string
	GameData
}

type GameData struct {
	Weapons   []int
	Inventory []int
	Health    float64
}
