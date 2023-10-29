package user

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/rs/zerolog"
)

// confirmString is the string that the user must type to confirm the action
const characterSet = "abcdefghijklmnopqrstuvwxyz"
const confirmStringLength = 4

type AWSAccount struct {
	Account string
	Arn     string
	User    string
}

func (account AWSAccount) Report() (report string) {
	report += "Discovering IAM users using the following credentials:\n"
	report += "Account: " + account.Account + "\n"
	report += "Arn: " + account.Arn + "\n"
	report += "User: " + account.User + "\n"
	return report
}

// GetAWSAccount returns the current AWS account
func GetAWSAccount() AWSAccount {
	// Load your AWS credentials and configuration from a shared credentials file
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	// Create an STS client
	client := sts.NewFromConfig(cfg)

	// Create an input for GetCallerIdentity operation
	input := &sts.GetCallerIdentityInput{}

	// Call GetCallerIdentity to get the current identity
	resp, err := client.GetCallerIdentity(context.TODO(), input)
	return AWSAccount{
		Account: *resp.Account,
		Arn:     *resp.Arn,
		User:    *resp.UserId,
	}
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = characterSet[rand.Intn(len(characterSet))]
	}
	return string(b)
}

func ConfirmDisable(targets []string) {
	warning := "WARNING: This will disable the following IAM users:\n"
	for _, target := range targets {
		warning += "    " + target + "\n"
	}
	confirmString := generateRandomString(confirmStringLength)
	fmt.Println("\nType the characters: \"" + confirmString + "\" to continue:\n")
	reader := bufio.NewReader(os.Stdin)
	userInput, _ := reader.ReadString('\n')
	userInput = strings.TrimSpace(userInput) // Remove trailing newline character
	if strings.ToLower(userInput) != confirmString {
		panic("User input did not match confirmation string")
	}
}

func UpdateLogger(account AWSAccount, log *zerolog.Logger) zerolog.Logger {
	return log.With().Str("account", account.Account).Str("arn", account.Arn).Str("user", account.User).Logger()
}

type Config struct {
	Mode              string
	TargetsFile       string
	ReportFile        string
	TargetsFileExists bool
	AWSAccount        AWSAccount
}

func Usage() (usage string) {
	usage += "Usage:\n"
	usage += "  iam-disable [target file path]\n"
	usage += "\n"
	usage += "If run with no arguments, discover IAM user accounts and create a targets file if one doesn't\n"
	usage += "exist. example: 0123456789_targets.txt\n\n"
	usage += "create/overwrite a report file. example: 0123456789_report.txt\n"
	usage += "\n"
	usage += "If run with one argument (target file path), disable the IAM users in the targets file.\n"
	return usage
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// GetConfig returns the configuration for the program
func GetConfig() Config {
	// start building the config without user input
	account := GetAWSAccount()
	targetFile := account.Account + "_targets.txt"
	config := Config{
		TargetsFile:       targetFile,
		ReportFile:        account.Account + "_report.txt",
		TargetsFileExists: fileExists(targetFile),
		AWSAccount:        account,
	}
	// update the config with user input
	args := os.Args
	if len(args) > 2 {
		panic("Too many arguments. expected 0 or 1\n " + Usage())
	}
	// check for help flag
	if len(args) == 2 {
		if args[1] == "-h" || args[1] == "--help" || strings.ToLower(args[1]) == "help" {
			fmt.Println(Usage())
			os.Exit(0)
		}
	}
	if len(args) == 1 {
		config.Mode = "discover"
		return config
	}
	config.Mode = "disable"
	config.TargetsFile = args[1]
	return config
}
