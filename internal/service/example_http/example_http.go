package example_http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/kondohiroki/go-grpc-boilerplate/config"
	"github.com/kondohiroki/go-grpc-boilerplate/internal/logger"
	"github.com/kondohiroki/go-grpc-boilerplate/pkg/response"
	"github.com/kondohiroki/go-grpc-boilerplate/pkg/transport"
	"go.uber.org/zap"
)

type ExampleService interface {
}

type exampleService struct {
	endpoint string

	// authentication
	authEndpoint string
	email        string
	password     string
}

func NewExampleService(conf *config.Config) ExampleService {
	return &exampleService{
		endpoint:     conf.Services.Example.Endpoint,
		authEndpoint: conf.Services.Example.Authentication.Endpoint,
		email:        conf.Services.Example.Authentication.Username,
		password:     conf.Services.Example.Authentication.Password,
	}
}

type ExampleGRPCParams interface {
	ApiPath() string
}

func (s *exampleService) RequestHTTP(
	ctx context.Context,
	method string,
	params ExampleGRPCParams,
	result interface{},
	query map[string]string,
) (*http.Response, []byte, error) {
	payload, err := sonic.Marshal(params)
	if err != nil {
		return nil, nil, err
	}

	req := transport.HttpRequest{
		HttpClient: transport.NewHTTPClient(),
		Url:        s.endpoint + params.ApiPath(),
		Method:     method,
		Body:       payload,
	}

	if query != nil {
		req.Query = query
	}

	// Authentication
	// sessionId, err := s.obtainSessionId(ctx)
	// if err != nil {
	// 	return err
	// }
	// req.WithBearer(sessionId, s.authenticate)
	resp := &response.CommonResponse{}
	httpResp, respBody, err := transport.RequestAutoBodyParser(ctx, req, resp)
	if err != nil {
		return httpResp, respBody, fmt.Errorf("requestAutoBodyParser error %s: %w", params.ApiPath(), err)
	}

	if result != nil {
		if err := resp.UnwrapData(result); err != nil {
			logger.Log.Error("could not unwrap data", zap.Any("data", resp))
			return httpResp, respBody, fmt.Errorf("could not unwrap data: %w", err)
		}
	}

	return httpResp, respBody, nil
}
