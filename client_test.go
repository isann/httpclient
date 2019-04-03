package httpclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestRequestHttpWithFile(t *testing.T) {
	type args struct {
		requestUrl string
		files      []AttachFile
		postParam  map[string]string
		proxy      string
	}
	file, err := os.Open("/Users/zono/Desktop/00BFFF.jpg")
	if err != nil {
		panic(err)
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr bool
	}{
		{"", args{"https://112233:445566@www.zono.xyz/", []AttachFile{{"", "00BFFF.jpg", file}}, map[string]string{"aaa": "bbb"}, "http://localhost:8888/"}, &http.Response{StatusCode: 200}, false},
		{"", args{"https://112233:445566@www.zono.xyz/", []AttachFile{{"", "00BFFF.jpg", file}}, map[string]string{"aaa": "bbb"}, "http://localhost:8000/"}, &http.Response{StatusCode: 200}, true},
		{"", args{"https://112233:445566@www.zono.xyz/", []AttachFile{{"", "00BFFF.jpg", file}}, map[string]string{"aaa": "bbb"}, ""}, &http.Response{StatusCode: 200}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RequestHttpWithFile(tt.args.requestUrl, "", tt.args.postParam, nil, nil, tt.args.files, nil, tt.args.proxy)
			if (err != nil) != tt.wantErr {
				t.Errorf("RequestHttpWithFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				return
			}
			if got.StatusCode != tt.want.StatusCode {
				t.Errorf("RequestHttpWithFile() = %v, want %v", got, tt.want)
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("RequestHttpWithFile() = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestRequestHttp(t *testing.T) {
	type args struct {
		requestUrl    string
		method        string
		parameters    map[string]string
		cookies       []*http.Cookie
		requestHeader map[string]string
		rawData       []byte
		proxy         string
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RequestHttp(tt.args.requestUrl, tt.args.method, tt.args.parameters, tt.args.cookies, tt.args.requestHeader, tt.args.rawData, nil, tt.args.proxy)
			if (err != nil) != tt.wantErr {
				t.Errorf("RequestHttp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RequestHttp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPost(t *testing.T) {
	type args struct {
		requestUrl    string
		parameters    map[string]string
		cookies       []*http.Cookie
		requestHeader map[string]string
		rawData       []byte
		files         []AttachFile
		proxy         string
	}
	feature := time.Now()
	feature = feature.Add(10 * 24 * time.Hour)
	jar, _ := cookiejar.New(nil)
	cookies := []*http.Cookie{
		{Name: "abc", Value: "dedede", Expires: feature, Path: "/", Domain: ".www.zono.xyz"},
	}
	requestUrl := "https://www.zono.xyz/public/set-cookie/index.php"
	setCookieUrl, _ := url.Parse(requestUrl)
	jar.SetCookies(setCookieUrl, cookies)
	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr bool
	}{
		{"", args{requestUrl, map[string]string{"aaa": "bbb"}, cookies, nil, nil, nil, "http://localhost:8888/"}, &http.Response{StatusCode: 200}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Post(tt.args.requestUrl, tt.args.parameters, nil, tt.args.requestHeader, tt.args.rawData, tt.args.files, jar, tt.args.proxy)
			if (err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				_cookies := got.Cookies()
				fmt.Printf("**************************** %v\n", _cookies)
				body := got.Body
				bytes, _ := ioutil.ReadAll(body)
				println(string(bytes))
			}
			println(len(tt.args.cookies))
			cookies := jar.Cookies(setCookieUrl)
			fmt.Printf("%v\n", cookies)
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("Post() = %v, want %v", got, tt.want)
			//}
		})
	}
}
