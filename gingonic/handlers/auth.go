package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/matiasnu/go-jopit-toolkit/goauth"
	"github.com/matiasnu/go-jopit-toolkit/goutils/apierrors"
	"github.com/matiasnu/go-jopit-toolkit/tracing"
	"google.golang.org/api/option"
)

func JopitAuth(scopes []string) gin.HandlerFunc {
	return JopitAuthWithOptions(scopes)
}

// JopitAuthWithOptions returns an authentication middleware with options.
//
// This function creates the JopitAuth middleware with non-standard authorization options.
// These options typically weaken the authorization.
// Use this function ONLY when you fully understand the consequences of of such options. Otherwise use JopitAuth.
func JopitAuthWithOptions(scopes []string, opts ...goauth.AuthOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := authenticateRequest(c.Request, scopes, opts...); err != nil {
			c.Abort()
			c.JSON(err.(apierrors.ApiError).Status(), err)
		}
	}
}

func authenticateRequest(r *http.Request, scopes []string, opts ...goauth.AuthOption) error {
	if err := goauth.AuthenticateRequestWithOptions(r, opts...); err != nil {
		return err
	}

	return nil
}

func HeaderForwarding() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := tracing.ContextFromRequest(c.Request)
		c.Request = c.Request.WithContext(ctx)
	}
}

// ========================================================= FIREBASE AUTH =========================================================
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

func (fc *FirebaseCredential) loadCredentials() error {

	fc.Type = os.Getenv("TYPE")
	fc.ProjectId = os.Getenv("PROYECT_ID")
	fc.PrivateKeyId = os.Getenv("PRIVATE_KEY_ID")
	fc.PrivateKey = os.Getenv("PRIVATE_KEY")
	fc.ClientEmail = os.Getenv("CLIENT_EMAIL")
	fc.ClientId = os.Getenv("CLIENT_ID")
	fc.AuthUri = os.Getenv("AUTH_URI")
	fc.TokenUri = os.Getenv("TOKEN_URI")
	fc.AuthProviderX509CertUrl = os.Getenv("AUTH_PROVIDER_x509_CERT_URL")
	fc.ClientX509CertUrl = os.Getenv("CLIENT_x509_CERT_URL")

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

	firebaseCredential := FirebaseCredential{}
	firebaseCredential.loadCredentials()

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
		_, err := firebaseClient.AuthClient.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			apiErr := apierrors.NewInternalServerApiError("error getting token", err)
			c.AbortWithError(401, apiErr)
			return
		}

		c.Next()
	}
}
