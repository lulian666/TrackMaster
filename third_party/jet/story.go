package jet

import "encoding/json"

type Story struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type storyResponse struct {
	Data []Story `json:"data"`
}

var storyFetcher = &ThirdPartyDataFetcher{
	Host:    "https://jet-plus.midway.run",
	Path:    "/v1/internals/requirement",
	Query:   nil, //需要带上project参数，值是project id
	OnError: nil,
}

// GetStories
// 只取最新修改过的15条需求
func GetStories(projectID string) ([]Story, error) {
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
		return nil, err
	}
	response := storyResponse{}
	err = json.Unmarshal(body, &response)

	return response.Data, nil
}
