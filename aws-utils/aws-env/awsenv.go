package awsenv

import (
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// Functions that initialise the OS environment via AWS SSM.

// Load Environment Variables
func InitEnv(awssess *session.Session) {
	
	// Composite of retrieving SSM parameters and assigning to OS Env.
	ssmEnv := loadEnvFromSSM(awssess)
	loadSSMEnvIntoOS(ssmEnv)
}

// Load environment variables from AWS SSM Parameter Store path Bott-The-Pigeon
func loadEnvFromSSM(awssess *session.Session) (*ssm.GetParametersByPathOutput) {

	ssmsvc := ssm.New(awssess, aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")))
	paramPath, decrypt := os.Getenv("SSM_PARAMETER_PATH"), true
	ssmparams, err := ssmsvc.GetParametersByPath( &ssm.GetParametersByPathInput {
		Path: &paramPath,
		WithDecryption: &decrypt,
	})
	if err != nil {
		log.Fatal("Could not obtain application credentials from AWS: ", err)
	}

	return ssmparams
}

// Set OS Environment Variables based on AWS SSM Output.
func loadSSMEnvIntoOS(ssmparams *ssm.GetParametersByPathOutput) {

	// Iterate through SSM Parameters passed, assign env var of associated name and values.
	for i := 0; i < len(ssmparams.Parameters); i++ {
		k := strings.ReplaceAll(*ssmparams.Parameters[i].Name, os.Getenv("SSM_PARAMETER_PATH"), "") // Remove the path name
		v := *ssmparams.Parameters[i].Value
		os.Setenv(k, v)
	} 
}