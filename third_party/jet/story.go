package jet

import (
	"TrackMaster/pkg"
	"TrackMaster/third_party"
	"encoding/json"
)

type Story struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type storyResponse struct {
	Data []Story `json:"data"`
}

var storyFetcher = &third_party.ThirdPartyDataFetcher{
	Host:    "https://jet-plus.midway.run",
	Path:    "/v1/internals/requirement",
	Query:   nil, //需要带上project参数，值是project id
	OnError: nil,
}

// GetStories
// 只取最新修改过的15条需求
func GetStories(projectID string) ([]Story, *pkg.Error) {
	query := make(map[string][]string)
	query["project"] = []string{projectID}
	//{"sortKey": "updatedAt", "sortDirection": "DESC", "pageSize": 15, "page": 1}
	query["sortKey"] = []string{"updatedAt"}
	query["sortDirection"] = []string{"DESC"}
	query["pageSize"] = []string{"15"}
	query["page"] = []string{"1"}

	storyFetcher.Query = query
	body, err := storyFetcher.FetchData("")
	if err != nil {
		return nil, pkg.NewError(pkg.ServerError, err.Error())
	}
	response := storyResponse{}
	err1 := json.Unmarshal(body, &response)
	if err1 != nil {
		return nil, pkg.NewError(pkg.ServerError, err1.Error())
	}

	return response.Data, nil
}
