package http

type RequestIdReponse struct {
	Id int `json:"id"`
}

type Response struct {
	Message string `json:"message"`
}

type Error struct {
	Error string `json:"error"`
}

func newError(message string) Error {
	return Error{
		Error: message,
	}
}

func newResponse(message string) Response {
	return Response{
		Message: message,
	}
}
