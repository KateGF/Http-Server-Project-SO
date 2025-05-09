package handlers

import (
	"httpserver/core"
	"net/url"
	"testing"
)

func makeReq(method, raw string) *core.HttpRequest {
	u, _ := url.Parse(raw)
	return core.NewHttpRequest(method, u, map[string]string{}, "")
}

func TestReverse(t *testing.T) {
	res, _ := ReverseHandler(makeReq("GET", "/reverse?text=hola"))
	if res.Body != "aloh" {
		t.Errorf("Reverse incorrecto: %q", res.Body)
	}
	res2, _ := ReverseHandler(makeReq("GET", "/reverse"))
	if res2.StatusCode != 400 {
		t.Errorf("esperaba BadRequest si falta text")
	}
}

func TestToUpper(t *testing.T) {
	res, _ := ToUpperHandler(makeReq("GET", "/toupper?text=GoLang"))
	if res.Body != "GOLANG" {
		t.Errorf("ToUpper incorrecto: %q", res.Body)
	}
}

func TestHash(t *testing.T) {
	res, _ := HashHandler(makeReq("GET", "/hash?text=abc"))
	// SHA256("abc") = "..."
	want := "..." // pon aquí el valor literal
	if res.Body != want {
		t.Errorf("Hash incorrecto: got %q, want %q", res.Body, want)
	}
}
