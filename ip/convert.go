package ip

import (
	"net/url"
	"sort"
	"strings"
)

// Map2String преобразует строку в хеш.
// Пример: "a=foo&b=bar" -> {"a":"foo", "b":"bar"}
func Meta2Map(meta string) map[string]string {
	var rt = make(map[string]string)
	if meta == "" {
		return rt
	}

	v, err := url.ParseQuery(meta)
	if err != nil {
		return rt
	}

	for key := range v {
		rt[key] = v.Get(key)
	}
	return rt
}

// Map2String преобразует отсортированный по ключу хеш 'meta' в строку.
// Пример: {"a":"foo", "b":"bar"} -> "a=foo&b=bar"
func Map2String(meta map[string]string) string {
	var buf strings.Builder

	keys := make([]string, 0, len(meta))
	for k := range meta {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		vs := meta[k]
		keyEscaped := url.QueryEscape(k)
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(keyEscaped)
		buf.WriteByte('=')
		buf.WriteString(url.QueryEscape(vs))
	}

	return buf.String()
}
