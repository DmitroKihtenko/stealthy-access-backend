package main

import (
	"access-backend/api"
	"access-backend/api/controllers"
	"access-backend/api/services"
	"access-backend/base"
	"access-backend/docs"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	ory "github.com/ory/kratos-client-go"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
)

// @title Stealthy Access Backend
// @version 1.0.0
// @schemes http
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey User
// @In header
// @Name Access token
// @Description Access token of authenticated user

func processError(err error) {
	base.Logger.WithFields(logrus.Fields{
		"error": err.Error(),
	}).Fatal("Exiting due to fatal error")
	os.Exit(1)
}

func processPanic() {
	if r := recover(); r != nil {
		err, ok := r.(error)
		if ok {
			base.Logger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Fatal("Exiting due to fatal error")
		} else {
			base.Logger.WithFields(logrus.Fields{
				"detail": r,
			}).Fatal("Exiting due to fatal error")
		}
		os.Exit(1)
	}
}

func configureSwagger(swaggerRouter *gin.RouterGroup, config *base.BackendConfig) {
	base.Logger.Info("Configuring openapi")

	baseUrl := config.Server.Socket + config.Server.BasePath

	docs.SwaggerInfo.Host = config.Server.Socket
	docs.SwaggerInfo.BasePath = config.Server.BasePath
	docs.SwaggerInfo.Description = "Stealthy backend service. " +
		"REST API web application. Encapsulates user's service " +
		"business logic of Stealthy system." +
		"<br><br>API is based on JSON (JavaScript Object Notation) Web " +
		"Application Services and HTTPS transport, so is accessible from " +
		"any platform or operating system. Connection to the JSON API is " +
		"provided via HTTP/HTTPS. Authorization is performed using an " +
		"authorization access token. The example below illustrates " +
		"\"Get users list\" request with an access token: " +
		"<br><strong>curl -X GET " + baseUrl + "/v1/users " +
		"-H \"Authorization: Bearer " +
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9\"</strong>" +
		"<h4>How To Get Authorization Data</h4>" +
		"An authorization access token is static and installed in server " +
		"configuration file. So it is known only by a person who launched " +
		"the service."

	swaggerRouter.GET(
		config.Server.OpenapiBasePath+"/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)
}

func setLogger(config *base.BackendConfig) {
	base.Logger = base.CreateLogger(config)
}

func runServer(engine *gin.Engine, config *base.BackendConfig) {
	base.Logger.Info("Starting server")
	if err := engine.Run(config.Server.Socket); err != nil {
		panic(err)
	}
	base.Logger.Info("Server stopped")
}

func createKratosClient(config *base.BackendConfig) *ory.APIClient {
	serverConfig := ory.NewConfiguration()
	serverConfig.Servers = []ory.ServerConfiguration{
		{URL: config.Kratos.AdminApiUrl},
	}
	return ory.NewAPIClient(serverConfig)
}

func main() {
	defer processPanic()

	schemaValidator := base.CreateValidator()
	config, err := base.LoadConfiguration(base.ConfigFile)
	if err != nil {
		processError(err)
	}

	setLogger(config)

	client := createKratosClient(config)
	contextObject := context.TODO()

	authController := controllers.AuthController{
		AuthService: &services.AuthService{AuthConfig: &config.Auth},
	}
	userController := controllers.UserController{
		Service: &services.UserService{
			Context:      &contextObject,
			KratosClient: client,
		},
		SchemaValidator: schemaValidator,
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.NoRoute(api.NoRouteHandler)
	router.NoMethod(api.NoMethodHandler)
	router.Use(api.LogsHandler)
	router.Use(api.ErrorHandler)
	router.Use(api.CORSHandler)

	applicationGroup := router.Group(config.Server.BasePath)
	v1 := applicationGroup.Group("/v1")

	v1.GET("/health", controllers.CheckHealth)

	usersGroup := v1.Group("/users").Use(authController.Authorize)
	usersGroup.POST("", userController.AddUser)
	usersGroup.GET("", userController.GetUsers)
	usersGroup.DELETE(
		fmt.Sprintf("/:%s", base.UserIdPathParam),
		userController.DeleteUser,
	)

	configureSwagger(applicationGroup, config)

	runServer(router, config)
}
