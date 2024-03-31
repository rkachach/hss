package main

import (
	"fmt"
	"strings"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"github.com/spf13/cobra"
	"github.com/c-bata/go-prompt"
	"github.com/spf13/pflag"
)

// Base URL for the REST API
const baseURL = "http://localhost:9000"

var currentDirectory = "/"

var rootCmd = &cobra.Command{
	Use:   "greet",
	Short: "A simple CLI tool to greet the user",
	Long:  `greet is a CLI tool that greets the user with a message.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello, World!")
	},
}

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all files and directories",
	Run: func(cmd *cobra.Command, args []string) {
		url := fmt.Sprintf("%s/%s/?type=%s&operation=%s", baseURL, currentDirectory, "directory", "list")
		showEntries(sendQuery(url), cmd.Flags().Changed("l"))
	},
}

var changeDirectory = &cobra.Command{
	Use:   "cd",
	Short: "Change current directory",
	Run: func(cmd *cobra.Command, args []string) {
		currentDirectory = args[0]
	},
}

var showCurrentDirectory = &cobra.Command{
	Use:   "pwd",
	Short: "Show current directory",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%v\n", currentDirectory)
	},
}

func showEntries(response []byte, longFormat bool){

	// Parse JSON response
	var files []struct {
		Name string `json:"name"`
	}
	err := json.Unmarshal(response, &files)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	if longFormat {
		for _, file := range files {
			fmt.Printf("- %v\n", file.Name)
		}
	} else {
		for _, file := range files {
			fmt.Printf("%v ", file.Name)
		}
		fmt.Println()
	}
}

func sendQuery(url string) []byte{

	// Send GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil
	}

	// Print response
	//fmt.Println("Response:", string(body))
	return body
}

func completer(d prompt.Document) []prompt.Suggest {

	// Set up suggestions slice
	var suggestions []prompt.Suggest

	// Populate suggestions from Cobra commands
	for _, cmd := range rootCmd.Commands() {
		suggestions = append(suggestions, prompt.Suggest{
			Text:        cmd.Name(),
			Description: cmd.Short,
		})


		// Add flags to suggestions
		flags := cmd.Flags()
		flags.VisitAll(func(flag *pflag.Flag) {
			suggestions = append(suggestions, prompt.Suggest{
				Text:        fmt.Sprintf("--%s", flag.Name),
				Description: flag.Usage,
			})
		})
	}


	// Add "exit" suggestion
	suggestions = append(suggestions, prompt.Suggest{
		Text:        "exit",
		Description: "Exit the application",
	})

	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}

func main() {

	rootCmd.AddCommand(listCmd)

	for {
		listCmd.ResetFlags()
		listCmd.Flags().BoolP("l", "l", false, "List in long format")

		input := prompt.Input("> ", completer)
		args := strings.Fields(input)
		if len(args) == 0 {
			continue
		}

		cmdName := args[0]
		cmdArgs := args[1:]

		switch cmdName {
		case "ls":
			listCmd.SetArgs(cmdArgs)
			listCmd.ParseFlags(cmdArgs)
			listCmd.Run(listCmd, cmdArgs)
		case "cd":
			changeDirectory.SetArgs(cmdArgs)
			changeDirectory.ParseFlags(cmdArgs)
			changeDirectory.Run(listCmd, cmdArgs)
		case "pwd":
			showCurrentDirectory.SetArgs(cmdArgs)
			showCurrentDirectory.ParseFlags(cmdArgs)
			showCurrentDirectory.Run(listCmd, cmdArgs)
		case "exit":
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			cmd, _, err := rootCmd.Find(args)
			if err == nil {
				cmd.SetArgs(cmdArgs)
				cmd.Execute()
			} else {
				fmt.Println("Unknown command:", cmdName)
			}
		}
	}
}
