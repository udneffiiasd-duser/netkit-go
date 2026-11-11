package tactics

import "context"

// Params adalah parameter serangan yang dikirim API ke node.
type Params struct {
	Target     string            `json:"target"`
	Duration   int               `json:"time"`
	Threads    int               `json:"threads"`
	Rate       int               `json:"rate"`
	Concurrent int               `json:"concurrent"`
	Method     string            `json:"http_method"`
	Geo        string            `json:"geo"`
	Ver        string            `json:"ver"`
	UA         string            `json:"ua"`
	Cookie     string            `json:"cookie"`
	Referer    string            `json:"ref"`
	PostData   string            `json:"postdata"`
	Header     string            `json:"header"`
	Extra      map[string]string `json:"extra"`
}

// Func adalah signature setiap metode serangan.
type Func func(ctx context.Context, p Params) (int64, error)

func withDefaults(p Params) Params {
	if p.Threads <= 0 {
		p.Threads = 100
	}
	if p.Duration <= 0 {
		p.Duration = 60
	}
	if p.Method == "" {
		p.Method = "GET"
	}
	return p
}
