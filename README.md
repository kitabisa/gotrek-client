# gotrek-client

Client library for [GoTrek - Audit Trail Service](https://github.com/kitabisa/gotrek).

## Getting started

This project requires Go to be installed. On OS X with Homebrew you can just run `brew install go`. Then run `go get ./...` from project dir to fetch all dependencies. Don't forget to rename `gotrek.toml.sample` to `gotrek.toml` and modify it based on your need.

## Using the library

- `import github.com/kitabisa/gotrek-client`
- Define your config for HTTP Client, 
```console
    httpClientConf := &TrekHttpClient{
		Timeout:               1 * time.Second,
		BackoffInterval:       5 * time.Microsecond,
		MaximumJitterInterval: 5 * time.Microsecond,
		RetryCount:            4,
	}
```
- Set GoTrek host, `url := "http://some-url/vx"` , where `x` is API version based on `gotrek-client` major version number. 
ex: using `gotrek-client` `v1.0.0` you should set the value to `"http://some-url/v1"`  
- Set GoTrek secret key, `secret := "really-secret"`
- Init client, `client := NewTrekClient(url, secret, httpClientConf)`
- Start publish event, `client.publish("some-id", someMapInterface, time.Now().Unix(), "TAG,WITH_TAG")`