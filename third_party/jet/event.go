package jet

import (
	"TrackMaster/pkg"
	"TrackMaster/third_party"
	"encoding/json"
)

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

var eventFetcher = &third_party.ThirdPartyDataFetcher{
	Host:    "https://jet-plus.midway.run",
	Path:    "/v1/internals/requirement",
	Query:   nil,
	OnError: nil,
}

//需要带上story id作为param

func GetEvents(storyID string) ([]Event, *pkg.Error) {
	body, err := eventFetcher.FetchData(storyID)
	if err != nil {
		return nil, pkg.NewError(pkg.ServerError, err.Error())
	}

	response := eventResponse{}
	err1 := json.Unmarshal(body, &response)
	if err1 != nil {
		return nil, pkg.NewError(pkg.ServerError, err1.Error())
	}

	return response.Events, nil
}
