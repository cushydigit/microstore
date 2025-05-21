package redis

func RateLimiterKey(ip string) string {
	return "rl:" + ip
}
