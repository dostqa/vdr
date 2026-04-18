package handlers

type Response struct {
	Message string `json:"message"`
}

func newOK() Response {
	return Response{
		Message: "OK",
	}
}

func newError(message string) Response {
	return Response{
		Message: message,
	}
}
