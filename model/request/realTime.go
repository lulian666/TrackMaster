package request

import "TrackMaster/model"

type Start struct {
	AccountIDs []string `json:"accounts" binding:"required"`
	EventIDs   []string `json:"events" binding:"required"`
	Project    string   `json:"project" binding:"required"`
}

type Update struct {
	RecordID   string   `json:"record" binding:"required"`
	AccountIDs []string `json:"accounts" binding:"required"`
	EventIDs   []string `json:"events" binding:"required"`
}

type UpdateResult struct {
	RecordID string        `json:"record" binding:"required"`
	Fields   []model.Field `json:"fields"`
	EventID  string        `json:"event"`
}
