package thordb

type AccountSession struct {
	UserID int    `json:"uid"`
	Token  string `json:"token"`
}
