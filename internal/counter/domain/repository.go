package counter 

import (
	"context"
)

type Repository interface {
	// maybe someday, reseting formater in Counter is necessary
	AddCounter(ctx context.Context, c *Counter) error
	GetNextFromCounter(
		ctx context.Context,
		counterUUID string,
	) (string, error)
	ResetCounter(
		ctx context.Context,
		counterUUID string,
	) error

}