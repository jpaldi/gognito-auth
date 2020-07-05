package handlers

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/jpaldi/gognito-auth/auth"
)

type LoginHandler struct {
	Authenticator auth.CognitoAuth
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
	username := "set-me"
	password := "set-me"

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

	r.ParseForm()

	flow := aws.String("ADMIN_NO_SRP_AUTH")

	authTry := &cognito.AdminInitiateAuthInput{
		AuthFlow: flow,
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(username),
			"PASSWORD": aws.String(password),
		},
		UserPoolId: aws.String("set-me"),
		ClientId:   aws.String("set-me"),
	}

	res, err := l.Authenticator.CognitoClient.AdminInitiateAuth(authTry)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(err.Error()))
		return nil
	}

	fmt.Println(res)
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("login route"))

	return nil
}
