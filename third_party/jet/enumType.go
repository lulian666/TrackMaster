package jet

import "encoding/json"

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

var enumTypeFetcher = &ThirdPartyDataFetcher{
	Host:    "https://jet-plus.midway.run",
	Path:    "/v1/internals/enumType",
	Query:   nil, //需要带上project参数，值是project id
	OnError: nil,
}

func GetEnumTypes(projectID string) ([]Type, error) {
	query := make(map[string][]string)
	query["project"] = []string{projectID}
	enumTypeFetcher.Query = query
	body, err := enumTypeFetcher.FetchData(nil)
	if err != nil {
		return nil, err
	}
	response := enumResponse{}
	err = json.Unmarshal(body, &response)

	return response.Data, nil
}
