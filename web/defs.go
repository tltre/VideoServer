package main

type ApiBody struct {
	Url     string `json:"url"`
	Method  string `json:"method"`
	ReqBody string `json:"req_body"`
}

type Err struct {
	Error     string `json:"error"`
	ErrorCode string `json:"error_code"`
}

var (
	ErrorRequestNotRecognize = Err{
		Error:     "API not recognize, bad request",
		ErrorCode: "001",
	}

	ErrorRequestBodyParseFailed = Err{
		Error:     "Request Body not correct",
		ErrorCode: "002",
	}

	ErrorInternalFaults = Err{
		Error:     "Internal service error",
		ErrorCode: "003",
	}
)

const StreamServerUrl = "http://127.0.0.1:9000"
