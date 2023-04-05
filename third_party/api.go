package third_party

import (
	"TrackMaster/pkg"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type DataFetcher interface {
	FetchData(params map[string]string) ([]byte, error)
}

type ThirdPartyDataFetcher struct {
	Host    string
	Path    string
	Query   url.Values
	OnError func() ([]byte, error)
}

func (f *ThirdPartyDataFetcher) FetchData(params string) ([]byte, *pkg.Error) {
	// 构造URL
	u, err := url.Parse(f.Host + f.Path)
	if err != nil {
		return nil, pkg.NewError(pkg.ServerError, err.Error())
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
		return nil, pkg.NewError(pkg.ServerError, "something went wrong with third party api").WithDetails(err.Error())
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
			if r, err := f.OnError(); err == nil {
				return r, nil
			}
			return nil, pkg.NewError(pkg.ServerError, err.Error()).WithDetails(fmt.Sprintf("failed to fetch data, status code: %d \nrequest url is %s", resp.StatusCode, u))
		}
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, pkg.NewError(pkg.ServerError, err.Error())
	}

	return body, nil
}

func (f *ThirdPartyDataFetcher) PostData(params string, body []byte) ([]byte, *pkg.Error) {
	// 构造URL
	u, err := url.Parse(f.Host + f.Path)
	if err != nil {
		return nil, pkg.NewError(pkg.ServerError, err.Error())
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

	resp, err := http.Post(u.String(), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, pkg.NewError(pkg.ServerError, err.Error())
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	// 处理响应
	if resp.StatusCode != http.StatusOK {
		if f.OnError != nil {
			if r, err := f.OnError(); err == nil {
				return r, nil
			}
			return nil, pkg.NewError(pkg.ServerError, err.Error()).WithDetails(fmt.Sprintf("failed to fetch data, status code: %d \nrequest url is %s", resp.StatusCode, u))
		}
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, pkg.NewError(pkg.ServerError, err.Error())
	}

	return resBody, nil
}

func (f *ThirdPartyDataFetcher) PatchData(method string, params string, body []byte) ([]byte, *pkg.Error) {
	// 构造URL
	u, err := url.Parse(f.Host + f.Path)
	if err != nil {
		return nil, pkg.NewError(pkg.ServerError, err.Error())
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

	// 创建请求
	req, err := http.NewRequest(strings.ToUpper(method), u.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	// 处理响应
	if resp.StatusCode != http.StatusOK {
		if f.OnError != nil {
			if r, err := f.OnError(); err == nil {
				return r, nil
			}
			return nil, pkg.NewError(pkg.ServerError, err.Error()).WithDetails(fmt.Sprintf("failed to fetch data, status code: %d \nrequest url is %s", resp.StatusCode, u))
		}
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, pkg.NewError(pkg.ServerError, err.Error())
	}

	return resBody, nil
}
