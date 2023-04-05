package jet

import (
	"TrackMaster/pkg"
	"TrackMaster/third_party"
	"encoding/json"
)

type EnumValue struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"desc"`
	Status      string `json:"status"`
	Deprecated  bool   `json:"deprecated"`
}

type Type struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Status     string      `json:"status"`
	EnumValues []EnumValue `json:"enumValues"`
}

type enumResponse struct {
	Data []Type `json:"data"`
}

var enumTypeFetcher = &third_party.ThirdPartyDataFetcher{
	Host:    "https://jet-plus.midway.run",
	Path:    "/v1/internals/enumType",
	Query:   nil, //需要带上project参数，值是project id
	OnError: nil,
}

func GetEnumTypes(projectID string) ([]Type, *pkg.Error) {
	query := make(map[string][]string)
	query["project"] = []string{projectID}
	enumTypeFetcher.Query = query

	body, err := enumTypeFetcher.FetchData("")
	if err != nil {
		return nil, pkg.NewError(pkg.ServerError, err.Error())
	}

	response := enumResponse{}
	err1 := json.Unmarshal(body, &response)
	if err1 != nil {
		return nil, pkg.NewError(pkg.ServerError, err1.Error())
	}

	return response.Data, nil
}
