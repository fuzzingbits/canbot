# CanBot

A tool that watches and send alerts when users are deleted from your instance of Slack. This can be helpful in large (or any sized) organizations to know when users are no longer employed.

[![GitHub Actions](https://github.com/fuzzingbits/canbot/workflows/Go/badge.svg)](https://github.com/fuzzingbits/canbot/actions)
[![Coverage Status](https://coveralls.io/repos/github/fuzzingbits/canbot/badge.svg?branch=main)](https://coveralls.io/github/fuzzingbits/canbot?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/fuzzingbits/canbot)](https://goreportcard.com/report/github.com/fuzzingbits/canbot)
[![License](https://img.shields.io/github/license/fuzzingbits/canbot)](https://github.com/fuzzingbits/canbot/blob/main/LICENSE)
[![Docker Image](https://img.shields.io/badge/container-Docker-blue)](https://hub.docker.com/r/fuzzingbits/canbot)


## Config Options
| Environment Variable | Default | Required | Description |
| -------------------- | ------- | -------- | ----------- |
| SLACK_TOKEN | `n/a` | `true` | Your Slack API Token (Scopes needed: `users:read`, `chat:write:bot`) |
| SLACK_TARGETS | `n/a` | `true` | A comma separated list of Slack IDs to send alerts to. Example: `C029XH96S,C50FW7CER,D2X7AC3QR` |
| SLACK_USERNAME | `CanBot` | `false` | The username used to send the Slack message |
| SLACK_ICON_EMOJI | `:flushed:` | `false` | The emoji used as the profile picture of the Slack message |
| STATE_FILE | `state.json` | `false` | The filename of the state json file |

## Running Locally
- First time: run `make`
    - This will download dependencies and make sure everything is working
- To test changes: run `make dev`
    - This will just build and run your changes
    - Will attempt to read from the `.env` file for config options
- To build the docker image: run `make docker`
