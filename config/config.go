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
}

type CredentialConfig struct {
	Username          string `yaml:"username"`
	Password          string `yaml:"password"`
	SecondaryPassword string `yaml:"secondary_password"`
}

type XpathConfig struct {
	LoginButton          string `yaml:"login_button"`
	SecondaryLoginButton string `yaml:"secondary_login_button"`
}

type IDConfig struct {
	Username          string `yaml:"username_box_id"`
	Password          string `yaml:"password_box_id"`
	SecondaryPassword string `yaml:"secondary_password_box_id"`
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
