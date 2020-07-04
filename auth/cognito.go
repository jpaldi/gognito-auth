package auth

import (
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type CognitoAuth struct {
	CognitoClient *cognito.CognitoIdentityProvider
	UserPoolID    string
	AppClientID   string
}
