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
	"github.com/matiasnu/go-jopit-toolkit/goutils/apierrors"
	"google.golang.org/api/option"
)

const (
	UserIDMock = "TEST-MOCK-USER"
)

var (
	fbClient *firebaseClient
	once     sync.Once
)

type firebaseCredential struct {
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

type firebaseClient struct {
	AuthClient *auth.Client
}

// init initiates the firebase client ONCE
func init() {
	once.Do(InitFirebase)
}

func InitFirebase() {

	opt := option.WithCredentialsFile("./config/credentials.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Println("Error connecting to firebase" + err.Error())
	}

	authClient, errAuth := app.Auth(context.Background())
	if errAuth != nil {
		log.Println("Error connecting to firebase" + errAuth.Error())
	}

	fbClient = &firebaseClient{
		AuthClient: authClient,
	}
}

func GetEmailFromUserID(c *gin.Context) (string, error) {

	userID, exist := c.Get("user_id")
	if !exist {
		return "", fmt.Errorf("expected to receive an user_id, but it was empty")
	}

	userRecord, err := fbClient.AuthClient.GetUser(c, userID.(string))
	if err != nil {
		return "", err
	}

	userEmail := userRecord.UserInfo.Email

	return userEmail, nil
}

func AuthWithFirebase() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		idToken := strings.TrimSpace(strings.Replace(header, "Bearer", "", 1))
		decodedToken, err := fbClient.AuthClient.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			c.AbortWithStatusJSON(401, err.Error())
			return
		}

		c.Set("user_id", decodedToken.UID)
		c.Next()
	}
}

func CheckFirebaseCredentials() error {
	var fields []string
	firebaseCredentials := firebaseCredential{}

	bytes, err := ioutil.ReadFile("./config/credentials.json")
	if err != nil {
		return fmt.Errorf("file not found")
	}

	err = json.Unmarshal(bytes, &firebaseCredentials)
	if err != nil {
		return fmt.Errorf("error unmarshalling the credentials")
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

func GetUserId(c *gin.Context) string {
	userID, exist := c.Get("user_id")
	if !exist {
		return ""
	}
	return userID.(string)
}

func MockAuthWithFirebase() gin.HandlerFunc {
	return func(c *gin.Context) {

		userID := c.GetHeader("Authorization")
		if userID == "" {
			userID = UserIDMock
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

type firebaseAccountManager struct {
	AuthClient *auth.Client
}

func NewFirebaseAccountManager() FirebaseAccountManager {
	return firebaseAccountManager{AuthClient: fbClient.AuthClient}
}

type FirebaseAccountManager interface {
	VerificationEmail(c *gin.Context, userEmail string) (string, apierrors.ApiError)
	ResetPassword(c *gin.Context, userEmail string) (string, apierrors.ApiError)
}

func (fam firebaseAccountManager) VerificationEmail(c *gin.Context, userEmail string) (string, apierrors.ApiError) {

	link, err := fam.AuthClient.EmailVerificationLink(c, userEmail)
	if err != nil {
		return "", apierrors.NewApiError("error on firebase verification . ", err.Error(), 500, apierrors.CauseList{})
	}

	return link, nil
}

func (fam firebaseAccountManager) ResetPassword(c *gin.Context, userEmail string) (string, apierrors.ApiError) {

	link, err := fam.AuthClient.PasswordResetLink(c, userEmail)
	if err != nil {
		return "", apierrors.NewApiError("error on firebase PasswordResetLink. ", err.Error(), 500, apierrors.CauseList{})
	}

	return link, nil
}
