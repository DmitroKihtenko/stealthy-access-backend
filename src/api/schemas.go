package api

type HealthcheckResponse struct {
	Status string `json:"status" example:"ok"`
} //@name HealthcheckResponse

type AddUserRequest struct {
	User
} //@name AddUserRequest

type UserResponse struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
} //@name UserResponse

type GetUsersResponse struct {
	List           []UserResponse `json:"list"`
	FirstPageToken *string        `json:"first_page_token" example:"euKoY1BqY3J8GVax"`
	NextPageToken  *string        `json:"next_page_token" example:"QLux4Tu5gb8JfW70"`
} //@name GetUsersResponse

type ErrorResponse struct {
	Summary string `json:"summary" validate:"required" example:"Invalid authorization token"`
	Detail  any    `json:"detail"`
} //@name ErrorResponse

type PaginationQueryParameters struct {
	PageToken string `query:"page_token" example:"euKoY1BqY3J8GVax" default:""`
	Limit     int64  `validate:"gte=1" query:"limit" example:"20" default:"20"`
} //@name PaginationQueryParameters
