package catalog

import (
	"fmt"
	"net/url"
	"strings"
)

// BuildCommand maps query parameters to the method's command template.
func BuildCommand(cfg MethodConfig, target string, timeSec int, q url.Values) string {
	cmdParts := []string{cfg.Command, target, fmt.Sprintf("%d", timeSec)}

	for paramKey := range cfg.Params {
		if paramKey == "target" || paramKey == "time" || paramKey == "host" || paramKey == "duration" {
			continue
		}
		if val := q.Get(paramKey); val != "" {
			cmdParts = append(cmdParts, val)
		}
	}

	for k, vs := range q {
		if len(vs) == 0 || vs[0] == "" {
			continue
		}
		switch strings.ToLower(k) {
		case "key", "token", "target", "host", "time", "duration", "method", "concurrent":
			continue
		default:
			if _, exists := cfg.Params[k]; !exists {
				cmdParts = append(cmdParts, vs[0])
			}
		}
	}

	return strings.Join(cmdParts, " ")
}
