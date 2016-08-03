# cwlogs-tee

`tee` command for CloudWatch Logs.

## Usage

```
Usage of cwlogs-tee:
  -g string
      log group name
  -s string
      log stream name
```

```sh
while true; do
  date
  sleep 1
done | cwlogs-tee -g my-group -s my-stream

# LogGroup/LogStream is created automatically
```

## Installation

```
brew install https://raw.githubusercontent.com/winebarrel/cwlogs-tee/master/homebrew/cwlogs-tee.rb
```
