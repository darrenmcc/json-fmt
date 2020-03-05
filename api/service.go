package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/darrenmcc/gizmo"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

type service struct {
	sec string
}

func NewService(secret string) gizmo.Service {
	return &service{
		sec: secret,
	}
}

func (s service) HTTPEndpoints() map[string]map[string]gizmo.HTTPEndpoint {
	return map[string]map[string]gizmo.HTTPEndpoint{
		"/{url:.*}": {
			"GET": {
				Endpoint: s.Fmt,
				Encoder:  fmtEncoder,
			},
		},
	}
}

func (s service) Fmt(ctx context.Context, req interface{}) (interface{}, error) {
	r := req.(*http.Request)

	if r.URL.Query().Get("sec") != s.sec {
		return nil, gizmo.NewJSONStatusResponse("now way jose", 401)
	}

	u := strings.Replace(mux.Vars(r)["url"], "https:/", "https://", 1)
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got status %d from downstream", resp.StatusCode)
	}

	var m map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func fmtEncoder(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	b, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

func (s service) HTTPMiddleware(h http.Handler) http.Handler       { return h }
func (s service) Middleware(e endpoint.Endpoint) endpoint.Endpoint { return e }
func (s service) HTTPRouterOptions() []gizmo.RouterOption          { return nil }
func (s service) HTTPOptions() []httptransport.ServerOption        { return nil }
func (s service) RPCMiddleware() grpc.UnaryServerInterceptor       { return nil }
func (s service) RPCOptions() []grpc.ServerOption                  { return nil }
func (s service) RPCServiceDesc() *grpc.ServiceDesc                { return nil }
