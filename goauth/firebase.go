package goauth

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
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

func LoadFirebaseCredentialFile(Type, ProjectId, PrivateKeyId, PrivateKey, ClientEmail, ClientId, AuthUri, TokenUri, AuthProviderX509CertUrl, ClientX509CertUrl string) error {

	fc := FirebaseCredential{}

	fc.Type = Type
	fc.ProjectId = ProjectId
	fc.PrivateKeyId = PrivateKeyId
	fc.PrivateKey = PrivateKey
	fc.ClientEmail = ClientEmail
	fc.ClientId = ClientId
	fc.AuthUri = AuthUri
	fc.TokenUri = TokenUri
	fc.AuthProviderX509CertUrl = AuthProviderX509CertUrl
	fc.ClientX509CertUrl = ClientX509CertUrl

	bytes, err := json.Marshal(fc)
	if err != nil {
		return err
	}

	_ = ioutil.WriteFile("credentials.json", bytes, 0644)

	return nil
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
		a, err := firebaseClient.AuthClient.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			apiErr := apierrors.NewInternalServerApiError("error getting token", err)
			c.AbortWithError(401, apiErr)
			return
		}

		c.Set("user_id", a.UID)
		c.Next()
	}
}
