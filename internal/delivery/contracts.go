package delivery

import "context"

type HttpService interface {
	Start(address string) error
	Shutdown(ctx context.Context) error
}
