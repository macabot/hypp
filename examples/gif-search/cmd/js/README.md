# Gif search JS

## Setup

```shell
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" public
GOOS=js GOARCH=wasm go build -ldflags="-X 'main.APIKey=<api_key>'" -o public/main.wasm main.go
```

where `<api_key>` is the [Giphy API](https://developers.giphy.com/) key.

Start a server in the `public` directory. E.g.:

```shell
python3 -m http.server
```
