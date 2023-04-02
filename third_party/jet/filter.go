package jet

import (
	"TrackMaster/third_party"
	"encoding/json"
)

const (
	READY     = "READY"
	RECORDING = "RECORDING"
	PAUSED    = "PAUSED"
	STOPPED   = "STOPPED"
)

type status string

type Filter struct {
	Events  []string `json:"events"`  // 需要event name
	Project string   `json:"project"` // id
	UserIDs []string `json:"userIds"`
	ID      string   `json:"id"`
	Status  string   `json:"status"`
}

type FilterRes struct {
	Events  []string `json:"events"`
	Project Project  `json:"project"`
	UserID  string   `json:"userId"`
	UserIDs []string `json:"userIds"`
	ID      string   `json:"id"`
	Status  string   `json:"status"`
}

var filterFetcher = &third_party.ThirdPartyDataFetcher{
	Host:    "https://jet-plus.midway.run",
	Path:    "/v1/internals/eventLog/filter",
	Query:   nil,
	OnError: nil,
}

// Update
// filter的patch方法不返回任何数据
func (f Filter) Update() error {
	reqBody, err := json.Marshal(f)
	if err != nil {
		return err
	}

	_, err = filterFetcher.PatchData(f.ID, reqBody)
	if err != nil {
		return err
	}

	return nil
}

func (f Filter) Create() (FilterRes, error) {
	body, err := json.Marshal(f)
	if err != nil {
		return FilterRes{}, err
	}

	resBody, err := filterFetcher.PostData("", body)
	if err != nil {
		return FilterRes{}, err
	}

	filter := FilterRes{}
	err = json.Unmarshal(resBody, &filter)

	return filter, err
}
