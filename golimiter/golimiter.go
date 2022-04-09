package golimiter

import (
	"errors"
	"time"

	"github.com/matiasnu/go-jopit-toolkit/golimiter/node"
)

var OverQuotaError = errors.New("over quota")

type Limiter struct {
	maxRPM uint64
	period time.Duration
	node   *node.TokenRateNode
}

func New(maxRPM uint64, period time.Duration) *Limiter {
	return &Limiter{maxRPM, period, node.New(maxRPM, uint64(period.Nanoseconds()/1000000))}
}

func (l *Limiter) Action(weight uint64, f func() (interface{}, error)) (interface{}, error) {
	if weight <= 0 {
		return nil, errors.New("weight must be positive")
	}

	if !l.node.Reject(weight) {
		return f()
	} else {
		return nil, OverQuotaError
	}
}
