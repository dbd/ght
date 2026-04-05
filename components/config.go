package components

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
)

type Config struct {
	Repo  string      `mapstructure:"repo"`
	Pr    PrConfig    `mapstructure:"pr"`
	Issue IssueConfig `mapstructure:"issue"`
}

type PrConfig struct {
	Searches []Search `mapstructure:"searches"`
}

type IssueConfig struct {
	Searches   []Search        `mapstructure:"searches"`
	Milestones []MilestoneRepo `mapstructure:"milestones"`
}

type MilestoneRepo struct {
	Name string `mapstructure:"name"`
	Repo string `mapstructure:"repo"`
}

type Search struct {
	Name  string `mapstructure:"name"`
	Query string `mapstructure:"query"`
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

func SaveIssueSearch(name string, query string) error {
	config := GetConfig()

	for i, search := range config.Issue.Searches {
		if search.Name == name {
			config.Issue.Searches[i].Query = query
			return writeConfig(config)
		}
	}

	config.Issue.Searches = append(config.Issue.Searches, Search{Name: name, Query: query})
	return writeConfig(config)
}

func SaveMilestoneRepo(name string, repo string) error {
	config := GetConfig()

	for i, m := range config.Issue.Milestones {
		if m.Name == name {
			config.Issue.Milestones[i].Repo = repo
			return writeConfig(config)
		}
	}

	config.Issue.Milestones = append(config.Issue.Milestones, MilestoneRepo{Name: name, Repo: repo})
	return writeConfig(config)
}

func writeConfig(config Config) error {
	viper.Set("pr.searches", config.Pr.Searches)
	viper.Set("issue.searches", config.Issue.Searches)
	viper.Set("issue.milestones", config.Issue.Milestones)
	configPath := os.ExpandEnv("$HOME/.config/ght/config.yaml")
	
	// Create directory if it doesn't exist
	configDir := os.ExpandEnv("$HOME/.config/ght")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}
	
	return viper.WriteConfigAs(configPath)
}
