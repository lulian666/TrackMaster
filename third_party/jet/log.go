package jet

import (
	"TrackMaster/third_party"
	"encoding/json"
)

type logResponse struct {
	Data []Log `json:"data"`
}

type Log struct {
	Event  string `json:"event"` //event name
	Filter string `json:"filter"`
	ID     string `json:"id"`  //可以用这个id来区别log是否已经记录过
	Log    string `json:"log"` //log具体内容在这个字段里
	UserID string `json:"userId"`
}

var logFetcher = &third_party.ThirdPartyDataFetcher{
	Host:    "https://jet-plus.midway.run",
	Path:    "/v1/internals/eventLog",
	Query:   nil, //需要传filter参数，值为filter的id
	OnError: nil,
}

func GetLogs(filterID string) ([]Log, error) {
	query := make(map[string][]string)
	query["filter"] = []string{filterID}
	logFetcher.Query = query

	body, err := logFetcher.FetchData("")
	if err != nil {
		return nil, err
	}

	response := logResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}
