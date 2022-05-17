package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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
	ProyectId               string `json:"proyect_id"`
	PrivateKeyId            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientId                string `json:"client_id"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
}

func LoadFirebaseCredentialFile(Type, ProyectId, PrivateKeyId, PrivateKey, ClientEmail, ClientId, AuthUri, TokenUri, AuthProviderX509CertUrl, ClientX509CertUrl string) error {

	fc := FirebaseCredential{}

	fc.Type = Type
	fc.ProyectId = ProyectId
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
		_, err := firebaseClient.AuthClient.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			apiErr := apierrors.NewInternalServerApiError("error getting token", err)
			c.AbortWithError(401, apiErr)
			return
		}

		c.Next()
	}
}

/*  TYPE="service_account"

PROYECT_ID="jopit-334523"

PRIVATE_KEY_ID=1e0b465476ac7c874f58b83f5da81a428450b532

PRIVATE_KEY=-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDvK44Srj8tPMzc\nYNFPH9aAt6WOcQogWtUSiqoicfk74ycr4CfE1tT0FpnDqBsGOpgDL785mpGikVrC\nx8t+bBBNU92htpgy3wSA2bHJseu4vish5kmb+prF9Y56y7nBity/OvNICUWtpui1\ngC7hRfuKD8vpgfB920CmnGtC5v4mbTfeiIGK5sfiK+IJ9S+PrBplv7ollyTppYge\nnGohX0m2gm29cMynWQ/nEku7Npcl3RbFOwI1dRWrJgztdRaAQonPxHbN323m5tWU\n2Jo+0t0YqA3gHs8wDHNR/f2kHbW3x/XakjiDOTZyAW+VyKhCejw0oXv03ttNvkfs\n+df1wSz/AgMBAAECggEAE/QioUUE+PrhQzhB84Oo0LUBofWlA6KwXiClexmJJ/RB\nQ17V0o1IglNptrIitw1SZgmovt86FpcL5xjmIwD/By81BsCGHXsPvsmVafS8bChU\nJDrjRJkkKoSHGOWVn3kSBlOeZ7x7oguSefoHbzyjH2yb9Y+lXHVbWnXNrmzjzUQW\nyf+yEbI37VNlYmtvB0CXoKL+4RRXpnaWUjOBTH8vkOKBtLqlvT8zputqBgMX4gv+\nkpPM5wOJAD5O8HQUXMt82AVx2UVwrx7LakoUbGPa/vaxBa/fCKzIJVTbZgeUXBWY\nr9YwjrM6SueOYroaW9+iy/4iCAako0KpEvjI6gWmEQKBgQD+G+TBScTyX0rTk9/F\n/7qXPASsGSX0HHtFcMC6Eh8fRiB8LEsGOguBlLc1cjF5BA63hZo2fMTZOJbpph89\nT52FHoBdqLEmaNN8ngdvtIllB1hKeAWvu6gt+enOKpTK9jiAdpAARQbF6BK7IuLT\naSN0EZobtF0A4G14flyYmmGRVwKBgQDw8zOEwhgY7FXfQk73GbErtPI/eN5Xo/ek\nS2r/g6PMG4pCgknyEZ17zmV34u4GPzccuIRuYqCI7oKgzO32PYb98SmqjWvdvr5n\nwWsUS0DV/5jF5tDIR7e+trAQcBv6bPEbJZo6YuIgeB3O+59gvBeG3tx0tIdQ11HD\nIq1757cwmQKBgHfoF6ixs7KfDsMJ+0UGKDknBBllgIhRTEX3L6pd1LvALqIJkJlu\nmHFcCJu6t+ef54XoEF33wDe0QVodno4i3HStcOtBJ961Bl6+f0nRyieXMK1rN1u2\nVGTlkdNMOowPZZgQ2mCWMdz/Zp9RlfEmoqkpiSvbjytTu7RpBC6qYrHfAoGBANVW\n0JfMbx+qKgJKqLY7BlbnmNJAR8WbhYnvyOZB/UacX9exmp19oU3cmpUC1bOsRzTj\n30YJh5CWdgzARjYPljXEURsUqpwk5lvhsti+JMMV04PucY/TiEqRYOS9Dti1mulp\natwlb4hGMkZxHlW9VKtzKgkfSux2KYu4aZjauqWJAoGBAMsY1UCco1A3wqJG+XxF\n7Z0ZP/+eVe/3GoK8UUP3S8Ee7NqhANB4JPzDtjBQriO1He+L1s5p3WK6SnKtoBlo\nbkirhzClIzclL2913dpZQmg9o2Y8GsdwgvofvQA5pYyI2Au0uNzuVK+kN9vDcpV6\nXCYgzf8iGg2mX3HrtG5VAXf2\n-----END PRIVATE KEY-----\n"

CLIENT_EMAIL=firebase-adminsdk-c7sfj@jopit-334523.iam.gserviceaccount.com

CLIENT_ID=101930044915607302886

AUTH_URI=https://accounts.google.com/o/oauth2/auth

TOKEN_URI=https://oauth2.googleapis.com/token

AUTH_PROVIDER_x509_CERT_URL=https://www.googleapis.com/oauth2/v1/certs

CLIENT_x509_CERT_URL=https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-c7sfj%40jopit-334523.iam.gserviceaccount.com
*/
