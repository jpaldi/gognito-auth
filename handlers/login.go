package handlers

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/jpaldi/gognito-auth/auth"
)

const flowUsernamePassword = "USER_PASSWORD_AUTH"
const flowRefreshToken = "REFRESH_TOKEN_AUTH"

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
	username := r.Form.Get("username")
	password := r.Form.Get("password")

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

	refresh := r.Form.Get("refresh")
	refreshToken := r.Form.Get("refresh_token")

	flow := aws.String(flowUsernamePassword)
	params := map[string]*string{
		"USERNAME": aws.String(username),
		"PASSWORD": aws.String(password),
	}

	if refresh != "" {
		flow = aws.String(flowRefreshToken)
		params = map[string]*string{
			"REFRESH_TOKEN": aws.String(refreshToken),
		}
	}

	authTry := &cognito.InitiateAuthInput{
		AuthFlow:       flow,
		AuthParameters: params,
		ClientId:       aws.String(l.Authenticator.AppClientID),
	}

	res, err := l.Authenticator.CognitoClient.InitiateAuth(authTry)
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
