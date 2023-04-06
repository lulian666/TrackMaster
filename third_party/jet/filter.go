package jet

import (
	"TrackMaster/pkg"
	"TrackMaster/third_party"
	"encoding/json"
)

const (
	READY     status = "READY"
	RECORDING status = "RECORDING"
	PAUSED    status = "PAUSED"
	STOPPED   status = "STOPPED"
)

type status string

type Filter struct {
	Events  []string `json:"events"`  // 需要event name
	Project string   `json:"project"` // id
	UserIDs []string `json:"userIds"`
	ID      string   `json:"id"`
	Status  status   `json:"status"`
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
func (f Filter) Update() *pkg.Error {
	reqBody, err1 := json.Marshal(f)
	if err1 != nil {
		return pkg.NewError(pkg.ServerError, err1.Error())
	}

	_, err := filterFetcher.PatchData("PATCH", f.ID, reqBody)
	if err != nil {
		return err
	}

	return nil
}

func (f Filter) Create() (FilterRes, *pkg.Error) {
	body, err1 := json.Marshal(f)
	if err1 != nil {
		return FilterRes{}, pkg.NewError(pkg.ServerError, err1.Error())
	}

	resBody, err := filterFetcher.PatchData("POST", "", body)
	if err != nil {
		return FilterRes{}, err
	}

	filter := FilterRes{}
	err1 = json.Unmarshal(resBody, &filter)

	if err1 != nil {
		return FilterRes{}, pkg.NewError(pkg.ServerError, err1.Error())
	}

	return filter, nil
}
