package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
)

type HttpRequest struct {
	Method string
	Target *url.URL
	Header map[string]string
	Body   string
}

func NewHttpRequest(method string, target *url.URL, header map[string]string, body string) *HttpRequest {
	return &HttpRequest{
		Method: method,
		Target: target,
		Header: header,
		Body:   body,
	}
}

func ReadRequest(conn net.Conn) (*HttpRequest, error) {
	lines := make([]string, 0)

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, err
		}

		line = strings.TrimSuffix(line, "\r\n")
		line = strings.TrimSuffix(line, "\n")

		if line == "" {
			break
		}

		lines = append(lines, line)
	}

	if len(lines) == 0 {
		return nil, fmt.Errorf("empty request")
	}

	headersPart := strings.Join(lines, "\n") + "\n\n"

	request, err := ParseRequest(headersPart)
	if err != nil {
		return nil, fmt.Errorf("can't parse request: %w", err)
	}

	if contentLengthStr, ok := request.Header["Content-Length"]; ok {
		var contentLength int

		_, err := fmt.Sscan(contentLengthStr, &contentLength)
		if err != nil {
			return nil, fmt.Errorf("bad content length format: %w", err)
		}

		if contentLength > 0 {
			body := make([]byte, contentLength)

			_, err := io.ReadAtLeast(reader, body, contentLength)
			if err != nil {
				return nil, fmt.Errorf("can't read body: %w", err)
			}

			request.Body = string(body)
		}
	}

	return request, nil
}

func ParseRequest(headersPart string) (*HttpRequest, error) {
	headerEndIndex := strings.Index(headersPart, "\n\n")
	if headerEndIndex == -1 {
		return nil, fmt.Errorf("no header end")
	}

	headerPart := headersPart[:headerEndIndex]

	lines := strings.Split(headerPart, "\n")
	if len(lines) < 1 {
		return nil, fmt.Errorf("no start line")
	}

	start := strings.SplitN(lines[0], " ", 3)
	if len(start) < 2 {
		return nil, fmt.Errorf("no method or target")
	}

	method := start[0]
	if method == "" {
		return nil, fmt.Errorf("no method")
	}

	targetStr := start[1]
	if targetStr == "" {
		return nil, fmt.Errorf("no target")
	}

	target, err := url.Parse(targetStr)
	if err != nil {
		return nil, fmt.Errorf("bad target format: %w", err)
	}

	headers := make(map[string]string)

	for _, line := range lines[1:] {
		parts := strings.SplitN(line, ":", 2)

		if len(parts) != 2 {
			continue
		}

		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		if k == "" {
			continue
		}

		headers[k] = v
	}

	body := ""

	return NewHttpRequest(method, target, headers, body), nil
}
