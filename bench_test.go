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
		if strings.IndexRune(chars, r) < 0 {
			return -1
		}
		return r
	}
	rset, err := NewRuneSet("a-zA-Z0-9_.-")
	if err != nil {
		panic(err)
	}

	benchs := []struct {
		name string
		f    func(in string) string
	}{
		{
			name: "Regexp",
			f:    func(in string) string { return re.ReplaceAllString(in, "") },
		},
		{
			name: "Check runes",
			f: func(in string) string {
				res := make([]rune, 0)
				for _, r := range in {
					if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_' || r == '.' || r == '-' {
						res = append(res, r)
					}
				}
				return string(res)
			},
		},
		{
			name: "strings.Map",
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
