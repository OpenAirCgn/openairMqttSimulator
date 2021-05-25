package openairMqttSimulator

import (
	"encoding/json"
	"os"
)

type ClientCertConfig struct {
	ClientId    string `json:"clientID"`
	CRTFilename string `json:"crt"`
	PEMFilename string `json:"pem"`
	MAC         string `json:"mac"`
	// typically, the clientId will be the hashed mac
	// we can use arbitrary names loading the list,
	// but may want to keep track of macs.
}

func LoadClientCertConfigList(fn string) (list []ClientCertConfig, err error) {
	file, err := os.Open(fn)
	if err != nil {
		return list, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&list)
	return
}
