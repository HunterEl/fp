# CLI Tool to Proxy Function Calls

## Start:

- `go run main.go configure` and supply a config repo you have access to
- Check `fp.rc` to make sure `commandsRepo` is set to the correct repo

## Current Use:

- _DISCLAIMER_ None of the below commands work just yet (2019/03/15) since moving the test scripts to [fp-test-scripts](https://github.com/HunterEl/fp-test-scripts/)
- `go run main.go ls {additional args}`
- `go run main.go node-hello`
- `go run main.go go-hello`

## Tests:

- run the current test suite with `go test -v ./cmd/` in the root of project

## Config:

- `config.json`
- command is relative to the config file's location
- Schema:

```{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Config Schema",
  "description": "FP Config Schema",
  "type": "object",
  "required": ["commands"],
  "properties": {
    "commands": {
      "type": "object",
      "additionalProperties": {
        "type": "object",
        "required": ["command", "lang"],
        "properties": {
          "command": {
            "type": "string"
          },
          "environment": {
            "type": "string"
          },
          "lang": {
            "type": "string"
          },
          "runCommands": {
            "type": "array",
            "items": {
              "type": "string"
            }
          },
          "install": {
            "type": "string"
          }
        }
      }
    }
  }
}
```
