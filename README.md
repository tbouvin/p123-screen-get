# Portfolio123 Screen Utility
## Prerequisites
Download GoLang:
https://golang.org/dl/

#### (macOS) Download Homebrew: 
* Copy following in terminal and run:

/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
#### (macOS) Download Selenium: 
* brew install selenium-server-standalone

## Clone Repo
git clone https://github.com/tbouvin/p123-screen-get.git

## Build config
Copy config from config/examples to resources/local/config.yml

Replace file paths, username, password, secondary password

If running Windows, replace selenium command/arguments

## Run
cd p123-screen-get
go run cmd/main.go
