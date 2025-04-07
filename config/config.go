package config

import (
    "os"
    "fmt"
    "encoding/json"
)

type Config struct {
    Clients Clients `json:"clients"`
}

type Clients struct {
    Databases []Database `json:"databases"`
    Bucket Bucket `json:"bucket"`
}

type Database struct {
    Host string `json:"host"`
    Name string `json:"name"`
    Username string `json:"username"`
    Password string `json:"password"`
    Port string `json:"port"`
}

type Bucket struct {
    Name string `json:"name"`
}


func New(configPath string) (Config, error) {
    c, err := os.ReadFile(configPath)
    if err != nil {
        return Config{}, fmt.Errorf("error, when reading config file. Error: %v", err)
    }

    result := Config{}
    err = json.Unmarshal(c, &result)
    if err != nil {
        return Config{}, fmt.Errorf("error, when decoding config file. Error: %v", err)
    }
    return result, nil
}
