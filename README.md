# cwlogs-tee

`tee` command for CloudWatch Logs.

[![Build Status](https://travis-ci.org/winebarrel/cwlogs-tee.svg?branch=master)](https://travis-ci.org/winebarrel/cwlogs-tee)

## Usage

```
Usage of cwlogs-tee:
  -g string
      log group name
  -s string
      log stream name
  -v	show version
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

## Demo

[![asciicast](https://asciinema.org/a/81712.png)](https://asciinema.org/a/81712)

