package goauth

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/matiasnu/go-jopit-toolkit/goutils/apierrors"
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

	opt := option.WithCredentialsFile("credentials.json")
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
			apiErr := apierrors.NewInternalServerApiError("error getting token", err)
			c.AbortWithError(401, apiErr)
			return
		}

		c.Set("user_id", decodedToken.UID)
		c.Next()
	}
}

func LoadFirebaseCredentials() error {

	firebaseCredentials := FirebaseCredential{
		Type:                    os.Getenv("FB_TYPE"),
		ProjectId:               os.Getenv("FB_PROJECT_ID"),
		PrivateKeyId:            os.Getenv("FB_PRIVATE_KEY_ID"),
		PrivateKey:              os.Getenv("FB_PRIVATE_KEY"),
		ClientEmail:             os.Getenv("FB_CLIENT_EMAIL"),
		ClientId:                os.Getenv("FB_CLIENT_ID"),
		AuthUri:                 os.Getenv("FB_AUTH_URI"),
		TokenUri:                os.Getenv("FB_TOKEN_URI"),
		AuthProviderX509CertUrl: os.Getenv("FB_AUTH_PROVIDER_X509_CERT_URL"),
		ClientX509CertUrl:       os.Getenv("FB_CLIENT_X509_CERT_URL"),
	}

	if firebaseCredentials.Type == "" || firebaseCredentials.ProjectId == "" || firebaseCredentials.PrivateKeyId == "" || firebaseCredentials.PrivateKey == "" || firebaseCredentials.ClientEmail == "" || firebaseCredentials.ClientId == "" || firebaseCredentials.TokenUri == "" || firebaseCredentials.AuthProviderX509CertUrl == "" || firebaseCredentials.ClientX509CertUrl == "" {
		return fmt.Errorf("error loading credentials. Some of them are nil")
	}

	return nil
}
