package tactics

import "strings"

var registry = map[string]Func{
	"HFREE":  HFree,
	"HGET":   HGet,
	"HRAX":   HRax,
	"HLEGIT": HLegit,
	"HFULL":  HFull,
	"HTOR":   HTor,
	"HBROW":  HBrow,
}

// Get mengembalikan fungsi metode berdasarkan nama (case-insensitive).
func Get(name string) (Func, bool) {
	fn, ok := registry[strings.ToUpper(strings.TrimSpace(name))]
	return fn, ok
}

// Names mengembalikan daftar metode yang tersedia.
func Names() []string {
	out := make([]string, 0, len(registry))
	for k := range registry {
		out = append(out, k)
	}
	return out
}
