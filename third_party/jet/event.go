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

func GetEvents(id string) ([]Event, error) {
	eventFetcher.Path = eventFetcher.Path + "/" + id
	body, err := eventFetcher.FetchData(nil)
	if err != nil {
		return nil, err
	}
	response := eventResponse{}
	err = json.Unmarshal(body, &response)

	return response.Events, nil
}
