package main

import "strings"

func GermanContains(a, b string) bool {
	r := strings.NewReplacer(
		"ä", "ae",
		"ö", "oe",
		"ü", "ue",
		"ß", "ss",
		"Ä", "Ae",
		"Ö", "Oe",
		"Ü", "Ue",
	)

	a = strings.TrimSpace(r.Replace(a))
	b = strings.TrimSpace(r.Replace(b))
	return strings.Contains(strings.ToLower(a), strings.ToLower(b))
}

func GermanEquals(a, b string) bool {
	r := strings.NewReplacer(
		"ä", "ae",
		"ö", "oe",
		"ü", "ue",
		"ß", "ss",
		"Ä", "Ae",
		"Ö", "Oe",
		"Ü", "Ue",
	)

	a = strings.TrimSpace(r.Replace(a))
	b = strings.TrimSpace(r.Replace(b))
	return strings.EqualFold(a, b)
}
