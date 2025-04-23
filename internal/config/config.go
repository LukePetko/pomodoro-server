package config

import (
    "encoding/json"
    "os"
)

type Config struct {
    ShortBreakTime int `json:"short_break_time"`
    LongBreakTime int `json:"long_break_time"`
    WorkTime int `json:"work_time"`
    NumberOfSessions int `json:"number_of_sessions"`
}

func LoadConfig(path string) (*Config, error) {
    config := Config{}

    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&config); err != nil {
        return nil, err
    }
    return &config, nil
}

func SaveConfig(path string, config *Config) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    if err := encoder.Encode(config); err != nil {
        return err
    }

    return nil
}
