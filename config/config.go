package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Host     string `json:"db_host"`
	User     string `json:"db_user"`
	Password string `json:"db_password"`
	DB_Name     string `json:"db_name"`
}

type AccountConfig struct {
	Role string `json:"super"`
	Account string `json:"acc_super"`
	Password string `json:"pw_super"`
}

func LoadConfig(path string) (*Config, error) {
	var Cfg Config
	bts, err := ioutil.ReadFile(path)
	if err != nil {
		return &Cfg, err
	}

	if err := json.Unmarshal(bts, &Cfg); err != nil {
		return &Cfg, err
	}
	return &Cfg, nil

}

func LoadConfigAccount (path string) (*AccountConfig,error){
	var Cfg AccountConfig
	bts, err := ioutil.ReadFile(path)
	if err != nil {
		return &Cfg, err
	}

	if err := json.Unmarshal(bts, &Cfg); err != nil {
		return &Cfg, err
	}
	return &Cfg, nil
}
