package config

import (
  "encoding/json"
  "io/ioutil"
)

const config_file = "config.json"

type Config_t struct {
  Self string
  Contacts []string
  Backups []string
}

var Config *Config_t = Load();

func Load() *Config_t {
  data, err := ioutil.ReadFile(config_file)
  if err != nil {
    panic(err)
  }

  var config Config_t
  err = json.Unmarshal(data, &config)
  if err != nil {
    panic(err)
  }

  return &config
}

func (config Config_t) Save() error {
  config_b, err := json.MarshalIndent(config, "", "  ")
  if err != nil {
    return err
  }

  err = ioutil.WriteFile(config_file, config_b, 0644)
  if err != nil {
    return err
  }
  return nil
}
