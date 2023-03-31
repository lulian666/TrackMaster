package jet

import "encoding/json"

type Event struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Desc             string  `json:"desc"`
	EventDefinitions []Field `json:"eventDefinitions"`
}

type Field struct {
	ID     string    `json:"id"`
	Name   string    `json:"name"`
	Type   FieldType `json:"type"`
	Values []string  `json:"values"`
	Note   string    `json:"note"`
}

type FieldType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type eventResponse struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Desc   string  `json:"desc"`
	Events []Event `json:"events"`
}

var eventFetcher = &ThirdPartyDataFetcher{
	Host:    "https://jet-plus.midway.run",
	Path:    "/v1/internals/requirement",
	Query:   nil,
	OnError: nil,
}

//需要带上story id作为param

func GetEvents(storyID string) ([]Event, error) {
	body, err := eventFetcher.FetchData(storyID)
	if err != nil {
		return nil, err
	}
	response := eventResponse{}
	err = json.Unmarshal(body, &response)

	return response.Events, nil
}
