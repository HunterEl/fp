# CLI Tool to Proxy Function Calls

## Start:

- `go run main.go configure` and supply a config repo you have access to
- Check `fp.rc` to make sure `commandsRepo` is set to the correct repo

## Current Use:

- `go run main.go ls {additional args}`
- `go run main.go node-hello`
- `go run main.go go-hello`

## Config:

- `config.json`
- command is relative to the config file's location
