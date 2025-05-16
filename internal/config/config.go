package config

import (
    "encoding/json"
    "os"
    "fmt"
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
    base, err := LoadConfig(path)
    fmt.Println(base)

    if err != nil {
        return err
    }

    file, err := os.Create(path)

    if err != nil {
        return err
    }
    defer file.Close()

    base = PatchConfig(base, config)
    fmt.Println(base)

    encoder := json.NewEncoder(file)
    if err := encoder.Encode(base); err != nil {
        return err
    }

    return nil
}

func PatchConfig(base *Config, partial *Config) *Config {
    if (partial.ShortBreakTime != 0) {
        base.ShortBreakTime = partial.ShortBreakTime
    }

    if (partial.LongBreakTime != 0) {
        base.LongBreakTime = partial.LongBreakTime
    }

    if (partial.WorkTime != 0) {
        base.WorkTime = partial.WorkTime
    }

    if (partial.NumberOfSessions != 0) {
        base.NumberOfSessions = partial.NumberOfSessions
    }

    return base 
}
