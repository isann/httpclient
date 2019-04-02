package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func RequestHttpWithFile(requestUrl string, files []AttachFile, postParam map[string]string) (*http.Response, error) {
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
	if postParam != nil {
		for key, val := range postParam {
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
	req, err := http.NewRequest("POST", requestUrl, &b)
	if err != nil {
		fmt.Println("requestMsgRegister", err)
		return nil, err
	}
	req.Header.Add("Content-Type", w.FormDataContentType())
	//req.SetBasicAuth("112233", "445566")
	// 自動リダイレクトのオフ、プロキシ設定
	proxyUrl, err := url.Parse("http://localhost:8888")
	client := http.Client{
		// Go HTTP Client NOT Follow Redirects Automatically.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
	}
	return client.Do(req)
}

func requestHttp(requestUrl string, method string, getParam map[string]string,
	postParam map[string]string, cookie *http.Cookie, isRaw bool,
	requestHeader map[string]string) (*http.Response, error) {
	var req *http.Request
	var err error
	if strings.ToLower(method) == "post" {
		if isRaw {
			postData, err := json.Marshal(postParam)
			if err != nil {
				return nil, err
			}
			req, err = http.NewRequest("POST", requestUrl, strings.NewReader(string(postData)))
		} else {
			values := url.Values{}
			if postParam != nil {
				for key, val := range postParam {
					values.Add(key, val)
				}
			}
			req, err = http.NewRequest("POST", requestUrl, strings.NewReader(values.Encode()))
		}
		if err != nil {
			return nil, err
		}
	} else if strings.ToLower(method) == "get" {
		req, err = http.NewRequest("GET", requestUrl, nil)
		if err != nil {
			return nil, err
		}
	}
	// GET parameter
	if getParam != nil {
		values2 := url.Values{}
		for key, val := range getParam {
			values2.Add(key, val)
		}
		req.URL.RawQuery = values2.Encode()
	}
	// Request Header
	if requestHeader != nil {
		for key, val := range requestHeader {
			req.Header.Add(key, val)
		}
	}
	if req.Header.Get("User-Agent") == "" {
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; MAFSJS; rv:11.0) like Gecko")
	}
	if strings.ToLower(method) == "post" && req.Header.Get("Content-Type") == "" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	//req.SetBasicAuth("112233", "445566")
	if cookie != nil {
		req.AddCookie(cookie)
	}
	// 自動リダイレクトのオフ、プロキシ設定
	//proxyUrl, err := url.Parse("http://192.168.20.177:8888")
	client := http.Client{
		// Go HTTP Client NOT Follow Redirects Automatically.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		//Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
	}
	return client.Do(req)
}
