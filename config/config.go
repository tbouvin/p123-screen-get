package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Credentials CredentialConfig `yaml:"credentials"`
	Xpaths      XpathConfig      `yaml:"xpaths"`
	IDs         IDConfig         `yaml:"ids"`
	Screens     ScreenConfig     `yaml:"screens"`
	URLs        URLConfig        `yaml:"urls"`
}

type CredentialConfig struct {
	Username          string `yaml:"username"`
	Password          string `yaml:"password"`
	SecondaryPassword string `yaml:"secondary_password"`
}

type XpathConfig struct {
	LoginButton          string `yaml:"login_button"`
	SecondaryLoginButton string `yaml:"secondary_login_button"`
	ShowAllScreenButton  string `yaml:"show_all_screen_button"`
	ScreenDownload       string `yaml:"screen_download"`
}

type IDConfig struct {
	Username          string `yaml:"username_box_id"`
	Password          string `yaml:"password_box_id"`
	SecondaryPassword string `yaml:"secondary_password_box_id"`
}

type ScreenConfig struct {
	Monday    []ScreenPart `yaml:"monday"`
	Tuesday   []ScreenPart `yaml:"tuesday"`
	Wednesday []ScreenPart `yaml:"wednesday"`
	Thursday  []ScreenPart `yaml:"thursday"`
	Friday    []ScreenPart `yaml:"friday"`
	Saturday  []ScreenPart `yaml:"saturday"`
	Sunday    []ScreenPart `yaml:"sunday"`
}

type ScreenPart struct {
	Names []string `yaml:"names"`
}

type URLConfig struct {
	Login  string `yaml:"login"`
	Screen string `yaml:"screen"`
}

func GetConfig() Config {
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return Config{}
	}

	var conf Config
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return Config{}
	}

	return conf
}
