package controllers

import (
	"access-backend/api"
	"access-backend/api/services"
	"access-backend/base"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
)

type UserController struct {
	Service         services.BaseUserService
	SchemaValidator *validator.Validate
}

// AddUser Add new user godoc
// @Summary      Add new user
// @Description  This method adds a new user
// @Tags         Users
// @Security     User
// @Accept       json
// @Produce      json
// @Param   	 request  body  api.AddUserRequest true "User sign-up schema"
// @Success      201  {object}  api.UserResponse
// @Failure      400  {object}  api.ErrorResponse
// @Failure      422  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /v1/users [post]
func (controller UserController) AddUser(c *gin.Context) {
	base.Logger.Info("Requested creating user")

	var request api.AddUserRequest
	if err := c.BindJSON(&request); err != nil {
		return
	}

	err := controller.SchemaValidator.Struct(request)
	if err != nil {
		c.Error(base.WrapValidationErrors(err))
		return
	}

	user, err := controller.Service.AddUser(&request)
	if err != nil {
		c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusCreated, user)
}

// GetUsers Get users list godoc
// @Summary      Get list of users
// @Description  This method returns list of users
// @Tags         Users
// @Security     User
// @Accept       json
// @Produce      json
// @Param 		 _ 	  query     api.PaginationQueryParameters false "Pagination parameters"
// @Success      201  {object}  api.GetUsersResponse
// @Failure      400  {object}  api.ErrorResponse
// @Failure      422  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /v1/users [get]
func (controller UserController) GetUsers(c *gin.Context) {
	base.Logger.Info("Requested list of users")

	queryParams := api.PaginationQueryParameters{}

	pageToken := c.DefaultQuery(base.PageTokenQueryParam, "")
	val := strconv.FormatInt(20, 10)
	limit, err := strconv.ParseInt(
		c.DefaultQuery(base.LimitQueryParam, val), 10, 64)
	if err != nil {
		c.Error(base.NewPathParamError(base.LimitQueryParam, err))
		return
	}

	queryParams.PageToken = pageToken
	queryParams.Limit = limit
	if err = controller.SchemaValidator.Struct(queryParams); err != nil {
		c.Error(base.WrapValidationErrors(err))
		return
	}

	response, err := controller.Service.GetUsers(&queryParams)
	if err != nil {
		c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusOK, &response)
}

// DeleteUser Delete user godoc
// @Summary      Delete user by id
// @Description  This method removes user
// @Tags         Users
// @Security     User
// @Accept       json
// @Produce      json
// @Param 		 id path string true "User id" example(6e98ca78-d3ea-4682-adf1-51c12585e7d7)
// @Success      204
// @Failure      400  {object}  api.ErrorResponse
// @Failure      422  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /v1/users/{user_id} [delete]
func (controller UserController) DeleteUser(c *gin.Context) {
	base.Logger.Info("Requested deleting user")

	fileId := c.Param(base.UserIdPathParam)
	if fileId == "" {
		c.Error(base.NewPathParamRequiredError(base.UserIdPathParam))
		return
	}

	err := controller.Service.DeleteUser(fileId)
	if err != nil {
		c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusNoContent, nil)
}
