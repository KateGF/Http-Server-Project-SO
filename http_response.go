package main

import "fmt"

type HttpResponse struct {
	StatusCode int
	StatusText string
	Headers    map[string]string
	Body       string
}

func NewHttpResponse(statusCode int, statusText string, body string) *HttpResponse {
	return &HttpResponse{
		StatusCode: statusCode,
		StatusText: statusText,
		Headers:    make(map[string]string),
		Body:       body,
	}
}

func Ok() *HttpResponse {
	return NewHttpResponse(200, "OK", "")
}

func NotFound() *HttpResponse {
	return NewHttpResponse(404, "Not Found", "")
}

func (response *HttpResponse) SetStatusCode(code int) *HttpResponse {
	response.StatusCode = code
	return response
}

func (response *HttpResponse) SetStatusText(text string) *HttpResponse {
	response.StatusText = text
	return response
}

func (response *HttpResponse) SetHeader(key, value string) *HttpResponse {
	response.Headers[key] = value
	return response
}

func (response *HttpResponse) SetBody(body string) *HttpResponse {
	response.Body = body
	return response
}

func (response *HttpResponse) SetContentType(contentType string) *HttpResponse {
	response.SetHeader("Content-Type", contentType)
	return response
}

func (response *HttpResponse) Text(text string) *HttpResponse {
	response.SetContentType("text/plain")
	response.SetBody(text)
	return response
}

func (response *HttpResponse) Json(json string) *HttpResponse {
	response.SetContentType("application/json")
	response.SetBody(json)
	return response
}

func (response *HttpResponse) String() string {
	contentLength := len(response.Body)
	response.SetHeader("Content-Length", fmt.Sprint(contentLength))

	headersStr := ""
	for key, value := range response.Headers {
		headersStr += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	return fmt.Sprintf("HTTP/1.1 %d %s\r\n%s\r\n%s", response.StatusCode, response.StatusText, headersStr, response.Body)
}
