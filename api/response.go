package api

// errorResponse is a typed error response that avoids fiber.Map JSON encoding issues
type errorResponse struct {
	Error string `json:"error"`
}

// errorResponseWithIP includes the source IP in the error response
type errorResponseWithIP struct {
	Error string `json:"error"`
	IP    string `json:"ip"`
}

// errorResponseWithLimit includes limit details in the error response
type errorResponseWithLimit struct {
	Error     string `json:"error"`
	Maximum   int    `json:"maximum"`
	Requested int    `json:"requested"`
}

// errorResponseWithURL includes URL details in the error response
type errorResponseWithURL struct {
	Error string `json:"error"`
	Index int    `json:"index"`
	URL   string `json:"url"`
}
