package main

import (
	"io/ioutil"
	"strconv"
	"strings"
)

type ServerConfig struct {
	ServerName       string
	SecureServerIP   string
	SecureServerPort string
	ServerPort       string
	PrudpVersion     int
	SignatureVersion int
	KerberosKeySize  int
	AccessKey        string
	MongoAddress     string
}

func ImportConfigFromFile(path string) (*ServerConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	indexes := strings.Split(string(data), "\r\n")
	config := &ServerConfig{
		ServerName:      "server",
		KerberosKeySize: 32}
	for i := 0; i < len(indexes); i++ {
		index := strings.Split(indexes[i], "=")
		if len(index) != 2 {
			continue
		}
		switch index[0] {
		case "ServerName":
			config.ServerName = index[1]
			break
		case "SecureServerIP":
			config.SecureServerIP = index[1]
			break
		case "SecureServerPort":
			config.SecureServerPort = index[1]
			break
		case "ServerPort":
			config.ServerPort = index[1]
			break
		case "PrudpVersion":
			config.PrudpVersion, err = strconv.Atoi(index[1])
			if err != nil {
				return nil, err
			}
			break
		case "SignatureVersion":
			config.SignatureVersion, err = strconv.Atoi(index[1])
			if err != nil {
				return nil, err
			}
			break
		case "KerberosKeySize":
			config.KerberosKeySize, err = strconv.Atoi(index[1])
			if err != nil {
				return nil, err
			}
			break
		case "AccessKey":
			config.AccessKey = index[1]
			break
		case "MongoAddress":
			config.MongoAddress = index[1]
			break
		default:
			break
		}
	}
	return config, nil
}
