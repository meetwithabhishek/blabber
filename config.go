package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Username      string `yaml:"username"`
	ServerAddress string `yaml:"serverAddress"`
}

var conf Config

const ClientConfigFilename = "client.conf"

var promptStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4"))

	// Define a style for the input echoed back
var _ = lipgloss.NewStyle().
	Italic(true).
	Foreground(lipgloss.Color("#FF6188")) // pink/red

var infoStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#A9B7C6")).
	Italic(true)

func ensureConfigExists() error {
	// create a sample config file, if it doesn't exists
	_, err := os.Stat(GetPlayPath())
	if err != nil {
		err := os.MkdirAll(GetPlayPath(), 0755)
		if err != nil {
			panic(err)
		}
	}

	// Try to open the file (only for checking existence)
	if _, err := os.Stat(GetPlayPath(ClientConfigFilename)); os.IsNotExist(err) {
		// File does not exist, create it
		file, createErr := os.Create(GetPlayPath(ClientConfigFilename))
		if createErr != nil {
			return createErr
		}

		fmt.Println(infoStyle.Render("First time user, creating config, please provide the following details:"))
		fmt.Println()

		// Print styled prompt
		fmt.Print(promptStyle.Render("Enter Username") + ": ")

		// Read user input (one line)
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return err
		}

		username := strings.TrimSpace(input) // clean newline

		// Print styled prompt
		fmt.Print(promptStyle.Render("Enter Server Address") + ": ")

		// Read user input (one line)
		reader = bufio.NewReader(os.Stdin)
		input, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return err
		}

		serverAddress := strings.TrimSpace(input) // clean newline

		d, err := yaml.Marshal(Config{
			Username:      username,
			ServerAddress: serverAddress,
		})

		if err != nil {
			return err
		}

		if _, err := file.Write(d); err != nil {
			return err
		}

		file.Close()
	} else if err != nil {
		return err
	}

	// Load Conf

	path := GetPlayPath(ClientConfigFilename)

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return err
	}

	return nil
}
