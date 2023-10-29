package discover

import (
	"context"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/rs/zerolog"
)

type AccessKey struct {
	AccessKeyID string    `json:"accessKeyId"`
	Enabled     bool      `json:"enabled"`
	LastUsed    time.Time `json:"lastUsed"`
}
type IAMUser struct {
	Username          string      `json:"username"`
	HasConsoleProfile bool        `json:"hasConsoleProfile"`
	AccessKeys        []AccessKey `json:"accessKeys"`
	IAMGroups         []string    `json:"iamGroups"`
}

func (user IAMUser) Report() (report string) {
	report += "Username: " + user.Username + "\n"
	report += "HasConsoleProfile: " + strconv.FormatBool(user.HasConsoleProfile) + "\n"
	report += "Access Keys:\n"
	for _, key := range user.AccessKeys {
		report += "    AccessKeyID: " + key.AccessKeyID + "\n"
		report += "    Enabled: " + strconv.FormatBool(key.Enabled) + "\n"
		report += "    LastUsed: " + key.LastUsed.String() + "\n\n"
	}
	report += "IAM Groups:\n"
	for _, group := range user.IAMGroups {
		report += "    " + group + "\n"
	}
	report += "\n"
	return report
}

// GetAccessKeyLastUsedTime returns the last used time for an access key
// this function ignores response errors and returns a zero time.Time and an error
// the calling function should check for a zero time.Time and handle/log the error
// If the access key has never been used, the time will be time.Time{}
func GetAccessKeyLastUsedTime(accessKeyID string) (lastUsed time.Time, err error) {
	// Load your AWS credentials and configuration from a shared credentials file
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	// Create an IAM client
	client := iam.NewFromConfig(cfg)

	// Create an input struct for GetAccessKeyLastUsed operation
	input := &iam.GetAccessKeyLastUsedInput{
		AccessKeyId: aws.String(accessKeyID),
	}

	// Call the GetAccessKeyLastUsed operation
	resp, err := client.GetAccessKeyLastUsed(context.TODO(), input)
	//
	if err != nil {
		return time.Time{}, err
	}

	// Check if the access key has been used before
	if resp.AccessKeyLastUsed.LastUsedDate != nil {
		return *resp.AccessKeyLastUsed.LastUsedDate, err
	}
	return time.Time{}, err
}

// GetIAMGroupsForUser returns a list of IAM groups for a user
func GetIAMGroupsForUser(userName string) (groups []string) {
	// Load your AWS credentials and configuration from a shared credentials file
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	// Create an IAM client
	client := iam.NewFromConfig(cfg)

	// Create an input struct for ListGroupsForUser operation
	input := &iam.ListGroupsForUserInput{
		UserName: &userName,
	}

	// Call the ListGroupsForUser operation to get the IAM groups for the user
	resp, err := client.ListGroupsForUser(context.TODO(), input)
	if err != nil {
		return groups
	}

	for _, group := range resp.Groups {
		groups = append(groups, *group.GroupName)
	}

	return groups
}

// GetAccessKeys returns a list of access keys for a user
// errors are logged, but not returned
func GetAccessKeys(username string, log *zerolog.Logger) (accessKeys []AccessKey) {
	// Load your AWS credentials and create a new IAM client.
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := iam.NewFromConfig(cfg)

	// Create input for listing access keys for the user.
	input := &iam.ListAccessKeysInput{
		UserName: aws.String(username),
	}

	// List access keys for the user.
	result, err := client.ListAccessKeys(context.TODO(), input)
	if err != nil {
		log.Error().Err(err).Msg("error listing access keys")
		return accessKeys
	}

	for _, key := range result.AccessKeyMetadata {
		thisKey := AccessKey{
			AccessKeyID: *key.AccessKeyId,
			Enabled:     key.Status == types.StatusType("Active"),
		}
		lastUsed, err := GetAccessKeyLastUsedTime(*key.AccessKeyId)
		if err != nil {
			log.Error().Err(err).Msgf("error getting last used time for %s : %s", username, *key.AccessKeyId)
		}
		thisKey.LastUsed = lastUsed
		accessKeys = append(accessKeys, thisKey)
	}

	return accessKeys
}

// UserHasConsoleProfile returns true if the user has a console profile
// if there is an error in the response, the function returns false and the error
// the calling function should check for false and handle/log the error
func UserHasConsoleProfile(username string) (hasProfile bool, err error) {
	// Load your AWS credentials and create a new IAM client.
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := iam.NewFromConfig(cfg)

	// Create input for getting the login profile for the user.
	input := &iam.GetLoginProfileInput{
		UserName: aws.String(username),
	}

	// Get the login profile for the user.
	resp, err := client.GetLoginProfile(context.TODO(), input)
	if err != nil {
		return false, err
	}

	// Check if the login profile is enabled.
	if resp == nil {
		return false, err
	}

	return true, err
}

// Users returns a list of IAM users
func Users(log *zerolog.Logger) (users []IAMUser) {

	// Load your AWS credentials and configuration from a shared credentials file
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	// Create an IAM client
	client := iam.NewFromConfig(cfg)

	// Create a paginator to list all IAM users
	paginator := iam.NewListUsersPaginator(client, &iam.ListUsersInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(context.TODO())
		if err != nil {
			panic("failed to list users, " + err.Error())
		}

		// Process the list of IAM users on this page
		for _, user := range resp.Users {
			thisUser := IAMUser{
				Username: *user.UserName,
			}
			hasProfile, err := UserHasConsoleProfile(*user.UserName)
			if err != nil {
				log.Error().Err(err).Msgf("error checking for console profile: %s", *user.UserName)
			}
			thisUser.HasConsoleProfile = hasProfile
			accessKeys := GetAccessKeys(*user.UserName, log)
			thisUser.AccessKeys = accessKeys
			thisUser.IAMGroups = GetIAMGroupsForUser(*user.UserName)
			users = append(users, thisUser)
		}
	}

	return users
}
