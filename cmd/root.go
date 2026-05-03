package cmd

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/spf13/cobra"
)

var client *http.Client

var rootCmd = &cobra.Command{
	Use:   "sleuth [username]",
	Short: "Scan various websites for a username match",
	Long: `sleuth is a CLI tool for inspecting different websites for users with a given username.
	
It makes requests to each website listed in data/sites.json, and outputs what matches were found.
Supports multiple output formats, such as CSV, JSON, or plain-text.`,

	Args:             cobra.ArbitraryArgs,
	Run:              scan,
	PersistentPreRun: setupClient,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sleuth.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func setupClient(cmd *cobra.Command, args []string) {
	timeoutDuration := 10 * time.Second

	client = &http.Client{
		Timeout: timeoutDuration,
	}
}

func scan(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		return
	}

	for _, username := range args {
		url := fmt.Sprintf("https://github.com/%s", username)

		resp, err := client.Get(url)
		if err != nil {
			fmt.Printf("error making request")
			return
		}

		resp.Body.Close()

		// body, err := io.ReadAll(resp.Body)
		// if err != nil {
		// 	fmt.Printf("error reading bdy")
		// }

		if resp.StatusCode == 200 {
			lipgloss.Printf("%s (%s) - Account Found\n", successStyle.Render("[+] GitHub"), url)
		} else {
			lipgloss.Printf("%s (%s) - Account NOT Found\n", errorStyle.Render("[-] GitHub"), url)
		}
	}
}
