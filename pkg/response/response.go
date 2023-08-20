package response

import (
	"github.com/bytedance/sonic"
	"github.com/kondohiroki/go-grpc-boilerplate/pkg/exception"
	"github.com/kondohiroki/go-grpc-boilerplate/pkg/pagination"
)

type DataUnwrapper interface {
	UnwrapData(interface{}) error
}

// Standard Response
type CommonResponse struct {
	ResponseCode    int                        `json:"response_code"`
	ResponseMessage string                     `json:"response_message"`
	Errors          *exception.ExceptionErrors `json:"errors,omitempty"`
	Pagination      *pagination.PaginationDTO  `json:"pagination,omitempty"`
	Data            any                        `json:"data,omitempty"`
	RequestID       string                     `json:"request_id,omitempty"`
	LastUpdatedAt   string                     `json:"last_updated_at,omitempty"`
}

func (resp *CommonResponse) UnwrapData(target interface{}) error {
	bs, err := sonic.Marshal(resp.Data)
	if err != nil {
		return err
	}

	if err := sonic.Unmarshal(bs, target); err != nil {
		return err
	}

	return nil
}
