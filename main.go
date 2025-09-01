package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
)

type Config struct {
	AppName string `mapstructure:"app_name"`
	Logging struct {
		StructuredLogging bool   `mapstructure:"structured-logging`
		DebugLevel        string `mapstructure:"debuglevel"`
	} `mapstructure:"logging"`
	Service struct {
		Port  int `mapstructure:"port"`
		Debug int `mapstructure:"debug"`
	} `mapstructure:"service"`
	DB struct {
		Variant  string `mapstructure:"variant"`
		Username string `mapstructure:"username"`
		Dbname   string `mapstructure:"dbname"`
	} `mapstructure:"db"`
}

type model struct {
	recordTransaction bool
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	fmt.Printf("Proceeding with input : %v", msg)
	return m, nil
}
func (m model) View() string {
	return "Record Transaction? (y/n)"
}

func (m model) Init() tea.Cmd {
	return nil
}

func loadConfig() *Config {
	viper.AddConfigPath("./config")
	viper.SetConfigFile("./config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("unable to read in configuration: %v", err)
	}
	credentialsScope := viper.New()
	credentialsScope.SetConfigFile("./config/credentials.yaml")
	if err := credentialsScope.ReadInConfig(); err != nil {
		log.Fatalf("unable to read in configuration: %v", err)
	}
	viper.MergeConfigMap(credentialsScope.AllSettings())
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("unable to unmarshal the configuration: %v", err)
	}
	return &config
}

func main() {
	var config *Config = loadConfig()
	var logLevel slog.Level
	if err := logLevel.UnmarshalText([]byte(config.Logging.DebugLevel)); err != nil {
		// Default behavior of the log level
		logLevel = slog.LevelDebug
		log.Printf("Invalid log level '%s' in config, defaulting to debug", config.Logging.DebugLevel)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	}))
	slog.SetDefault(logger)
	logger.Debug("Version 1 of Wallety!\n")
	logger.Debug("Add your expenses below \n DD-MM-YYYY\tTYPE\tCATEGORY\tMOP\tSOURCE\tDESCRIPTION\n")

	bteaProgram := tea.NewProgram(model{})
	if _, err := bteaProgram.Run(); err != nil {
		log.Fatalf("running term app failed. :%v\n", err)
	}
}
