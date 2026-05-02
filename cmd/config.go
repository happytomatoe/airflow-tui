package cmd

import (
	"fmt"
	"os"

	"github.com/airflow-tui/airflow-tui/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

var configAddCmd = &cobra.Command{
	Use:   "add [name] [url]",
	Short: "Add a server",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		url := args[1]

		authType, _ := cmd.Flags().GetString("auth")
		username, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("pass")
		token, _ := cmd.Flags().GetString("token")
		apiVersion, _ := cmd.Flags().GetString("api")

		cfgPath := config.ConfigPath()
		v := config.LoadConfig(cfgPath)

		server := config.ServerConfig{
			Name:       name,
			URL:        url,
			APIVersion: apiVersion,
			Auth: config.AuthConfig{
				Type:     authType,
				Username: username,
				Password: password,
				Token:    token,
			},
		}

		v.Servers = append(v.Servers, server)
		if v.ActiveServer == "" {
			v.ActiveServer = name
		}

		config.SaveConfig(cfgPath, v)
		fmt.Printf("Added server %s\n", name)
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List servers",
	Run: func(cmd *cobra.Command, args []string) {
		cfgPath := config.ConfigPath()
		v := config.LoadConfig(cfgPath)

		if len(v.Servers) == 0 {
			fmt.Println("No servers configured")
			return
		}

		for _, s := range v.Servers {
			active := ""
			if s.Name == v.ActiveServer {
				active = " (active)"
			}
			fmt.Printf("  %s%s\n    URL: %s\n    API: %s\n", s.Name, active, s.URL, s.APIVersion)
		}
	},
}

var configRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		cfgPath := config.ConfigPath()
		v := config.LoadConfig(cfgPath)

		found := false
		newServers := []config.ServerConfig{}
		for _, s := range v.Servers {
			if s.Name == name {
				found = true
				continue
			}
			newServers = append(newServers, s)
		}

		if !found {
			fmt.Printf("Server %s not found\n", name)
			os.Exit(1)
		}

		v.Servers = newServers
		if v.ActiveServer == name && len(newServers) > 0 {
			v.ActiveServer = newServers[0].Name
		} else if len(newServers) == 0 {
			v.ActiveServer = ""
		}

		config.SaveConfig(cfgPath, v)
		fmt.Printf("Removed server %s\n", name)
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set-active [name]",
	Short: "Set active server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		cfgPath := config.ConfigPath()
		v := config.LoadConfig(cfgPath)

		found := false
		for _, s := range v.Servers {
			if s.Name == name {
				found = true
				break
			}
		}

		if !found {
			fmt.Printf("Server %s not found\n", name)
			os.Exit(1)
		}

		v.ActiveServer = name
		config.SaveConfig(cfgPath, v)
		fmt.Printf("Set active server to %s\n", name)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configAddCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configRemoveCmd)
	configCmd.AddCommand(configSetCmd)

	configAddCmd.Flags().StringP("auth", "a", "basic", "Auth type (basic, token)")
	configAddCmd.Flags().StringP("user", "u", "", "Username for basic auth")
	configAddCmd.Flags().StringP("pass", "p", "", "Password for basic auth")
	configAddCmd.Flags().StringP("token", "t", "", "Token for token auth")
	configAddCmd.Flags().StringP("api", "i", "v1", "API version (v1, v2)")
}
