package podcast

import (
	"TrackMaster/third_party"
	"encoding/json"
)

type Content struct {
	Title string `json:"title"`
	Desc  string `json:"description"`
}

type contentResponse struct {
	Data Content `json:"data"`
}

func onError() ([]byte, error) {
	content := Content{
		Title: "未找到相应id的内容",
		Desc:  "可能是type或id传错，也可能是该type不支持查询",
	}
	return json.Marshal(content)
}

var podcastFetcher = &third_party.ThirdPartyDataFetcher{
	Host:    "http://podcast-service.podcast-prod.svc.cluster.local:3000",
	Path:    "/internal/podcast/get",
	Query:   nil,
	OnError: onError,
}

func GetPodcast(pid string) (Content, error) {
	query := make(map[string][]string)
	query["pid"] = []string{pid}

	podcastFetcher.Query = query
	body, err := podcastFetcher.FetchData("")
	if err != nil {
		return Content{}, err
	}

	response := contentResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return Content{}, err
	}

	return response.Data, nil
}

var episodeFetcher = &third_party.ThirdPartyDataFetcher{
	Host:    "http://podcast-service.podcast-prod.svc.cluster.local:3000",
	Path:    "/internal/episode/get",
	Query:   nil,
	OnError: onError,
}

func GetEpisode(eid string) (Content, error) {
	query := make(map[string][]string)
	query["eid"] = []string{eid}

	episodeFetcher.Query = query
	body, err := episodeFetcher.FetchData("")
	if err != nil {
		return Content{}, err
	}

	response := contentResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return Content{}, err
	}

	return response.Data, nil
}

var collectionFetcher = &third_party.ThirdPartyDataFetcher{
	Host:    "http://podcast-service.podcast-prod.svc.cluster.local:3000",
	Path:    "/internal/collection/get",
	Query:   nil,
	OnError: onError,
}

func GetCollection(id string) (Content, error) {
	query := make(map[string][]string)
	query["id"] = []string{id}

	collectionFetcher.Query = query
	body, err := collectionFetcher.FetchData("")
	if err != nil {
		return Content{}, err
	}

	response := contentResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return Content{}, err
	}

	return response.Data, nil
}
