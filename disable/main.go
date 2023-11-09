package disable

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/natemarks/iam-disable/discover"
	"github.com/rs/zerolog"
)

const (
	tagKey = "iam-disabled-date"
)

func disableIAMUserConsolePassword(username string) error {
	// Load your AWS credentials and create a new IAM client.
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := iam.NewFromConfig(cfg)

	// Create input for deleting the login profile (console password) for the user.
	input := &iam.DeleteLoginProfileInput{
		UserName: aws.String(username),
	}

	// Delete the user's login profile (console password).
	_, err = client.DeleteLoginProfile(context.TODO(), input)
	if err != nil {
		return err
	}
	return nil
}

func disableAccessKey(accessKeyID string, userName string) error {
	// Load your AWS credentials and configuration from a shared credentials file
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	// Create an IAM client
	client := iam.NewFromConfig(cfg)

	// Create an input struct for UpdateAccessKey operation
	input := &iam.UpdateAccessKeyInput{
		AccessKeyId: &accessKeyID,
		Status:      types.StatusTypeInactive,
		UserName:    &userName,
	}

	// Call the UpdateAccessKey operation to disable the access key
	_, err = client.UpdateAccessKey(context.TODO(), input)
	return err
}

func Today() string {
	// Get the current time in the local timezone
	currentTime := time.Now()

	// Format the current time as a string in the "2006-01-02" layout
	dateString := currentTime.Format("2006-01-02")

	return dateString
}
func TagIamUser(username string) (err error) {
	// Load your AWS credentials and configuration from a shared credentials file
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	client := iam.NewFromConfig(cfg)

	// List the tags for the IAM user
	listTagsOutput, err := client.ListUserTags(ctx, &iam.ListUserTagsInput{
		UserName: aws.String(username),
	})
	if err != nil {
		return err
	}

	// Check if the tag key already exists
	tagExists := false
	for _, tag := range listTagsOutput.Tags {
		if *tag.Key == tagKey {
			tagExists = true
			break
		}
	}

	// If the tag key doesn't exist, create the tag
	if !tagExists {
		_, err := client.TagUser(ctx, &iam.TagUserInput{
			UserName: aws.String(username),
			Tags: []types.Tag{
				{
					Key:   aws.String(tagKey),
					Value: aws.String(Today()),
				},
			},
		})
		if err != nil {
			return err
		}
	}

	return err
}

// IamUser disables an IAM user
func IamUser(username string, log *zerolog.Logger) {
	userLog := log.With().Str("username", username).Logger()
	err := TagIamUser(username)
	if err != nil {
		userLog.Error().Msgf("failed to tag user %s, %s", username, err)
	}
	userLog.Info().Msgf("disabling user %s", username)
	err = disableIAMUserConsolePassword(username)
	if err != nil {
		userLog.Error().Msgf("failed to disable console password for user %s, %s", username, err)
	} else {
		userLog.Info().Msgf("disabled console password for user %s", username)
	}
	accessKeys := discover.GetAccessKeys(username, &userLog)

	for _, key := range accessKeys {
		err := disableAccessKey(key.AccessKeyID, username)
		if err != nil {
			userLog.Error().Msgf("failed to disable access key %s for user %s, %s", key.AccessKeyID, username, err)
		} else {
			userLog.Info().Msgf("disabled access key %s for user %s", key.AccessKeyID, username)
		}
	}

}
