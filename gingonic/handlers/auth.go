package handlers

import (
	"context"
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

type FirebaseClient struct {
	AuthClient *auth.Client
}

//initiates the firebase client ONCE
func Newfirebase() *FirebaseClient {
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
	InitFirebase()
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
