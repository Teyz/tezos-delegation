package poller_v1

import (
	"context"
)

type TezosPoller interface {
	StartPolling(ctx context.Context) error
}
