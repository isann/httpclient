package httpclient

import (
	"net/http"
	"os"
	"reflect"
	"testing"
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
			got, err := RequestHttpWithFile(tt.args.requestUrl, tt.args.files, tt.args.postParam, tt.args.proxy)
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
			got, err := RequestHttp(tt.args.requestUrl, tt.args.method, tt.args.parameters, tt.args.cookies, tt.args.requestHeader, tt.args.rawData, tt.args.proxy)
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
