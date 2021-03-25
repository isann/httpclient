# About
This is HTTP Client library.


# Getting start
Install go module.

```
go get github.com/isann/httpclient
```

# Usage
## GET
```
// WIP
```

## POST（application/x-www-form-urlencoded）
```
// WIP
```

## POST (REST API, application/json)
```
m := map[string]interface{}{
    "param1": "value1",
    "param2": "value2",
}
decodeJson, err := json.Marshal(m)
if err != nil {
    return
}
_, err = httpclient.RequestHTTP(
    "https://foobar/user",
    "POST",
    nil,
    []*http.Cookie{},
    map[string]string{
        "Content-Type":  "application/json",
        "Authorization": "Bearer xxxxxxxxxxxx",
    },
    decodeJson,
    nil,
    "")
if err != nil {
    return
}
```

## POST（With file, multipart/form-data）

```
file, err := os.Open("/path/to/binary-file")
if err != nil {
    return
}
_, err = httpclient.RequestHTTPWithFile(
    "https://foobar/user/icon",
    "POST",
    map[string]string{"param1": "value1", "param2": "value2", "param3": "value3"},
    []*http.Cookie{},
    map[string]string{},
    []httpclient.AttachFile{{"file", "foobar.jpg", file}},
    nil,
    "")
if err != nil {
    return
}
```

## CookieJar
```
// WIP
jar, _ := cookiejar.New(nil)
```
