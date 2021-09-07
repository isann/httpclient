# About
This is HTTP Client library.


# Getting started
Install go module.

```
go get github.com/isann/httpclient/v2
```

## Old version

### v1
```
go get github.com/isann/httpclient
```


# Usage
## GET
```
client := &httpclient.HttpClient{}
client.Url = "https://aaaaa:bbbbb@foobar/path/to/"
client.GetParameters = map[string]string{"a": "1", "b": "2"}
response, err := client.Get()
if err != nil {
    return
}
```

## POST（application/x-www-form-urlencoded）
```
client := &httpclient.HttpClient{}
client.Url = "https://aaaaa:bbbbb@foobar/path/to/"
client.Parameters = map[string]string{"a": "1", "b": "2"}
client.CookieJar = nil
client.Proxy = "http://localhost:8888/"
response, err := client.Post()
if err != nil {
    return
}
```

## POST (REST API, application/json)
```
m := map[string]interface{}{
    "param1": "value1",
    "param2": "value2",
    "param3": "value3",
}
decodeJson, err := json.Marshal(m)
if err != nil {
    return
}
client := &httpclient.HttpClient{}
client.Url = "https://aaaaa:bbbbb@foobar/path/to/"
client.Headers = map[string]string{
    "Content-Type":  "application/json",
    "Authorization": "Bearer xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
}
client.RawData = decodeJson
client.Proxy = ""
response, err := client.Post()
if err != nil {
    return
}
```

## POST（With file, multipart/form-data）

```
file, err := os.Open("/path/to/binary-file")
if err != nil {
    panic(err)
}
attachFiles := []AttachFile{{"file001", "00BFFF.jpg", file}}
client := &httpclient.HttpClient{}
client.Url = "https://aaaaa:bbbbb@foobar/path/to/"
client.Parameters = map[string]string{"a": "aaa",}
client.Files = attachFiles
client.CookieJar = nil
client.Proxy = ""
response, err := client.Post()
if err != nil {
    return
}
```

## CookieJar
```
// WIP
jar, _ := cookiejar.New(nil)
```
