package idempotency

import (
	"errors"
	"github.com/negarciacamilo/deuna_challenge/application/logger"
	c "github.com/patrickmn/go-cache"
	"time"
)

/*
This is needed since the client can send multiple requests if anything fails. So to avoid this, in a productive application I WOULDN'T RECOMMEND USING A GLOBAL VARIABLE
Please read the docs to check how to do this in the -what I consider- the right way to it
*/

var cache *c.Cache

// Since this is just an example, I wouldn't normally use init as well
func init() {
	if cache == nil {
		// I believe that both Paypal and Stripe IK has a TTL like 24 hours or so
		cache = c.New(1*time.Minute, 1*time.Minute)
	}
}

func IdempotencyKeyExists(key string) bool {
	if cache == nil {
		msg := "cache is nil"
		logger.Panic(msg, "check-idempotency-key", errors.New(msg), nil)
	}

	_, found := cache.Get(key)
	if !found {
		// Store the key
		err := cache.Add(key, true, 1*time.Minute)
		// Best effort
		if err != nil {
			logger.Error("something happened storing the idempotency key", "check-idempotency-key", err, nil)
		}
		return false
	}

	return true
}
