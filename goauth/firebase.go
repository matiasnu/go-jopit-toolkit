package goauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

var (
	firebaseClient *FirebaseClient
	once           sync.Once
)

type FirebaseCredential struct {
	Type                    string `json:"type"`
	ProjectId               string `json:"project_id"`
	PrivateKeyId            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientId                string `json:"client_id"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
}

type FirebaseClient struct {
	AuthClient *auth.Client
}

//initiates the firebase client ONCE
func NewfirebaseService() *FirebaseClient {
	once.Do(InitFirebase)

	return firebaseClient
}

func InitFirebase() {

	opt := option.WithCredentialsFile("./config/credentials.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Println("Error connecting to firebase" + err.Error())
	}

	auth, err2 := app.Auth(context.Background())
	if err2 != nil {
		log.Println("Error connecting to firebase" + err2.Error())
	}

	firebaseClient = &FirebaseClient{
		AuthClient: auth,
	}
}

func AuthWithFirebase() gin.HandlerFunc {
	return func(c *gin.Context) {

		header := c.GetHeader("HeaderAuthorization")
		idToken := strings.TrimSpace(strings.Replace(header, "Bearer", "", 1))
		decodedToken, err := firebaseClient.AuthClient.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			c.JSON(401, err)
			return
		}

		c.Set("user_id", decodedToken.UID)
		c.Next()
	}
}

func CheckFirebaseCredentials() error {

	var fields []string
	firebaseCredentials := FirebaseCredential{}

	bytes, err := ioutil.ReadFile("./config/credentials.json")
	if err != nil {
		return fmt.Errorf("file not found")
	}

	err = json.Unmarshal(bytes, &firebaseCredentials)
	if err != nil {
		return fmt.Errorf("error unmarshalling the credentials.json")
	}

	if firebaseCredentials.Type == "" {
		fields = append(fields, "type is nil")
	}
	if firebaseCredentials.ProjectId == "" {
		fields = append(fields, "projectId is nil")
	}
	if firebaseCredentials.PrivateKeyId == "" {
		fields = append(fields, "privateKeyId is nil")
	}
	if firebaseCredentials.PrivateKey == "" {
		fields = append(fields, "privateKey is nil")
	}
	if firebaseCredentials.ClientEmail == "" {
		fields = append(fields, "clientEmail is nil")
	}
	if firebaseCredentials.ClientId == "" {
		fields = append(fields, "clientId is nil")
	}
	if firebaseCredentials.AuthUri == "" {
		fields = append(fields, "authUri is nil")
	}
	if firebaseCredentials.TokenUri == "" {
		fields = append(fields, "tokenUri is nil")
	}
	if firebaseCredentials.AuthProviderX509CertUrl == "" {
		fields = append(fields, "authProviderX509CertUrl is nil")
	}
	if firebaseCredentials.ClientX509CertUrl == "" {
		fields = append(fields, "clientX509CertUrl is nil")
	}

	if len(fields) != 0 {
		return fmt.Errorf("some credentials values are nil: %s", fields)
	}

	return nil
}
