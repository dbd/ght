package components

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
)

type Config struct {
	Repo  string   `mapstructure:Repo`
	Pr    PrConfig `mapstructure:pr`
	Issue IssueConfig
}

type PrConfig struct {
	Searches []Search `mapstructure:searches`
}
type IssueConfig struct {
	Searches []Search
}

type Search struct {
	Name  string `mapstructure:name`
	Query string `mapstructure:query`
}

func GetConfigCmd() tea.Cmd {
	return func() tea.Msg {
		return GetConfig()
	}
}

func GetInitialConfig() tea.Cmd {
	return func() tea.Msg {
		return SetupConfig()
	}
}

func SetupConfig() Config {
	var c Config
	viper.AddConfigPath("$HOME/.config/ght")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()
	err := viper.Unmarshal(&c)
	if err != nil {
		fmt.Printf("unable to decode into struct, %v", err)
	}
	return c
}

func GetConfig() Config {
	var c Config
	viper.ReadInConfig()
	err := viper.Unmarshal(&c)
	if err != nil {
		fmt.Printf("unable to decode into struct, %v", err)
	}
	return c
}

func SaveSearch(name string, query string) error {
	config := GetConfig()
	
	// Check if search already exists
	for i, search := range config.Pr.Searches {
		if search.Name == name {
			config.Pr.Searches[i].Query = query
			return writeConfig(config)
		}
	}
	
	// Add new search
	config.Pr.Searches = append(config.Pr.Searches, Search{Name: name, Query: query})
	return writeConfig(config)
}

func writeConfig(config Config) error {
	viper.Set("pr.searches", config.Pr.Searches)
	configPath := os.ExpandEnv("$HOME/.config/ght/config.yaml")
	return viper.WriteConfigAs(configPath)
}
