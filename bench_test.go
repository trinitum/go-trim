package trim

import (
	"regexp"
	"strings"
	"testing"
)

func BenchmarkTrims(b *testing.B) {
	testString := `foo\bar-baz.?你们你们`
	testTrimmed := `foobar-baz.`

	re := regexp.MustCompile(`[^a-zA-Z0-9_.-]`)
	chars := "abcdefghigklmnopqrstuvwxyzABCDEFGHIGKLMNOPQRSTUVWXYZ0123456789_.-"
	mapFunc := func(r rune) rune {
		if !strings.ContainsRune(chars, r) {
			return -1
		}
		return r
	}
	checkFunc := func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '.' || r == '-' {
			return r
		}
		return -1
	}
	rset := NewRuneSetMust("a-zA-Z0-9_.-")

	benchs := []struct {
		name string
		f    func(in string) string
	}{
		{
			name: "Regexp",
			f:    func(in string) string { return re.ReplaceAllString(in, "") },
		},
		{
			name: "strings.Map(with if)",
			f: func(in string) string {
				return strings.Map(checkFunc, in)
			},
		},
		{
			name: "strings.Map(IndexRune)",
			f: func(in string) string {
				return strings.Map(mapFunc, in)
			},
		},
		{
			name: "trim",
			f: func(in string) string {
				return rset.Trim(in)
			},
		},
	}
	for _, bench := range benchs {
		if res := bench.f(testString); res != testTrimmed {
			panic("benchmark function " + bench.name + " returned incorrect result: " + res)
		}
		b.Run(bench.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bench.f(testString)
			}
		})
	}
}
