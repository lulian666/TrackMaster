package request

type Request struct {
	AccountIDs []string `json:"accounts" binding:"required"`
	EventIDs   []string `json:"events" binding:"required"`
	Project    string   `json:"project" binding:"required"`
}
