package middlewares

import (
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/cushydigit/microstore/shared/helpers"
	myredis "github.com/cushydigit/microstore/shared/redis"
)

const (
	Limit     = 5
	WindowSec = 60
)

func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)
		key := myredis.RateLimiterKey(ip)

		count, err := myredis.Client.Incr(myredis.Ctx, key).Result()
		if err != nil {
			helpers.ErrorJSON(w, errors.New("Resis error"), http.StatusInternalServerError)
			return
		}

		if count == 1 {
			myredis.Client.Expire(myredis.Ctx, key, time.Duration(WindowSec)*time.Second)
		}

		if count > Limit {
			helpers.ErrorJSON(w, errors.New("too many request"), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)

	})
}

// TODO feature add Nginx for trusted reverse proxy
// this method can not detect user that use proxy for hiding their ip addres
func getIP(r *http.Request) string {
	//
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // fallback
	}
	return ip
}
