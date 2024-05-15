# Memoized list JS

## Setup

```shell
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" public
GOOS=js GOARCH=wasm go build -o public/main.wasm main.go
```

Start a server in the `public` directory. E.g.:

```shell
python3 -m http.server
```
