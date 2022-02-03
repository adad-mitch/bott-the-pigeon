package main

import (
	awsenv "bott-the-pigeon/aws-utils/aws-env"
	aws "bott-the-pigeon/aws-utils/init"
	bot "bott-the-pigeon/bot-utils/init"

	"flag"
	"os"
)

// Project root/start; primarily init function caller

// MAIN
func main() {

	initEnv()
	config := *flagHandler()

	//Works with the state of the --prod flag
	botTokenKey := getBotTokenKey(*config.prod);

	// Create AWS session and initialise environment stage
	awssess := aws.InitAws()
	awsenv.InitEnv(awssess)

	// Return a bot instance. This is merely an "artifact" to be closed, everything happens inside bot func
	bot := bot.InitBot(botTokenKey)

	// Close bot at EOF
	defer bot.Close()
}


// FLAGS/CONFIGS - (TODO) Maybe isolate in their own package?

// Miscellaneous (and non-confidential) environment variable initialisation (That doesn't need AWS) goes here
func initEnv() {
	os.Setenv("AWS_REGION", "eu-west-2") // AWS SDK Session Region
	os.Setenv("SSM_PARAMETER_PATH", "/btp/") // SSM Parameter Store location of project-related variables.
	//Note: The instance we're running on shouldn't have permission to use parameters outside this path. 
}

// Struct containing the flag configurations for the application
type flagConfig struct {
	prod *bool
}

//Parse and return the flag configurations for the application
func flagHandler() (*flagConfig) {

	flags := &flagConfig{
		prod: flag.Bool("prod", false, "Should the production bot application be used?"),
	}

	flag.Parse()
	return flags
}

// Determine what the key for the bot token needed is, based on if running in prod
// We use this rather than manually passing in strings because it gives a single and easily-traceable "source of truth"
// TODO: Perhaps there's a less error-prone way to do this, that doesn't require fighting against Go's lack of configs support?
func getBotTokenKey(isProd bool) (string) {

	botTokenKey := "BOT_TOKEN_TEST"
	if (isProd) {
		botTokenKey = "BOT_TOKEN"
	}

	return botTokenKey
}