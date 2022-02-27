package main

import (
	bot "bott-the-pigeon/app/session"
	aws "bott-the-pigeon/lib/aws/session"
	ssm "bott-the-pigeon/lib/aws/ssm-env"
	"fmt"
	"log"

	"flag"
	"os"
	"os/signal"
	"syscall"
)

// Note: main should give a high-level overview of the E2E flow of the application.

func main() {

	config := *flagHandler()
	botTokenKey := getBotTokenKey(*config.prod)
	
	// This is the only place where logs can (should) be fatal, and terminate the app.
	err := setEnvs(getConfigs())
	if err != nil {
		log.Fatal(err)
	}

	awssess, err := aws.GetAWSSession()
	if err != nil {
		log.Fatal(err)
	}

	ssmEnv, err := ssm.GetEnv(awssess, os.Getenv("AWS_SSM_PARAMETER_PATH"))
	if err != nil {
		log.Fatal(err)
	}

	err = setEnvs(ssmEnv)
	if err != nil {
		log.Fatal(err)
	}

	bot, err := bot.GetBotSession(os.Getenv(botTokenKey))
	if err != nil {
		log.Fatal(err)
	}

	defer bot.Close()
	addCloseListener()

}

// Returns a k,v map of base configs. NON-SENSITIVE CONFIGS GO HERE.
func getConfigs() map[string]string {
	env := make(map[string]string)
	env["GITHUB_REPO_ACCOUNT"] = "BottThePigeon"
	env["GITHUB_PROJECT_ID"] = "1"
	env["GITHUB_SUGGESTIONS_COLUMN_ID"] = "17803319"
	env["AWS_SSM_PARAMETER_PATH"] = "/btp/"
	return env
}

// Initialises environment based upon provided k,v map.
func setEnvs(env map[string]string) error {
	errs := []error{}
	for k, v := range env {
		err := os.Setenv(k, v)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("failed env variable initialisation. error(s): %v", errs)
	}
	return nil
}

// Flag configurations for the application.
type flagConfig struct {
	prod *bool
}

// Parses the flag configurations for the application.
func flagHandler() *flagConfig {
	flags := &flagConfig{
		prod: flag.Bool("prod",
			false,
			"Should the production bot application be used?"),
	}
	flag.Parse()
	return flags
}

// Returns the token environment variable key based on isProd.
func getBotTokenKey(isProd bool) string {
	botTokenKey := "BOT_TOKEN_TEST"
	if isProd {
		botTokenKey = "BOT_TOKEN"
	}
	return botTokenKey
}

// Waits for a termination/kill etc. signal (Holding the application open).
func addCloseListener() {

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sigChan
}
