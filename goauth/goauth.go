package goauth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	jsonLib "github.com/json-iterator/go"
	"github.com/matiasnu/go-jopit-toolkit/goutils/apierrors"
	"github.com/matiasnu/go-jopit-toolkit/rest"
	"github.com/matiasnu/go-jopit-toolkit/rest/retry"
)

const BASE_URL string = "http://internal.jopit.com"

type authRequestData struct {
	UserId     *string `json:"user_id"`
	Status     *string `json:"status"`
	AdminId    *string `json:"admin_id"`
	ClientId   *int64  `json:"client_id"`
	IsTest     *bool   `json:"is_test"`
	OperatorID *int    `json:"operator_id"`
	DetachedId *string `json:"detached_id"`
	RootId     *int64  `json:"root_id"`
}

var (
	restClient     *rest.RequestBuilder
	privateParams  = [8]string{"caller.id", "caller.scopes", "caller.status", "client.id", "admin.id", "caller.siteId", "operator.id", "root.id"}
	privateHeaders = [10]string{"X-Caller-Id", "X-Caller-Scopes", "X-Caller-Status", "X-Client-Id", "X-Test-Token", "X-Admin-Id", "X-Caller-SiteId", "X-Operator-Id", "X-Detached-Id", "X-Root-Id"}
	useMock        bool
)

type authOptions struct {
	allowNonActiveUser bool
}

type AuthOption func(ao *authOptions)

func AuthenticateRequest(request *http.Request) error {
	return AuthenticateRequestWithOptions(request)
}

// AllowNonActiveUser lets you indicate whether a user with a status different than active would be allowed to authenticate.
func AllowNonActiveUser(allow bool) AuthOption {
	return func(ao *authOptions) {
		ao.allowNonActiveUser = allow
	}
}

func AuthenticateRequestWithOptions(request *http.Request, opts ...AuthOption) error {
	if useMock {
		return nil
	}

	if IsHandledByMiddleware(request) {
		return nil
	}

	isPublic := IsPublic(request)
	if isPublic {
		cleanRequest(request)
	}

	accessToken := request.URL.Query().Get("access_token")
	if accessToken == "" {
		return nil
	}

	var authOptions authOptions
	for _, opt := range opts {
		opt(&authOptions)
	}

	url := fmt.Sprintf("%s/auth/access_token", BASE_URL)

	headers := make(http.Header)
	headers.Set("X-Access-Token", accessToken)

	response := restClient.Get(url, rest.Context(request.Context()), rest.Headers(headers))
	if err := response.Err; err != nil {
		if isPublic {
			err = errors.New("internal network error")
		}
		return apierrors.NewInternalServerApiError("Error validating access token", err)
	}

	if response.StatusCode == http.StatusOK {
		authData := &authRequestData{}
		if marshalError := jsonLib.Unmarshal(response.Bytes(), authData); marshalError != nil {
			return apierrors.NewInternalServerApiError("Invalid json response calling auth api", nil)
		}

		if err := validateRequest(authData, authOptions); err != nil {
			return err
		}

		addParametersToRequest(request, authData)
		addHeadersToRequest(request, authData)

		return nil
	}

	if response.StatusCode == http.StatusNotFound {
		return apierrors.NewApiError("invalid_token", "not_found", http.StatusUnauthorized, apierrors.CauseList{})
	}

	// Unknown status code returned by rest API call.
	message := response.String()
	if isPublic {
		message = "unknown_error"
	}
	return apierrors.NewApiError(message, "", response.StatusCode, apierrors.CauseList{})
}

func validateRequest(data *authRequestData, authOptions authOptions) error {

	if data.UserId == nil {
		return apierrors.NewForbiddenApiError("Invalid user id")
	}

	if data.RootId == nil {
		return apierrors.NewForbiddenApiError("Invalid root id")
	}

	if data.Status == nil || (!authOptions.allowNonActiveUser && *(data.Status) != "active") {
		return apierrors.NewForbiddenApiError("User not active")
	}

	return nil
}

func addParametersToRequest(request *http.Request, data *authRequestData) {
	query := request.URL.Query()

	query.Add("caller.id", *data.UserId)
	query.Add("caller.status", "ACTIVE")
	query.Add("root.id", fmt.Sprint(*data.RootId))

	if data.ClientId != nil {
		query.Add("client.id", strconv.FormatInt(*data.ClientId, 10))
	}

	if data.AdminId != nil {
		query.Add("admin.id", *data.AdminId)
	}

	if data.OperatorID != nil {
		query.Add("operator.id", strconv.Itoa(*data.OperatorID))
	}

	request.URL.RawQuery = query.Encode()
}

func addHeadersToRequest(request *http.Request, data *authRequestData) {
	request.Header.Set("X-Caller-Id", *data.UserId)
	request.Header.Set("X-Caller-Status", "ACTIVE")

	if data.IsTest != nil {
		request.Header.Set("X-Test-Token", strconv.FormatBool(*data.IsTest))
	}

	if data.ClientId != nil {
		request.Header.Set("X-Client-Id", strconv.FormatInt(*data.ClientId, 10))
	}

	if data.AdminId != nil {
		request.Header.Set("X-Admin-Id", *data.AdminId)
	}

	if data.OperatorID != nil {
		request.Header.Set("X-Operator-Id", strconv.Itoa(*data.OperatorID))
	}

	if data.DetachedId != nil {
		request.Header.Set("X-Detached-Id", *data.DetachedId)
	}

	if data.RootId != nil {
		request.Header.Set("X-Root-Id", fmt.Sprint(*data.RootId))
	}
}

func cleanRequest(request *http.Request) {
	query := request.URL.Query()

	for i := 0; i < len(privateParams); i++ {
		query.Del(privateParams[i])
	}

	for i := 0; i < len(privateHeaders); i++ {
		request.Header.Del(privateHeaders[i])
	}

	request.URL.RawQuery = query.Encode()
}

func init() {
	customPool := &rest.CustomPool{
		MaxIdleConnsPerHost: 100,
	}

	restClient = &rest.RequestBuilder{
		Timeout:        500 * time.Millisecond,
		ContentType:    rest.JSON,
		DisableTimeout: false,
		RetryStrategy:  retry.NewSimpleRetryStrategy(2, 20*time.Millisecond),
		CustomPool:     customPool,
		MetricsConfig:  rest.MetricsReportConfig{TargetId: "auth-api"},
	}

	useMock = !(os.Getenv("GO_ENVIRONMENT") == "production")
}

func PasswordMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		if os.Getenv("ADMIN_PASSWORD") == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "password is not setted in the repository")
			return
		}

		if os.Getenv("ADMIN_USERNAME") == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "username is not setted in the repository")
			return
		}

		headerUsername := c.GetHeader("admin_username")
		if headerUsername == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "username is empty, please provide one")
			return
		}

		headerPassword := c.GetHeader("admin_password")
		if headerPassword == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "password is empty, please provide one")
			return
		}

		c.Set("admin_username", headerUsername)
		c.Next()
	}
}
