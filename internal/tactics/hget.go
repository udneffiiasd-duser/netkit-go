package tactics

import "context"

// HGet — classic flooder untuk target umum (full power).
func HGet(ctx context.Context, p Params) (int64, error) {
	cfg := floodConfig{
		httpMethod: "GET",
		keepAlive:  true,
	}
	return runFlood(ctx, p, cfg)
}
