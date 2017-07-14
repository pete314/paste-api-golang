//Author: Peter Nagy <https://peternagy.ie>
//Since: 06, 2017
//Description:
package common

import (
	"io/ioutil"

	"github.com/pquerna/ffjson/ffjson"
)

var RuntimeConfig *Config

//Config - the configuration structure
type Config struct {
	HTTPAddress        string   `json:"http-address"`
	HTTPPort           int      `json:"http-port"`
	HTTPSslEnable      bool     `json:"http-ssl-enable"`
	HTTPSslKey         string   `json:"http-ssl-key"`
	HTTPSslCert        string   `json:"http-ssl-cert"`
	RedisAddress       string   `json:"redis-address"`
	RedisPort          int      `json:"redis-port"`
	RedisDb            int      `json:"redis-db"`
	RedisPassword      string   `json:"redis-password"`
	AuthAccessExpires  int      `json:"auth-access-expires"`
	AuthRefreshExpires int      `json:"auth-refresh-expires"`
	AuthAccessBase     string   `json:"auth-access-base"`
	AuthRefreshBase    string   `json:"auth-refresh-base"`
	AuthClientID       string   `json:"auth-client-id"`
	AuthClientSecret   string   `json:"auth-client-secret"`
	Admins             []string `json:"admins"`
}

//LoadConfig - load configuration into conf structure
func LoadConfig(confPath string) *Config {
	var err error

	if fBytes, err := ioutil.ReadFile(confPath); err == nil {
		var conf Config
		if err := ffjson.Unmarshal(fBytes, &conf); err == nil {
			RuntimeConfig = &conf
			return &conf
		}
	}

	CheckError("Configuration loading error from path \n "+confPath, err, true)
	panic(err)
}
