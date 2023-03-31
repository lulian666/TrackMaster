package jet

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

type DataFetcher interface {
	FetchData(params map[string]string) ([]byte, error)
}

type ThirdPartyDataFetcher struct {
	Host    string
	Path    string
	Query   url.Values
	OnError func(statusCode int)
}

func (f *ThirdPartyDataFetcher) FetchData(params string) ([]byte, error) {
	// 构造URL
	u, err := url.Parse(f.Host + f.Path)
	if err != nil {
		return nil, err
	}

	// 添加param
	if params != "" {
		u.Path = path.Join(u.Path, "/", params)
	}

	// 添加查询参数
	q := u.Query()
	for k, vs := range f.Query {
		for _, v := range vs {
			q.Add(k, v)
		}
	}
	u.RawQuery = q.Encode()

	// 发送请求
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	// 处理响应
	if resp.StatusCode != http.StatusOK {
		// 调用回调函数
		if f.OnError != nil {
			f.OnError(resp.StatusCode)
		}
		return nil, fmt.Errorf("failed to fetch data, status code: %d \nrequest url is %s", resp.StatusCode, u)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
