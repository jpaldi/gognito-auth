package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/jpaldi/gognito-auth/auth"
)

type LoginHandler struct {
	Authenticator auth.CognitoAuth
}

type LoginSuccessfulResponse struct {
	Session             string
	ChallengeParameters ChallengeParameters
}

type ChallengeParameters struct {
	userAttributes     string
	USER_ID_FOR_SRP    string
	requiredAttributes string
}

func (l *LoginHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		l.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

func (l *LoginHandler) post(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	if len(username) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request: body should have a username"))
		return nil
	}

	if len(password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request: body should have a password"))
		return nil
	}

	flow := aws.String("ADMIN_NO_SRP_AUTH")

	authTry := &cognito.AdminInitiateAuthInput{
		AuthFlow: flow,
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(username),
			"PASSWORD": aws.String(password),
		},
		UserPoolId: aws.String(os.Getenv("AWS_USER_POOL_ID")),
		ClientId:   aws.String(os.Getenv("AWS_APP_CLIENT_ID")),
	}

	res, err := l.Authenticator.CognitoClient.AdminInitiateAuth(authTry)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
		return nil
	}

	response, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(response)

	return nil
}
