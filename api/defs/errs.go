package defs

type Err struct {
	Error     string `json:"error"`
	ErrorCode string `json:"errorCode"`
}

type ErrorResponse struct {
	HttpSC int
	Error  Err
}

var (
	ErrorRequestBodyParseFailed = ErrorResponse{
		HttpSC: 400,
		Error:  Err{Error: "Request Body Not Correct", ErrorCode: "001"},
	}

	ErrorNotAuthUser = ErrorResponse{
		HttpSC: 401,
		Error:  Err{Error: "User authentication fail", ErrorCode: "002"},
	}

	ErrorDBError = ErrorResponse{
		HttpSC: 500,
		Error:  Err{Error: "DB ops failed", ErrorCode: "003"},
	}

	ErrorInternalFaults = ErrorResponse{
		HttpSC: 500,
		Error:  Err{Error: "Internal service error", ErrorCode: "004"},
	}

	ErrorNoSuchUser = ErrorResponse{
		HttpSC: 401,
		Error:  Err{Error: "No Such User, Please Create one", ErrorCode: "005"},
	}
)
