package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gofrs/flock"
	"github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
)

type Config struct {
	VaultAddress    string `mapstructure:"vault_address"`
	VaultToken      string `mapstructure:"vault_token"`
	VaultSecretPath string `mapstructure:"vault_secret_path"`
	OutputFile      string `mapstructure:"output_file"`
}

type Vault struct {
	Client *api.Client
}

func (v *Vault) GetSecret(secretPath string) (map[string]interface{}, error) {
	// Get the latest secret version
	secret, err := v.Client.Logical().Read(secretPath)
	if err != nil {
		return nil, err
	}
	secretData := secret.Data["data"].(map[string]interface{})
	return secretData, nil
}

func main() {
	// Load configuration
	// Define the command-line flag for the config file path
	configFile := flag.String("config", "", "path to the config file")
	flag.Parse()

	// Read the config file path from the command-line argument
	if *configFile == "" {
		log.Fatal("config file path is required")
	}

	// Set the config file path for Viper
	viper.SetConfigFile(*configFile)

	// Load the config file
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("error reading configuration file: %s", err)
	}

	// Initialize the config struct
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("error unmarshaling configuration file: %s", err)
	}

	// Initialize a lock file to ensure only one instance of the program is running at a time
	lock := flock.New("./vault-monitor.lock")

	// Try to obtain the lock
	locked, err := lock.TryLock()
	if err != nil {
		log.Fatalf("error obtaining lock: %s", err)
	}
	if !locked {
		log.Fatalf("unable to obtain lock, another instance of the program is running")
	}
	defer lock.Unlock()

	// Create a HashiCorp Vault client
	clientConfig := api.DefaultConfig()
	clientConfig.Address = config.VaultAddress
	client, err := api.NewClient(clientConfig)
	if err != nil {
		log.Fatalf("unable to initialize Vault client: %v", err)
	}
	client.SetToken(config.VaultToken)
	vault := &Vault{Client: client}

	// Initialize the previous secret version and data
	var previousVersion string
	var previousSecretData map[string]interface{}

	// Create a context and cancel function to allow for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Handle SIGINT and SIGTERM signals to gracefully shutdown the program
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			// If the context is done, exit the program
			log.Println("gracefully shutting down")
			return
		default:
			// Check for secret updates
			log.Println("Checking for secret updates...")
			secretData, err := vault.GetSecret(config.VaultSecretPath)
			if err != nil {
				log.Printf("unable to get secret from Vault: %v", err)
				time.Sleep(10 * time.Second)
				continue
			}
			latestVersion := fmt.Sprintf("%v", secretData["version"])

			// If the latest version is different from the previous version, update the secret data and write it to file
			if latestVersion != previousVersion || !isSecretDataEqual(secretData, previousSecretData) {
				log.Printf("New secret version detected: %s", latestVersion)
				for key, value := range secretData {
					os.Setenv(key, fmt.Sprintf("%v", value))
				}
				var outputData []string
				for key, value := range secretData {
					outputData = append(outputData, fmt.Sprintf("%s=%v", key, value))
				}
				outputString := strings.Join(outputData, "\n")
				err = ioutil.WriteFile(config.OutputFile, []byte(outputString), 0644)
				if err != nil {
					log.Printf("unable to write output file: %v", err)
				} else {
					log.Printf("Successfully updated secret file: %s", config.OutputFile)
				}

				// Update the previous version and secret data to the latest values
				previousVersion = latestVersion
				previousSecretData = secretData
			} else {
				log.Println("No secret updates detected")
			}

			// Sleep for 10 seconds before checking again
			time.Sleep(10 * time.Second)
		}
	}
}

func isSecretDataEqual(data1, data2 map[string]interface{}) bool {
	// Check if the number of keys is the same
	if len(data1) != len(data2) {
		return false
	}

	// Check if the values of all keys are equal
	for key, value1 := range data1 {
		value2, ok := data2[key]
		if !ok || value1 != value2 {
			return false
		}
	}

	return true
}
