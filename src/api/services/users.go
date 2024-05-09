package services

import (
	"access-backend/api"
	"access-backend/base"
	"context"
	ory "github.com/ory/kratos-client-go"
	"github.com/sirupsen/logrus"
	"github.com/tomnomnom/linkheader"
	"net/http"
	"net/url"
)

func kratosIdentityToUser(identity *ory.Identity) *api.UserResponse {
	defer func() {
		if err := recover(); err != nil {
			base.Logger.WithFields(logrus.Fields{
				"error":       err,
				"identity_id": identity.Id,
			}).Warn("Error parsing user from identity")
		}
	}()

	traits := identity.Traits.(map[string]interface{})

	return &api.UserResponse{
		Id:        identity.Id,
		Username:  traits[string(base.Username)].(string),
		Email:     traits[string(base.Email)].(string),
		FirstName: traits[string(base.FirstName)].(string),
		LastName:  traits[string(base.LastName)].(string),
	}
}

type BaseUserService interface {
	AddUser(request *api.AddUserRequest) (*api.UserResponse, error)
	GetUsers(request *api.PaginationQueryParameters) (
		*api.GetUsersResponse, error,
	)
	DeleteUser(userId string) error
}

type UserService struct {
	BaseUserService
	Context      *context.Context
	KratosClient *ory.APIClient
}

func (service *UserService) AddUser(request *api.AddUserRequest) (
	*api.UserResponse, error,
) {
	identityBody := ory.CreateIdentityBody{
		SchemaId: base.UserSchemaId,
		Credentials: &ory.IdentityWithCredentials{
			Password: &ory.IdentityWithCredentialsPassword{
				Config: &ory.IdentityWithCredentialsPasswordConfig{
					Password: &request.Password,
				},
			},
		},
		Traits: map[string]interface{}{
			string(base.Username):  request.Username,
			string(base.Email):     request.Email,
			string(base.FirstName): request.FirstName,
			string(base.LastName):  request.LastName,
		},
	}

	identity, response, err := service.KratosClient.IdentityAPI.CreateIdentity(
		*service.Context,
	).CreateIdentityBody(identityBody).Execute()
	if err != nil {
		if response.StatusCode == http.StatusConflict {
			return nil, base.ServiceError{
				Summary: "User with username '" + request.Username +
					"' already exist",
				Status: http.StatusConflict,
			}
		}

		return nil, base.NewKratosError(
			"Error creating new user",
			err,
		)
	}

	result := kratosIdentityToUser(identity)
	if result != nil {
		return result, nil
	} else {
		return nil, base.NewKratosError("User data not created", nil)
	}
}

func getPageTokenFromUrl(urlString string) *string {
	urlObj, _ := url.Parse(urlString)

	if urlObj != nil {
		values, _ := url.ParseQuery(urlObj.RawQuery)
		if values != nil {
			value := values.Get("page_token")
			if value != "" {
				return &value
			}
		}
	}

	return nil
}

func (service *UserService) GetUsers(request *api.PaginationQueryParameters) (
	*api.GetUsersResponse, error,
) {
	result := api.GetUsersResponse{}

	kratosRequest := service.KratosClient.IdentityAPI.ListIdentities(
		*service.Context,
	).PageSize(request.Limit).PageToken(request.PageToken)
	identities, response, err := kratosRequest.Execute()

	if err != nil {
		return nil, base.NewKratosError("Error retrieving users", err)
	}
	users := make([]api.UserResponse, 0, len(identities))
	for _, identity := range identities {
		users = append(users, *kratosIdentityToUser(&identity))
	}

	linkHeader := response.Header.Get(base.PaginationHeader)
	if linkHeader != "" {
		links := linkheader.Parse(linkHeader)
		if len(links) == 0 {
			base.Logger.Warn(
				"No links in '" + base.PaginationHeader + "' header",
			)
		}

		for _, link := range links {
			if link.Rel == "first" {
				result.FirstPageToken = getPageTokenFromUrl(link.URL)
			}
			if link.Rel == "next" {
				result.NextPageToken = getPageTokenFromUrl(link.URL)
			}
		}
	} else {
		base.Logger.Warn(
			"'" + base.PaginationHeader + "' header not found in Kratos response",
		)
	}
	result.List = users

	return &result, nil
}

func (service *UserService) DeleteUser(userId string) error {
	response, err := service.KratosClient.IdentityAPI.DeleteIdentity(
		*service.Context, userId,
	).Execute()

	if err != nil {
		if response.StatusCode == http.StatusNotFound {
			return base.ServiceError{
				Summary: "User with id '" + userId + "' not found",
				Status:  http.StatusBadRequest,
			}
		} else {
			return base.NewKratosError(
				"Error deleting user",
				err,
			)
		}
	}
	return nil
}
