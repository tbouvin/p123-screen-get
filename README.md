# Portfolio123 Screen Utility
This utility allows for automatically downloading running and downloading
screens from portfolio123.com

By default it retrieves the sets of screens to run for the day from the config.yml
file provided. It will attempt to merge those files together as a single
csv.

Screens are defined in sets in the config file. It is assumed that if 2 screens
are a part of the same set, then they contain the same tickers and sorted identically.
When the screens in the set are merged, their rows are merged together
Screens contained in a different set is placed in the rows below the first set of screens.
## Prerequisites
Download GoLang:
https://golang.org/dl/

#### Download chromedriver 89
https://chromedriver.chromium.org/downloads

#### Download Docker
https://www.docker.com/products/docker-desktop

#### Run selenium docker image
` docker run -d --restart unless-stopped -p 4444:4444 selenium/standalone-chrome:89.0`

## Clone Repo
git clone https://github.com/tbouvin/p123-screen-get.git

## Build config
Copy config from config/examples to resources/local/config.yml

Replace file paths, username, password, secondary password

## Run
cd p123-screen-get

go run cmd/main.go

## Help
go run cmd/main.go --help
### examples: 
#### Only download
go run cmd/main.go -convertfiles=false -mergefile=false
#### Only merge
go run cmd/main.go -getscreens=false -mergefile=false