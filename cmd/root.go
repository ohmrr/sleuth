package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type Site struct {
	Name    string `json:"name"`
	BaseURL string `json:"baseUrl"`
	URL     string `json:"url"`
}

var client *http.Client

var rootCmd = &cobra.Command{
	Use:   "sleuth [username]",
	Short: "Scan various websites for a username match",
	Long: `sleuth is a CLI tool for inspecting different websites for users with a given username.
	
It makes requests to each website listed in data/sites.json, and outputs what matches were found.
Supports multiple output formats, such as CSV, JSON, or plain-text.`,

	Args:             cobra.ArbitraryArgs,
	Run:              sleuth,
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

func loadSiteData() []Site {
	jsonFile, err := os.Open("data/sites.json")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	var sites []Site

	err = json.Unmarshal(byteValue, &sites)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil
	}

	return sites
}

func sleuth(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		return
	}

	sites := loadSiteData()
	for _, username := range args {
		fmt.Printf("\nBeginning Scan for %s\n\n", username)
		found, not_found := scan(username, sites)

		fmt.Printf("Found: %d | Not Found: %d\n", found, not_found)
	}
}

func scan(username string, sites []Site) (int, int) {
	found := 0

	for _, site := range sites {
		url := strings.Replace(site.URL, "{}", username, 1)

		resp, err := client.Get(url)
		if err != nil {
			fmt.Printf("Error:", err)
			os.Exit(1)
		}

		if resp.StatusCode == 200 {
			fmt.Printf("[+] %s (%s) - Account Found\n", site.Name, url)
			found++
		} else {
			fmt.Printf("[-] %s (%s) - Account NOT Found [%d]\n", site.Name, url, resp.StatusCode)
		}
	}

	return found, len(sites) - found
}
