package Config

import (
	"os"
	"path/filepath"
	"log"
)

var Config = map[string]string{
	"DB_USER": "nanollet",
	"DB_PASS": "",
	"DB_IP":   "127.0.0.1:3306",
	"DB_KEY":  "",

	"API_IP":      "127.0.0.1",
	"CALLBACK_IP": "127.0.0.1",

	"LOG": "",

	"DOMAIN_NAME":       "",
	"DOMAIN_PUBKEYHASH": "",
}

func Start() {
	for k := range Config {
		if conf := os.Getenv("NANOLLET_" + k); conf != "" {
			Config[k] = conf
		}
	}
}

func Dir() string {
	ex, err := os.Executable()
	if err != nil {
		log.Panic(err)
	}
	return filepath.Dir(ex)
}
