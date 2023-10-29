package main

import (
	"os"

	"github.com/natemarks/iam-disable/disable"

	"github.com/natemarks/iam-disable/discover"
	"github.com/natemarks/iam-disable/file"
	"github.com/natemarks/iam-disable/user"

	"github.com/natemarks/iam-disable/version"
	"github.com/rs/zerolog"
)

func Discover(config user.Config, log *zerolog.Logger) {
	var targetContents string
	users := discover.Users(log)
	report := ""
	for _, user := range users {
		report += user.Report()
	}
	file.WriteToFile(config.ReportFile, report, true)
	for _, user := range users {
		targetContents += user.Username + "\n"
	}
	file.WriteToFile(config.TargetsFile, targetContents, false)
}

func Disable(config user.Config, log *zerolog.Logger) {
	targets, err := file.TargetsFromFile(config.TargetsFile)
	if err != nil {
		log.Fatal().Err(err).Msgf("error getting targets from file %s", config.TargetsFile)
	}
	for _, target := range targets {
		disable.IamUser(target, log)
	}
}

func main() {
	config := user.GetConfig()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger := zerolog.New(os.Stderr).With().Str("version", version.Version).Timestamp().Logger()
	logger = user.UpdateLogger(config.AWSAccount, &logger)
	logger = logger.With().Str("mode", config.Mode).Logger()
	logger.Info().Msg("starting")
	if config.Mode == "discover" {
		Discover(config, &logger)
		os.Exit(0)
	}
	if config.Mode == "disable" {
		Disable(config, &logger)
		os.Exit(0)
	}
}
