package handlers

type Response struct {
	Message string `json:"message"`
}

type RequestIdReponse struct {
	Id int `json:"id"`
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
