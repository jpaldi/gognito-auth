package main

import (
	"fmt"
	"net/http"
	"os"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/joho/godotenv"
	"github.com/jpaldi/gognito-auth/auth"
	"github.com/jpaldi/gognito-auth/handlers"
)

func main() {
	godotenv.Load("secrets.env")

	conf := &aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))}
	sess, err := session.NewSession(conf)
	if err != nil {
		panic(fmt.Errorf("connection with aws failed, %w", err))
	}

	auth := auth.CognitoAuth{
		CognitoClient: cognito.New(sess),
	}

	loginHandler := handlers.LoginHandler{
		Authenticator: auth,
	}

	fmt.Println("api running")
	http.HandleFunc("/login", loginHandler.Handle)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
