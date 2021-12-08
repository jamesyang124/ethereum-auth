package handlers

import (
	"context"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func Test_responseErrorLogging(t *testing.T) {
	oneErr := fiber.NewError(1, "error code 1")
	lg := log.New(os.Stdout, "15:04:05 | ", 0)
	type args struct {
		code int
		errs []error
		l    *log.Logger
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "no-errs", args: args{code: 200, errs: []error{}, l: lg}},
		{name: "one-err", args: args{code: 401, errs: []error{oneErr}, l: lg}},
		{name: "many-errs", args: args{code: 500, errs: []error{oneErr, fiber.NewError(45021, "error code 45021")}, l: lg}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseErrorLogging(tt.args.code, tt.args.errs, tt.args.l)
		})
	}
}

// TODO: need mock server for this downstram testing
// TODO: https://pkg.go.dev/net/http/httptest
func Test_downstreamAuthRequest(t *testing.T) {
	type args struct {
		url   string
		extra map[string]interface{}
		cid   string
		nid   string
		paddr string
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 []byte
		want2 []error
	}{
		// TODO: Add test cases.
		{name: "case-1", args: args{url: "http://1.2.3.4", extra: map[string]interface{}{}, cid: "chain-id-1", nid: "network-id-1", paddr: "public-address-1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := downstreamAuthRequest(tt.args.url, tt.args.extra, tt.args.cid, tt.args.nid, tt.args.paddr)
			if got != tt.want {
				t.Errorf("downstreamAuthRequest() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("downstreamAuthRequest() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("downstreamAuthRequest() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

// TODO: may skip this part since we are expecting to test route instead of its handler impl.
func TestLoginHandler(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(nil)
	type args struct {
		ctx                context.Context
		rdb                *redis.Client
		l                  *log.Logger
		redisTTL           string
		signInTextTemplate string
		downstreamAuthUri  string
	}
	tests := []struct {
		name string
		args args
		want func(c *fiber.Ctx) error
	}{
		{name: "first-case", args: args{}, want: func(c *fiber.Ctx) error { return nil }},
		{name: "first-case", args: args{}, want: func(c *fiber.Ctx) error { return nil }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoginHandler(tt.args.ctx, tt.args.rdb, tt.args.l, tt.args.redisTTL, tt.args.signInTextTemplate, tt.args.downstreamAuthUri); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoginHandler() = %v, want %v", got(ctx), tt.want(ctx))
			}
		})
	}
}

//TODO: add not unit test level code
