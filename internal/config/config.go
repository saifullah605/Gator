package config

import ("os")
const configFileName = ".gatorconfig.json"

func GetConfigFilePath() (string, error){ 
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	return homeDir, nil

}

type Config struct {
	DBURL string `json:"db_url"`
	CurrUserName string `json:"current_user_name"`
}

func Read() (Config, error) {

	


	

	return Config{}, nil
}