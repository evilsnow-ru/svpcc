package config

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
)

const configFileName = "config.json"

type Config interface {
	ServerAddress() string
}

type ConfigImpl struct {
	Server string	`json:"server"`
	Port uint16		`json:"port"`
}

func (cfg *ConfigImpl) ServerAddress() string {
	return "http://" + cfg.Server + ":" + strconv.FormatUint(uint64(cfg.Port), 10)
}

func GetConfig() Config {
	_, err := os.Stat(configFileName)

	if err != nil {
		if os.IsNotExist(err) {
			cfgImpl := ConfigImpl{
				Server: "127.0.0.1",
				Port: 9080,
			}

			configFile, err := os.Create(configFileName)

			if err != nil {
				log.Panicf("Error creating file \"%s\": %s.", configFileName, err.Error())
			}

			defer configFile.Close()
			err = json.NewEncoder(configFile).Encode(cfgImpl)

			if err != nil {
				log.Panicf("Can't encode config structure to json: %s.", err.Error())
			}

			return &cfgImpl
		} else {
			log.Panicf("Error calling stat on file \"%s\": %s.", configFileName, err.Error())
		}
	}

	configFile, err := os.Open(configFileName)

	if err != nil {
		log.Panicf("Can't open \"%s\": %s.", configFileName, err.Error())
	}

	defer configFile.Close()

	cfg := ConfigImpl{}
	err = json.NewDecoder(configFile).Decode(&cfg)

	if err != nil {
		panic(err.Error())
	}

	return &cfg
}