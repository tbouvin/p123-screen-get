package config

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Credentials CredentialConfig `yaml:"credentials"`
	Xpaths      XpathConfig      `yaml:"xpaths"`
	IDs         IDConfig         `yaml:"ids"`
	Screens     ScreenConfig     `yaml:"screens"`
	URLs        URLConfig        `yaml:"urls"`
	FilePaths   FilePathsConfig  `yaml:"file_paths"`
	Selenium    SeleniumConfig   `yaml:"selenium"`
	SleepTime   int              `yaml:"sleep_time"`
}

type FilePathsConfig struct {
	DownloadPath string `yaml:"download_path"`
	CSVPath      string `yaml:"csv_path"`
	CombinedPath string `yaml:"combined_path"`
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
	MultipleLoginButton  string `yaml:"multiple_login_button"`
}

type IDConfig struct {
	Username               string `yaml:"username_box_id"`
	Password               string `yaml:"password_box_id"`
	SecondaryPassword      string `yaml:"secondary_password_box_id"`
	RunScreenButton        string `yaml:"run_screen_button"`
	TickerLink             string `yaml:"ticker_link"`
	SkipMultipleLoginCheck string `yaml:"skip_multiple_login_check"`
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
	Base   string `yaml:"base`
}

type SeleniumConfig struct {
	Port      string   `yaml:"port"`
	Command   string   `yaml:"command"`
	Arguments []string `yaml:"arguments"`
}

func GetConfig(configFile *string) (Config, error) {
	path, err := filepath.Abs(*configFile)
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return Config{}, err
	}

	var conf Config
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return Config{}, err
	}

	return conf, nil
}
