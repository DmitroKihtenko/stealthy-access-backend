package base

type SchemaProperty string

const ConfigFile string = "config.yaml"
const LimitQueryParam string = "limit"
const PageTokenQueryParam string = "page_token"
const UserIdPathParam string = "user_id"

const UserSchemaId string = "user"
const PaginationHeader string = "Link"

const (
	Username  SchemaProperty = "username"
	Email     SchemaProperty = "email"
	FirstName SchemaProperty = "firstname"
	LastName  SchemaProperty = "lastname"
)
