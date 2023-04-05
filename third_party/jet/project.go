package jet

import (
	"TrackMaster/pkg"
	"TrackMaster/third_party"
	"encoding/json"
)

type projectResponse struct {
	Data []Project `json:"data"`
}

// Project jet中拿到的project结构
type Project struct {
	Name       string `json:"name"`
	CnName     string `json:"cnName"`
	Type       string `json:"type"`
	Registered bool   `json:"registered"`
	ID         string `json:"id"`
}

// projectFetcher 获取所有projects
var projectFetcher = &third_party.ThirdPartyDataFetcher{
	Host:    "https://jet-plus.midway.run",
	Path:    "/v1/internals/project",
	Query:   nil,
	OnError: nil,
}

func GetProjects() ([]Project, *pkg.Error) {
	body, err := projectFetcher.FetchData("")
	if err != nil {
		return nil, pkg.NewError(pkg.ServerError, "something went wrong with jet").WithDetails(err.Error())
	}
	response := projectResponse{}

	err1 := json.Unmarshal(body, &response)
	if err1 != nil {
		return nil, pkg.NewError(pkg.ServerError, err1.Error())
	}

	return response.Data, nil
}
