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
    Database Database `json:"database"`
    Bucket Bucket `json:"bucket"`
}

type Database struct {
    Host string `json:"host"`
    Username string `json:"username"`
    Password string `json:"password"`
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
