# ipapi

Query https://ip-api.com including rate limit handling. The base usage of this
library is for low request applications without a paid https://ip-api.com
account. It is though possible to configure the library for usage with a paid
account.

This is not an official library for https://ip-api.com.

## Installation

Install into your project using `go get github.com/cmol/ipapi`.

## Usage

The library works by returning a channel containing the result.

```go
ipapi.Run()
c, err := ipapi.Lookup("8.8.8.8")
if err != nil {
  // handle error
}
response := ipapi.Response{}
select {
  case response = <-c:
    // do something with the response object
  case default:
    // channel was closed, handle error
}
```

You can set the fields requested via `ipapi.Fields`. These fields can be
configured using https://ip-api.com/docs/api:json

The default is:

```go
var Fields = "?fields=status,message,country,countryCode,region,regionName,city,zip,lat,lon,timezone,isp,org,as,query"
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to
discuss what you would like to change.

Please make sure to update tests as appropriate.

## Similar libraries

A similar library exists at: https://pkg.go.dev/github.com/BenB196/ip-api-go-pkg


## License

[MIT](https://choosealicense.com/licenses/mit/)
