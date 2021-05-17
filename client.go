package httpclient

import (
	"bytes"
	"crypto/tls"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

type HttpClient struct {
	Url                string
	Parameters         Parameters
	Cookies            []*http.Cookie
	Headers            RequestHeader
	RawData            []byte
	Files              []AttachFile
	CookieJar          http.CookieJar
	Proxy              string
	InsecureSkipVerify bool
}

func (c *HttpClient) Clear() {
	c.Url = ""
	c.Parameters = nil
	c.Cookies = nil
	c.Headers = nil
	c.RawData = nil
	c.Files = []AttachFile{}
	c.CookieJar = nil
	c.Proxy = ""
}

// Post は HTTP POST で送信します。 parameters と rawData は互いに排他的で、 rawData が優先されます。
// rawData は、「application/x-www-form-urlencoded」以外でリクエストする際に指定します。
// ファイルを添付した場合は、「multipart/form-data」 となり、 rawData は参照されません。
func (c *HttpClient) Post() (*http.Response, error) {
	if c.Files == nil || len(c.Files) == 0 {
		return c.RequestHTTP(c.Url, "POST", c.Parameters, c.Cookies, c.Headers, c.RawData, c.CookieJar, c.Proxy)
	} else {
		return c.RequestHTTPWithFile(c.Url, "POST", c.Parameters, c.Cookies, c.Headers, c.Files, c.CookieJar, c.Proxy)
	}
}

// Get は HTTP GET で送信します。  rawData は無視され、parameters の値が query string に変換されます。
func (c *HttpClient) Get() (*http.Response, error) {
	return c.RequestHTTP(c.Url, "GET", c.Parameters, c.Cookies, c.Headers, c.RawData, c.CookieJar, c.Proxy)
}

// Put は HTTP PUT で送信します。 parameters と rawData は互いに排他的で、 rawData が優先されます。
// rawData は、「application/x-www-form-urlencoded」以外でリクエストする際に指定します。
func (c *HttpClient) Put() (*http.Response, error) {
	return c.RequestHTTP(c.Url, "PUT", c.Parameters, c.Cookies, c.Headers, c.RawData, c.CookieJar, c.Proxy)
}

// Delete は HTTP DELETE で送信します。 parameters と rawData は互いに排他的で、 rawData が優先されます。
// rawData は、「application/x-www-form-urlencoded」以外でリクエストする際に指定します。
func (c *HttpClient) Delete() (*http.Response, error) {
	return c.RequestHTTP(c.Url, "DELETE", c.Parameters, c.Cookies, c.Headers, c.RawData, c.CookieJar, c.Proxy)
}

type Parameters map[string]string

// TODO: GET と POST のパラメータは分離する
//type GetParameters map[string]string
//type PostParameters map[string]string

type RequestHeader map[string]string

// AttachFile は HTTP リクエストに添付するファイルです。
type AttachFile struct {
	FieldName string
	FileName  string
	Reader    io.Reader
}

func (c *HttpClient) setRequestHeaders(requestHeader RequestHeader, req *http.Request) {
	if requestHeader != nil {
		for key, val := range requestHeader {
			req.Header.Add(key, val)
		}
	}
}

func (c *HttpClient) setCookies(cookies []*http.Cookie, req *http.Request, cookieJar http.CookieJar) {
	if cookies != nil {
		if cookieJar != nil {
			cookieJar.SetCookies(req.URL, cookies)
		} else {
			for _, cookie := range cookies {
				req.AddCookie(cookie)
			}
		}
	}
}

func (c *HttpClient) setDefaultRequestHeaders(req *http.Request, method string) {
	if req.Header.Get("User-Agent") == "" {
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; MAFSJS; rv:11.0) like Gecko")
	}
	if strings.ToLower(method) == "post" && req.Header.Get("Content-Type") == "" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	//req.SetBasicAuth("112233", "445566")
}

func (c *HttpClient) setupHttpClient(cookieJar http.CookieJar, proxy string) (http.Client, error) {
	// 自動リダイレクトのオフ
	client := http.Client{
		// Go HTTP Client NOT Follow Redirects Automatically.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	// プロキシ設定
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return http.Client{}, err
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}
	// SSL 証明書チェック
	if c.InsecureSkipVerify {
		c := &tls.Config{InsecureSkipVerify: true}
		if client.Transport == nil {
			client.Transport = &http.Transport{
				TLSClientConfig: c,
			}
		} else {
			transport := client.Transport.(*http.Transport)
			transport.TLSClientConfig = c
		}
	}
	// Cookie 設定
	if cookieJar != nil {
		client.Jar = cookieJar
	}
	return client, nil
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

// RequestHTTPWithFile はマルチパートフォームデータ「multipart/form-data」で送信します。画像をリクエストする場合などに使用します。
// parameters と rawData は互いに排他的で、 rawData が優先されます。
// rawData は、「application/x-www-form-urlencoded」以外でリクエストする際に指定します。
func (c *HttpClient) RequestHTTPWithFile(requestURL string, method string, parameters Parameters, cookies []*http.Cookie, requestHeader RequestHeader, files []AttachFile, cookieJar http.CookieJar, proxy string) (*http.Response, error) {
	var b bytes.Buffer
	var fw io.Writer
	var err error
	w := multipart.NewWriter(&b)

	// Add file
	if files != nil {
		for _, v := range files {
			// TODO: CreateFormFile ではなく CreatePart で file の content-type を octet-stream から変更できるようにする
			//h := make(textproto.MIMEHeader)
			//h.Set("Content-Disposition",
			//	fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			//		escapeQuotes(v.FieldName), escapeQuotes(v.FileName)))
			//h.Set("Content-Type", "image/jpeg")
			//fw, err = w.CreatePart(h)
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
	req, err := http.NewRequest(method, requestURL, &b)
	if err != nil {
		return nil, err
	}

	// Add Content-Type multipart/form-data header
	req.Header.Add("Content-Type", w.FormDataContentType())

	// Request Header
	c.setRequestHeaders(requestHeader, req)
	c.setDefaultRequestHeaders(req, method)
	// Cookie
	c.setCookies(cookies, req, cookieJar)
	client, err := c.setupHttpClient(cookieJar, proxy)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

// RequestHTTP は、リクエストヘッダ Content-Type の指定がない場合は「application/x-www-form-urlencoded」で送信します。
// フォームを送信する場合などに使用します。
// parameters と rawData は互いに排他的で、 rawData が優先されます。
// rawData は、「application/x-www-form-urlencoded」以外でリクエストする際に指定します。
func (c *HttpClient) RequestHTTP(requestURL string, method string, parameters Parameters, cookies []*http.Cookie, requestHeader RequestHeader, rawData []byte, cookieJar http.CookieJar, proxy string) (*http.Response, error) {
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
		req, err = http.NewRequest("POST", requestURL, data)
		if err != nil {
			return nil, err
		}
	} else if strings.ToLower(method) == "get" {
		req, err = http.NewRequest("GET", requestURL, nil)
		if err != nil {
			return nil, err
		}
		req.URL.RawQuery = values.Encode()
	} else {
		req, err = http.NewRequest(strings.ToUpper(method), requestURL, data)
		if err != nil {
			return nil, err
		}
	}

	// Request Header
	c.setRequestHeaders(requestHeader, req)
	c.setDefaultRequestHeaders(req, method)
	// Cookie
	c.setCookies(cookies, req, cookieJar)
	client, err := c.setupHttpClient(cookieJar, proxy)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}
