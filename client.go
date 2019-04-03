package httpclient

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

type AttachFile struct {
	FieldName string
	FileName  string
	Reader    io.Reader
}

// マルチパートフォームデータで送信します。
// 画像をリクエストする場合などに使用します。
func RequestHttpWithFile(requestUrl string, method string, parameters map[string]string, cookies []*http.Cookie,
	requestHeader map[string]string, files []AttachFile, proxy string) (*http.Response, error) {
	var b bytes.Buffer
	var fw io.Writer
	var err error
	w := multipart.NewWriter(&b)

	// Add file
	if files != nil {
		for _, v := range files {
			fw, err = w.CreateFormFile(v.FieldName, v.FileName)
			if err != nil {
				return nil, err
			}
			if _, err = io.Copy(fw, v.Reader); err != nil {
				return nil, err
			}
		}
	}
	// Add the other fields
	if parameters != nil {
		for key, val := range parameters {
			if fw, err = w.CreateFormField(key); err != nil {
				return nil, err
			}
			if _, err = fw.Write([]byte(val)); err != nil {
				return nil, err
			}
		}
	}

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	err = w.Close()
	if err != nil {
		return nil, err
	}

	// Request を生成
	req, err := http.NewRequest(method, requestUrl, &b)
	if err != nil {
		return nil, err
	}

	// Request Header
	setRequestHeaders(requestHeader, req)
	if req.Header.Get("User-Agent") == "" {
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; MAFSJS; rv:11.0) like Gecko")
	}
	req.Header.Add("Content-Type", w.FormDataContentType())
	//req.SetBasicAuth("112233", "445566")

	// Cookie
	setCookies(cookies, req)

	// 自動リダイレクトのオフ、プロキシ設定
	client := http.Client{
		// Go HTTP Client NOT Follow Redirects Automatically.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			return nil, err
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}
	return client.Do(req)
}

func setRequestHeaders(requestHeader map[string]string, req *http.Request) {
	if requestHeader != nil {
		for key, val := range requestHeader {
			req.Header.Add(key, val)
		}
	}
}

func setCookies(cookies []*http.Cookie, req *http.Request) {
	if cookies != nil {
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
	}
}

func Post(requestUrl string, parameters map[string]string, cookies []*http.Cookie,
	requestHeader map[string]string, rawData []byte, files []AttachFile, proxy string) (*http.Response, error) {
	if files == nil {
		return RequestHttp(requestUrl, "post", parameters, cookies, requestHeader, rawData, proxy)
	} else {
		return RequestHttpWithFile(requestUrl, "post", parameters, cookies, requestHeader, files, proxy)
	}
}

func Get(requestUrl string, parameters map[string]string, cookies []*http.Cookie,
	requestHeader map[string]string, rawData []byte, proxy string) (*http.Response, error) {
	return RequestHttp(requestUrl, "get", parameters, cookies, requestHeader, rawData, proxy)
}

func Put(requestUrl string, parameters map[string]string, cookies []*http.Cookie,
	requestHeader map[string]string, rawData []byte, proxy string) (*http.Response, error) {
	return RequestHttp(requestUrl, "put", parameters, cookies, requestHeader, rawData, proxy)
}

func Delete(requestUrl string, parameters map[string]string, cookies []*http.Cookie,
	requestHeader map[string]string, rawData []byte, proxy string) (*http.Response, error) {
	return RequestHttp(requestUrl, "delete", parameters, cookies, requestHeader, rawData, proxy)
}

// HTTP Request を行う関数です。
func RequestHttp(requestUrl string, method string, parameters map[string]string, cookies []*http.Cookie,
	requestHeader map[string]string, rawData []byte, proxy string) (*http.Response, error) {
	var req *http.Request
	var err error
	var data io.Reader
	values := url.Values{}
	if rawData != nil {
		data = bytes.NewReader(rawData)
	} else {
		if parameters != nil {
			for key, val := range parameters {
				values.Add(key, val)
			}
		}
		data = strings.NewReader(values.Encode())
	}
	if strings.ToLower(method) == "post" {
		req, err = http.NewRequest("POST", requestUrl, data)
		if err != nil {
			return nil, err
		}
	} else if strings.ToLower(method) == "get" {
		req, err = http.NewRequest("GET", requestUrl, nil)
		if err != nil {
			return nil, err
		}
		req.URL.RawQuery = values.Encode()
	} else {
		req, err = http.NewRequest(strings.ToUpper(method), requestUrl, data)
		if err != nil {
			return nil, err
		}
	}

	// Request Header
	setRequestHeaders(requestHeader, req)

	if req.Header.Get("User-Agent") == "" {
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; MAFSJS; rv:11.0) like Gecko")
	}
	if strings.ToLower(method) == "post" && req.Header.Get("Content-Type") == "" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	//req.SetBasicAuth("112233", "445566")

	// Cookie
	setCookies(cookies, req)

	// 自動リダイレクトのオフ、プロキシ設定
	// TODO: cookie jar
	client := http.Client{
		// Go HTTP Client NOT Follow Redirects Automatically.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			return nil, err
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}
	return client.Do(req)
}
