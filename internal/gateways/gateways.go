package gateways

import (
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
)

type gateways struct {
	tokenCache cmap.ConcurrentMap[string, TokenInfo]
}

type TokenInfo struct {
	Token      string
	Expiration time.Time
}
