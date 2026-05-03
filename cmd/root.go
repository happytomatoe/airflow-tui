package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/airflow-tui/airflow-tui/config"
	"github.com/airflow-tui/airflow-tui/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "airflow-tui",
	Short: "TUI for Apache Airflow",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTUI()
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the TUI",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTUI()
	},
}

func runTUI() error {
	cfgPath := config.ConfigPath()
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(cfgPath), 0755)
		defaultConfig := `servers: []
active_server: ""
poll_interval_ms: 2000
`
		os.WriteFile(cfgPath, []byte(defaultConfig), 0644)
		fmt.Printf("Created config at %s\n", cfgPath)
	}

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(cfgPath)
	if err := v.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
		os.Exit(1)
	}

	var cfg config.Config
	if err := v.Unmarshal(&cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing config: %v\n", err)
		os.Exit(1)
	}

	model := ui.NewModel(cfg)
	p := tea.NewProgram(model)
	p.Run()
	return nil
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func Execute() {
	rootCmd.Execute()
}
