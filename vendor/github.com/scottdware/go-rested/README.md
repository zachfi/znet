## go-rested
[![GoDoc](https://godoc.org/github.com/scottdware/go-rested?status.svg)](https://godoc.org/github.com/scottdware/go-rested) [![Travis-CI](https://travis-ci.org/scottdware/go-rested.svg?branch=master)](https://travis-ci.org/scottdware/go-rested)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/scottdware/go-rested/master/LICENSE)

A Go package that makes calling RESTful API's easy.

### Usage & Examples

Initiate the request, and set authentication (if applicable):

```Go
r := rested.NewRequest()
r.BasicAuth("user", "secret")
```

If you need to specify any headers or query parameters, then use a `map[string]string` type to define them:

```Go
headers := map[string]string{
	"Content-Type": "application/json",
	"Accept": "application/json",
}

query := map[string]string{
	"search_string": "dog",
	"results": "10",
}
```

Then, call the `Send()` function to issue the request. This function takes the following parameters:

|Parameter|Description|
|-----|-----------|
|method|The HTTP method (i.e. GET, POST, DELETE, etc.)|
|uri|The URI/URL that you will be calling.|
|body|The contents of your request. This must be a byte slice (`[]byte`).|
|headers|Any additional headers you want to send. Must be a `map[string]string` type.|
|query|Any additional query parameters you want to send. Must be a `map[string]string` type.|

```Go
data := r.Send("get", "https://someurl/api/v1.0/stuff?default_param=something", nil, headers, query)
```

If you need to send/post a form, just place your form values in a `map[string]string` type and use the `SendForm()` function. The parameters are
the same as the `Send()` function, except in place of the `body` parameter, you have your form values.

> Note: headers and query parameters are the same as above.

```Go
formValues := map[string]string{
	"name": "Scott Ware",
	"age": "unknown",
	"sport": "Hockey",
}

data := r.SendForm("post", "https://someplace/to/upload", formValues, nil, nil)
```

> If there was any type of error in your request, it will be defined in the `Error` field of the returned struct. You can check for errors similar to how you normally do in Go:
```Go
if data.Error != nil {
	fmt.Println(data.Error)
}
```

The returned data is a struct with the following fields:

```Go
type Response struct {
	Status  string
	Code    int
	Headers http.Header
	Body    []byte
	Error   error
}
```

You can see the payload by converting the `Body` field to a string:

```Go
fmt.Println(string(data.Body))
```