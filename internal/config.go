package internal

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type TelegramConfig struct {
	Token    string        `yaml:"token"`
	Channel  int64         `yaml:"channel"`
	MsgDelay time.Duration `yaml:"message-delay"`
}

type WykopConfig struct {
	WykopUserKey   string   `yaml:"userkey"`
	TagBlacklist   []string `yaml:"tag-blacklist"`
	MaxPage        int      `yaml:"max-page"`
	MaxNoOfEntries int      `yaml:"max-no-of-entries"`
	MinVotes       int      `yaml:"min-votes"`
}

type Config struct {
	Telegram TelegramConfig `yaml:"telegram"`
	Wykop    WykopConfig    `yaml:"wykop"`

	BadgerDirpath string `yaml:"badger-dirpath"`
}

func LoadConfig(filepath string) (Config, error) {
	cfgFile, err := os.Open(filepath)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := yaml.NewDecoder(cfgFile).Decode(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) String() string {
	return fmt.Sprintf("badger dir: '%s', tag blacklist: '%+v', max page: %d, max no of entries: %d, min votes: %d", c.BadgerDirpath, c.Wykop.TagBlacklist, c.Wykop.MaxPage, c.Wykop.MaxNoOfEntries, c.Wykop.MinVotes)
}
