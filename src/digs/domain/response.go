package domain


type UserLoginResponse struct {
	StatusCode int32
	SessionId string
	Name string
	Email string
	About string
}

type MessageReceivedResponse struct {
	StatusCode int32
}

type ErrorResponse struct {
	StatusCode int32
	Message string
}