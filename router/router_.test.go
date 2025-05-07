package router

import (
	"httpserver/core"
	"net/url"
	"testing"
)

// helper para crear HttpRequest sin pasar por socket
func makeReq(method, rawpath string) *core.HttpRequest {
	u, _ := url.Parse(rawpath)
	return core.NewHttpRequest(method, u, map[string]string{}, "")
}

func TestRouterBasic(t *testing.T) {
	r := New()
	r.Get("/foo", func(req *core.HttpRequest) (*core.HttpResponse, error) {
		return core.Ok().Text("OK"), nil
	})

	res, _ := r.Handle(makeReq("GET", "/foo"))
	if res.StatusCode != 200 || res.Body != "OK" {
		t.Fatalf("esperaba 200 OK, got %d %q", res.StatusCode, res.Body)
	}

	res404, _ := r.Handle(makeReq("GET", "/nope"))
	if res404.StatusCode != 404 {
		t.Errorf("esperaba 404 para ruta no registrada")
	}
}
