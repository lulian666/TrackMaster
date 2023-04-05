package jet

import (
	"TrackMaster/pkg"
	"TrackMaster/third_party"
	"encoding/json"
	"strconv"
	"strings"
)

type logResponse struct {
	Data []Log `json:"data"`
}

type Log struct {
	Event  string    `json:"event"` //event name
	Filter string    `json:"filter"`
	ID     string    `json:"id"` //可以用这个id来区别log是否已经记录过
	Log    LogDetail //log具体内容在这个字段里
	LogStr string    `json:"log"`
	UserID string    `json:"userId"`
}

// LogDetail
// 几乎涵盖了一个log里可能有的字段
// play_info和search_info暂时没写，不想结构太大了
// 如果发现取不到log里的值，那可能是需要维护这个struct了
type LogDetail struct {
	AppAddInfo     AppAddInfo     `json:"app_add_info"`
	AbtestInfo     AbtestInfo     `json:"abtest_info"`
	ContentAddInfo ContentAddInfo `json:"content_add_info"`
	ContentInfo    ContentInfo    `json:"content_info"`
	EventInfo      EventInfo      `json:"event_info"`
	DeviceInfo     DeviceInfo     `json:"device_info"`
	PageInfo       PageInfo       `json:"page_info"`
	WebInfo        WebInfo        `json:"web_info"`

	BuildCode          string  `json:"build_code"`
	Platform           string  `json:"platform"`
	AppID              string  `json:"$app_id"`
	AppName            string  `json:"$app_name"`
	AppVersion         string  `json:"$app_version"`
	Carrier            string  `json:"$carrier"`
	DeviceID           string  `json:"$device_id"`
	IDFV               string  `json:"$idfv"`
	IP                 string  `json:"$ip"`
	UA                 string  `json:"$ua"`
	IsFirstDay         bool    `json:"$is_first_day"`
	Lib                string  `json:"$lib"`
	LibMethod          string  `json:"$lib_method"`
	LibVersion         string  `json:"$lib_version"`
	Manufacturer       string  `json:"$manufacturer"`
	Model              string  `json:"$model"`
	NetworkType        string  `json:"$network_type"`
	OS                 string  `json:"$os"`
	OSVersion          string  `json:"$os_version"`
	ScreenHeight       int     `json:"$screen_height"`
	ScreenWidth        int     `json:"$screen_width"`
	TimezoneOffset     int     `json:"$timezone_offset"`
	Wifi               bool    `json:"$wifi"`
	EventDuration      float32 `json:"event_duration"`
	ExtraWebAbtestInfo string  `json:"extra_web_abtest_info"`
}

type ContentAddInfo struct {
	Amount  int    `json:"amount"`
	Content string `json:"content"`
	Id      string `json:"id"`
	Input   string `json:"input"`
	Source  string `json:"source"`
	Title   string `json:"title"`
	Type    string `json:"type"`
}

type ContentInfo struct {
	Id            string            `json:"id"`
	Type          string            `json:"type"`
	Content       string            `json:"content"`
	URL           string            `json:"url"`
	Title         string            `json:"title"`
	ReadTrackInfo map[string]string `json:"read_track_info"`
	Source        string            `json:"source"`
	Status        string            `json:"status"`
	Duration      int               `json:"duration"`
	Count         int               `json:"count"`
}

type EventInfo struct {
	Action          string `json:"action"`
	CurrentPageName string `json:"current_page_name"`
	Event           string `json:"event"`
	SourcePageName  string `json:"source_page_name"`
}

type AppAddInfo struct {
	IsPurchased bool   `json:"is_purchased"`
	IsScreenOff bool   `json:"is_screen_off"`
	Type        string `json:"type"`
	UtmSource   string `json:"utm_source"`
}

type AbtestInfo struct {
	DiscoveryFeedBeta              string `json:"discovery_feed_beta"`
	DiscoveryFeedProd              string `json:"discovery_feed_prod"`
	ExclusiveEpisodeRecommendation string `json:"exclusive_episode_recommendation"`
}

type DeviceInfo struct {
	CarBrand   string `json:"car_brand"`
	LocalID    string `json:"local_id"`
	LocalName  string `json:"local_name"`
	RemoteID   string `json:"remote_id"`
	RemoteName string `json:"remote_name"`
}

type PageInfo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
	URL   string `json:"url"`
}

type PageInfoWeb struct {
	ID    string `json:"page_info$$id"`
	Title string `json:"page_info$$title"`
	Type  string `json:"page_info$$type"`
	URL   string `json:"page_info$$url"`
}

type WebInfo struct {
	AbtestInfo      map[string]string `json:"abtest_info"`
	Action          string            `json:"action"`
	Campaign        string            `json:"campaign"`
	ExtraAbtestInfo string            `json:"extra_abtest_info"`
	Host            string            `json:"host"`
	ID              string            `json:"id"`
	Label           string            `json:"label"`
	PageName        string            `json:"page_name"`
	ShareDepth      int               `json:"share_depth"`
	ShareDistinctID string            `json:"share_distinct_id"`
	Source          string            `json:"source"`
}

type WebInfoWeb struct {
	AbtestInfo      map[string]string `json:"web_info$$abtest_info"`
	Action          string            `json:"web_info$$action"`
	Campaign        string            `json:"web_info$$campaign"`
	ExtraAbtestInfo string            `json:"web_info$$extra_abtest_info"`
	Host            string            `json:"web_info$$host"`
	ID              string            `json:"web_info$$id"`
	Label           string            `json:"web_info$$label"`
	PageName        string            `json:"web_info$$page_name"`
	ShareDepth      int               `json:"web_info$$share_depth"`
	ShareDistinctID string            `json:"web_info$$share_distinct_id"`
	Source          string            `json:"web_info$$source"`
}

var logFetcher = &third_party.ThirdPartyDataFetcher{
	Host:    "https://jet-plus.midway.run",
	Path:    "/v1/internals/eventLog",
	Query:   nil, //需要传filter参数，值为filter的id
	OnError: nil,
}

func GetLogs(filterID string) ([]Log, *pkg.Error) {
	query := make(map[string][]string)
	query["filter"] = []string{filterID}
	logFetcher.Query = query

	body, err := logFetcher.FetchData("")
	if err != nil {
		return nil, pkg.NewError(pkg.ServerError, err.Error())
	}

	response := logResponse{}
	err1 := json.Unmarshal(body, &response)
	if err1 != nil {
		return nil, pkg.NewError(pkg.ServerError, err1.Error())
	}

	// 因为body里给的log是字符串类型，所以要多这一步转化
	for i := range response.Data {
		logDetail := LogDetail{}
		err1 = json.Unmarshal([]byte(response.Data[i].LogStr), &logDetail)
		if err1 != nil {
			return nil, pkg.NewError(pkg.ServerError, err1.Error())
		}
		response.Data[i].Log = logDetail
	}

	return response.Data, nil
}

// ClearLogs
// 读取过log以后即刻清除，保证下次读取时全是未读的log
func ClearLogs(filterID string) *pkg.Error {
	query := make(map[string][]string)
	query["filter"] = []string{filterID}
	logFetcher.Query = query

	_, err := logFetcher.PatchData("DELETE", "", nil)
	if err != nil {
		return pkg.NewError(pkg.ServerError, err.Error())
	}

	return nil
}

// Get
// 默认key的结构一定是xx.xx
func (log *LogDetail) Get(key string, raw string) (string, bool) {
	keys := strings.Split(key, ".")
	// 不想大量使用反射的画可以用switch
	// Go语言的 switch 语句在性能方面非常高效
	if len(keys) != 2 {
		return "", false
	}

	value := ""

	switch keys[0] {
	case "app_add_info":
		st := log.AppAddInfo
		switch keys[1] {
		case "is_purchased":
			value = strconv.FormatBool(st.IsPurchased)
		case "is_screen_off":
			value = strconv.FormatBool(st.IsScreenOff)
		case "type":
			value = st.Type
		case "utm_source":
			value = st.UtmSource
		}
	case "abtest_info":
		st := log.AbtestInfo
		switch keys[1] {
		case "discovery_feed_beta":
			value = st.DiscoveryFeedBeta
		case "discovery_feed_prod":
			value = st.DiscoveryFeedProd
		case "exclusive_episode_recommendation":
			value = st.ExclusiveEpisodeRecommendation
		}
	case "content_add_info":
		st := log.ContentAddInfo
		switch keys[1] {
		case "amount":
			value = strconv.Itoa(st.Amount)
		case "content":
			value = st.Content
		case "id":
			value = st.Id
		case "input":
			value = st.Input
		case "source":
			value = st.Source
		case "title":
			value = st.Title
		case "type":
			value = st.Type
		}
	case "content_info":
		st := log.ContentInfo
		switch keys[1] {
		case "id":
			value = st.Id
		case "type":
			value = st.Type
		case "content":
			value = st.Content
		case "url":
			value = st.URL
		case "title":
			value = st.Title
		case "read_track_info":
			b, _ := json.Marshal(st.ReadTrackInfo)
			value = string(b)
		case "source":
			value = st.Source
		case "status":
			value = st.Status
		case "duration":
			value = strconv.Itoa(st.Duration)
		case "count":
			value = strconv.Itoa(st.Count)
		}
	case "event_info":
		st := log.EventInfo
		switch keys[1] {
		case "action":
			value = st.Action
		case "current_page_name":
			value = st.CurrentPageName
		case "event":
			value = st.Event
		case "source_page_name":
			value = st.SourcePageName
		}
	case "device_info":
		st := log.DeviceInfo
		switch keys[1] {
		case "car_brand":
			value = st.CarBrand
		case "local_id":
			value = st.LocalID
		case "local_name":
			value = st.LocalName
		case "remote_id":
			value = st.RemoteID
		case "remote_name":
			value = st.RemoteName
		}
	case "page_info":
		st := log.PageInfo
		switch keys[1] {
		case "id":
			value = st.ID
		case "title":
			value = st.Title
		case "type":
			value = st.Type
		case "url":
			value = st.URL
		}
		if value == "" {
			webLog := make(map[string]interface{})
			err := json.Unmarshal([]byte(raw), &webLog)
			if err != nil {
				return "", false
			}
			k := keys[0] + "$$" + keys[1]
			v, ok := webLog[k].(string)
			if ok {
				value = v
			}
		}
	case "web_info":
		st := log.WebInfo
		switch keys[1] {
		case "abtest_info":
			b, _ := json.Marshal(st.AbtestInfo)
			value = string(b)
		case "action":
			value = st.Action
		case "campaign":
			value = st.Campaign
		case "extra_abtest_info":
			value = st.ExtraAbtestInfo
		case "host":
			value = st.Host
		case "id":
			value = st.ID
		case "label":
			value = st.Label
		case "page_name":
			value = st.PageName
		case "share_depth":
			value = strconv.Itoa(st.ShareDepth)
		case "share_distinct_id":
			value = st.ShareDistinctID
		case "source":
			value = st.Source
		}
		if value == "" {
			webLog := make(map[string]interface{})
			err := json.Unmarshal([]byte(raw), &webLog)
			if err != nil {
				return "", false
			}
			k := keys[0] + "$$" + keys[1]
			value = webLog[k].(string)
		}
	}

	return value, true
}
