package example_grpc

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/kondohiroki/go-grpc-boilerplate/internal/logger"
	"github.com/kondohiroki/go-grpc-boilerplate/pkg/cache"
	"github.com/kondohiroki/go-grpc-boilerplate/pkg/exception"
	"github.com/kondohiroki/go-grpc-boilerplate/pkg/response"
	"github.com/kondohiroki/go-grpc-boilerplate/pkg/transport"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResult struct {
	SessionId string `json:"session_id"`
}

// authenticate login and obtains a newly session_id from Generic Authentication Service.
func (s *exampleService) authenticate(ctx context.Context) (string, error) {
	logger.Log.Debug(fmt.Sprintf("authenticating using %s:%s", s.email, s.password))

	// prepare request to Core Generic Authentication Service
	params := &LoginParams{
		Email:    s.email,
		Password: s.password,
	}
	payload, err := sonic.Marshal(params)
	if err != nil {
		logger.Log.Error("could not marshal params for calling auth because: %v", zap.Error(err))
		return "", exception.AuthCoreError
	}
	req := transport.HttpRequest{
		HttpClient: transport.NewHTTPClient(),
		Url:        s.authEndpoint,
		Method:     http.MethodPost,
		Body:       payload,
	}

	// read response and its result
	resp := &response.CommonResponse{}
	if err := transport.RequestAndParseJSONBody(ctx, req, resp); err != nil {
		logger.Log.Error("could not request and parse JSON body from auth because: %v", zap.Error(err))
		return "", exception.AuthCoreError
	}
	result := &LoginResult{}
	if err := resp.UnwrapData(result); err != nil {
		logger.Log.Error("could not unwrap auth resp because: %v", zap.Error(err))
		return "", exception.AuthCoreError
	}

	if err := cache.Set(ctx, cache.KEY_EXAMPLE_ACCESS_TOKEN, result.SessionId, 0); err != nil {
		logger.Log.Error("could not set session_id to cache because: %v", zap.Error(err))
		return "", exception.AuthCoreError
	}

	return result.SessionId, nil
}

// obtainSessionId gets session_id from cache. If it was invalid or expired, this func instead
// request authenticate to Generic Authentication Service and save to cache.
func (s *exampleService) obtainSessionId(ctx context.Context) (string, error) {
	sessionId, err := cache.Get(ctx, cache.KEY_EXAMPLE_ACCESS_TOKEN)
	if err != nil {
		if errors.Is(err, redis.Nil) { // token is not exist in cache or it was expired
			sessionId, err = s.authenticate(ctx)
			if err != nil {
				return "", err
			}
			return sessionId, nil
		}
		logger.Log.Error("could not get session_id from cache because: %v", zap.Error(err))
		return "", exception.AuthCoreError
	}

	logger.Log.Debug(fmt.Sprintf("obtained session_id: %s", sessionId))
	return sessionId, nil
}
