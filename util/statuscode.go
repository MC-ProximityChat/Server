package util

type Status struct {
	Code    int
	Message string
	Series  Series
}

type Series string

func (status Status) isInformational() bool {
	return status.Series == INFORMATIONAL
}

func (status Status) isSuccess() bool {
	return status.Series == SUCCESS
}

func (status Status) isRedirection() bool {
	return status.Series == REDIRECTION
}

func (status Status) isClientError() bool {
	return status.Series == CLIENT
}

func (status Status) isServerError() bool {
	return status.Series == SERVER
}

func (status Status) isError() bool {
	return status.isClientError() || status.isServerError()
}

// SERIES DECLARATION
//goland:noinspection ALL
const (
	INFORMATIONAL Series = "Informational Response"
	SUCCESS       Series = "Success"
	REDIRECTION   Series = "Redirection"
	CLIENT        Series = "Client Error"
	SERVER        Series = "Server Error"
)

// HTTP STATUS DECLARATION
//goland:noinspection ALL
var (
	CONTINUE            = Status{Code: 100, Message: "Continue", Series: INFORMATIONAL}
	SWITCHING_PROTOCOLS = Status{Code: 101, Message: "Switching Protocols", Series: INFORMATIONAL}

	OK                            = Status{Code: 200, Message: "OK", Series: SUCCESS}
	CREATED                       = Status{Code: 201, Message: "Created", Series: SUCCESS}
	ACCEPTED                      = Status{Code: 202, Message: "Accepted", Series: SUCCESS}
	NON_AUTHORITATIVE_INFORMATION = Status{Code: 203, Message: "Non-Authoritative Information", Series: SUCCESS}
	NO_CONTENT                    = Status{Code: 204, Message: "No Content", Series: SUCCESS}
	RESET_CONTENT                 = Status{Code: 205, Message: "Reset Content", Series: SUCCESS}

	MULTIPLE_CHOICES   = Status{Code: 300, Message: "Multiple Choices", Series: REDIRECTION}
	MOVED_PERMANENTLY  = Status{Code: 301, Message: "Moved Permanently", Series: REDIRECTION}
	FOUND              = Status{Code: 302, Message: "Found", Series: REDIRECTION}
	SEE_OTHER          = Status{Code: 303, Message: "See Other", Series: REDIRECTION}
	NOT_MODIFIED       = Status{Code: 304, Message: "Not Modified", Series: REDIRECTION}
	USE_PROXY          = Status{Code: 305, Message: "Use Proxy", Series: REDIRECTION}
	TEMPORARY_REDIRECT = Status{Code: 307, Message: "Temporary Redirect", Series: REDIRECTION}
	PERMANENT_REDIRECT = Status{Code: 308, Message: "Permanent Redirect", Series: REDIRECTION}

	BAD_REQUEST                   = Status{Code: 400, Message: "Bad Request", Series: CLIENT}
	UNAUTHORIZED                  = Status{Code: 401, Message: "Unauthorized", Series: CLIENT}
	PAYMENT_REQUIRED              = Status{Code: 402, Message: "Payment Required", Series: CLIENT}
	FORBIDDEN                     = Status{Code: 403, Message: "Forbidden", Series: CLIENT}
	NOT_FOUND                     = Status{Code: 404, Message: "Not Found", Series: CLIENT}
	METHOD_NOT_ALLOWED            = Status{Code: 405, Message: "Method Not Allowed", Series: CLIENT}
	NOT_ACCEPTABLE                = Status{Code: 406, Message: "Not Acceptable", Series: CLIENT}
	PROXY_AUTHENTICATION_REQUIRED = Status{Code: 407, Message: "Proxy Authentication Required", Series: CLIENT}
	REQUEST_TIMEOUT               = Status{Code: 408, Message: "Request Timeout", Series: CLIENT}
	CONFLICT                      = Status{Code: 409, Message: "Conflict", Series: CLIENT}
	GONE                          = Status{Code: 410, Message: "Gone", Series: CLIENT}
	BLAZE_IT                      = Status{Code: 420, Message: "Blaze It", Series: CLIENT}

	INTERNAL_SERVER_ERROR = Status{Code: 500, Message: "Internal Server Error", Series: SERVER}
	NOT_IMPLEMENTED       = Status{Code: 501, Message: "Not Implemented", Series: SERVER}
	BAD_GATEWAY           = Status{Code: 502, Message: "Bad Gateway", Series: SERVER}
	SERVICE_UNAVAILABLE   = Status{Code: 503, Message: "Service Unavailable", Series: SERVER}
	GATEWAY_TIMEOUT       = Status{Code: 504, Message: "Gateway Timeout", Series: SERVER}
)
