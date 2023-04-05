package podcast

import (
	"TrackMaster/pkg"
	"TrackMaster/third_party"
	"encoding/json"
	"strings"
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
		Title: "podcast内部接口调用失败",
		Desc:  "可能是type或id传错，可能是该type不支持查询，也可能是接口不稳定",
	}
	return json.Marshal(content)
}

var podcastFetcher = &third_party.ThirdPartyDataFetcher{
	Host:    "http://podcast-service.podcast-prod.svc.cluster.local:3000",
	Path:    "/internal/podcast/get",
	Query:   nil,
	OnError: onError,
}

func GetPodcast(pid string) (Content, *pkg.Error) {
	query := make(map[string][]string)
	query["pid"] = []string{pid}

	podcastFetcher.Query = query
	body, err := podcastFetcher.FetchData("")
	if err != nil {
		return Content{}, err
	}

	response := contentResponse{}
	err1 := json.Unmarshal(body, &response)
	if err1 != nil {
		return Content{}, pkg.NewError(pkg.ServerError, err1.Error())
	}

	// 有的内容太长了
	if len(response.Data.Desc) > 20 {
		response.Data.Desc = response.Data.Desc[:20]
	}

	return response.Data, nil
}

var episodeFetcher = &third_party.ThirdPartyDataFetcher{
	Host:    "http://podcast-service.podcast-prod.svc.cluster.local:3000",
	Path:    "/internal/episode/get",
	Query:   nil,
	OnError: onError,
}

func GetEpisode(eid string) (Content, *pkg.Error) {
	query := make(map[string][]string)
	query["eid"] = []string{eid}

	episodeFetcher.Query = query
	body, err := episodeFetcher.FetchData("")
	if err != nil {
		return Content{}, err
	}

	response := contentResponse{}
	err1 := json.Unmarshal(body, &response)
	if err1 != nil {
		return Content{}, pkg.NewError(pkg.ServerError, err1.Error())
	}

	if len(response.Data.Desc) > 20 {
		response.Data.Desc = response.Data.Desc[:20]
	}

	return response.Data, nil
}

var collectionFetcher = &third_party.ThirdPartyDataFetcher{
	Host:    "http://podcast-service.podcast-prod.svc.cluster.local:3000",
	Path:    "/internal/collection/get",
	Query:   nil,
	OnError: onError,
}

func GetCollection(id string) (Content, *pkg.Error) {
	query := make(map[string][]string)
	query["id"] = []string{id}

	collectionFetcher.Query = query
	body, err := collectionFetcher.FetchData("")
	if err != nil {
		return Content{}, err
	}

	response := contentResponse{}
	err1 := json.Unmarshal(body, &response)
	if err1 != nil {
		return Content{}, pkg.NewError(pkg.ServerError, err1.Error())
	}

	if len(response.Data.Desc) > 20 {
		response.Data.Desc = response.Data.Desc[:20]
	}

	return response.Data, nil
}

func GetContentByTypeAndID(t string, id string) (Content, *pkg.Error) {
	switch strings.ToLower(t) {
	case "collection":
		return GetCollection(id)
	case "podcast":
		return GetPodcast(id)
	case "episode":
		return GetEpisode(id)
	}
	return Content{}, nil
}
