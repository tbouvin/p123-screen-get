# Portfolio123 Screen Utility
## Prerequisites
Download GoLang:
https://golang.org/dl/

#### (macOS) Download Homebrew: 
* /usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
#### (macOS) Download Selenium: 
* brew install selenium-server-standalone

## Clone Repo
git clone https://github.com/tbouvin/p123-screen-get.git

## Build
cd p123-screen-get
go build -o p123screenget cmd/main.go

## Build config
Copy config from config/examples to resources/local/config.yml

Replace file paths, username, password, secondary password

## Run
./p123screenget