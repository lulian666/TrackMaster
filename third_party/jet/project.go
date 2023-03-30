package jet

import "encoding/json"

type Response struct {
	Data []Project `json:"data"`
}

// Project jet中拿到的project结构
type Project struct {
	Name       string `json:"name"`
	CnName     string `json:"cnName"`
	Type       string `json:"type"`
	Registered bool   `json:"registered"`
	ID         string `json:"id"`
	CreatedAt  int64  `json:"createdAt"`
	UpdatedAt  int64  `json:"updatedAt"`
}

// projectFetcher 获取所有projects
var projectFetcher = &ThirdPartyDataFetcher{
	Host:    "https://jet-plus.midway.run",
	Path:    "/v1/internals/project",
	Query:   nil,
	OnError: nil,
}

func GetProjects() ([]Project, error) {
	body, err := projectFetcher.FetchData(nil)
	if err != nil {
		return nil, err
	}
	response := Response{}
	err = json.Unmarshal(body, &response)

	return response.Data, nil
}
