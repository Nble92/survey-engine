package config

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	//"strconv"
)

// SurveyConfig basic config
type SurveyConfig struct {
	Database struct {
		Type string `json:"type"`
		Path string `json:"path"`
	} `json:"database"`
	Files struct {
		Path string `json:"path"`
	} `json:"files"`
	Server struct {
		ListenAddress string `json:"listenAddress"`
		ListenPort    int    `json:"listenPort"`
		JWTSigningKey string `json:"jwtSigningKey"`
	} `json:"server"`
	SendGrid struct {
		APIKey      string `json:"apiKey"`
		FromName    string `json:"fromName"`
		FromAddress string `json:"fromAddress"`
	} `json:"sendgrid"`
	Survey struct {
		APIURL       string `json:"APIURL"`
		Name         string `json:"Name"`
		FriendlyName string `json:"friendlyName"`
		WelcomeText  string `json:"welcomeText"`
		ContactName  string `json:"contactName"`
		ContactEmail string `json:"contactEmail"`
	} `json:"survey"`
}

// ConfigFile is the path to the JSON file, Config is the struct it is loaded in to
var (
	ConfigFile string
	Config     *SurveyConfig
)

// LoadConfig loads the file in to the struct
func LoadConfig(file string) (*SurveyConfig, error) {
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&Config)

	if err != nil {
		return nil, err
	}

	if Config.Server.ListenAddress == "" {
		Config.Server.ListenAddress = GetLocalIP()
	}

	// build API endpoint URL
	Config.Survey.APIURL = "/v1"

	return Config, err
}

// GetLocalIP is used to resolve the local IP address to use for the default ListenAddress
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
