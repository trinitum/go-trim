// Package trim provides functions to create sets of characters and trim
// strings removing characters outside of the specified set
package trim

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	maxRune      = '\U0010FFFF'
	surrogateMin = 0xD800
	surrogateMax = 0xDFFF
)

// RuneSet represents a set of characters
type RuneSet unicode.RangeTable

// Includes checks if character is in the set
func (rs *RuneSet) Includes(r rune) bool {
	return unicode.Is((*unicode.RangeTable)(rs), r)
}

// NewRuneSetMust creates a new set just as NewRuneSet, but in case of error it
// panics instead of returning an error
func NewRuneSetMust(chars string) *RuneSet {
	rset, err := NewRuneSet(chars)
	if err != nil {
		panic(err)
	}
	return rset
}

type rsState int

const (
	stateNew rsState = iota
	stateHaveRune
	stateRange
)

// NewRuneSet creates a new set from the provided string. String may list
// individual characters or ranges specified as two characters separated with
// '-'. If you want to include '-' specify it as the first or the last
// character. For example "0-9aeiou-" will create a set including digits from 0
// to 9, a, e, i, o, u, and '-'.
func NewRuneSet(chars string) (*RuneSet, error) {
	rs := &RuneSet{}
	var state rsState
	var lr rune
	for _, r := range chars {
		if state == stateRange {
			if r < lr {
				return nil, fmt.Errorf("invalid range %s-%s in rune set %s", string(lr), string(r), chars)
			}
			if lr < surrogateMin && r > surrogateMax {
				return nil, fmt.Errorf("range %s-%s includes surrogates", string(lr), string(r))
			}
			if err := rs.addRange(lr, r); err != nil {
				return nil, err
			}
			state = stateNew
			continue
		}
		if state == stateHaveRune {
			if r == '-' {
				state = stateRange
				continue
			}
			if err := rs.addRange(lr, lr); err != nil {
				return nil, err
			}
		}
		lr = r
		state = stateHaveRune
	}
	if state != stateNew {
		if err := rs.addRange(lr, lr); err != nil {
			return nil, err
		}
	}
	if state == stateRange {
		if err := rs.addRange('-', '-'); err != nil {
			return nil, err
		}
	}
	return rs, nil
}

func (rs *RuneSet) addRange(lo, hi rune) error {
	if lo <= 0xffff {
		if hi > 0xffff {
			if err := rs.addRange16(uint16(lo), 0xffff); err != nil {
				return err
			}
			return rs.addRange32(0x10000, uint32(hi))
		}
		return rs.addRange16(uint16(lo), uint16(hi))
	}
	return rs.addRange32(uint32(lo), uint32(hi))
}

func (rs *RuneSet) addRange16(lo, hi uint16) error {
	var i int
	for _, r := range rs.R16 {
		if lo < r.Lo {
			if hi >= r.Lo {
				return fmt.Errorf("range %s-%s overlaps with %s-%s", string(lo), string(hi), string(r.Lo), string(r.Hi))
			}
			break
		}
		if lo <= r.Hi {
			return fmt.Errorf("range %s-%s overlaps with %s-%s", string(lo), string(hi), string(r.Lo), string(r.Hi))
		}
		i++
	}
	if hi <= unicode.MaxLatin1 {
		rs.LatinOffset++
	}
	var r16 []unicode.Range16
	if i > 0 {
		r16 = append(r16, rs.R16[0:i]...)
	}
	r16 = append(r16, unicode.Range16{Lo: lo, Hi: hi, Stride: 1})
	r16 = append(r16, rs.R16[i:]...)
	rs.R16 = r16
	return nil
}

func (rs *RuneSet) addRange32(lo, hi uint32) error {
	var i int
	for _, r := range rs.R32 {
		if lo < r.Lo {
			if hi >= r.Lo {
				return fmt.Errorf("range %s-%s overlaps with %s-%s", string(lo), string(hi), string(r.Lo), string(r.Hi))
			}
			break
		}
		if lo <= r.Hi {
			return fmt.Errorf("range %s-%s overlaps with %s-%s", string(lo), string(hi), string(r.Lo), string(r.Hi))
		}
	}
	var r32 []unicode.Range32
	if i > 0 {
		r32 = append(r32, rs.R32[0:i]...)
	}
	r32 = append(r32, unicode.Range32{Lo: lo, Hi: hi, Stride: 1})
	r32 = append(r32, rs.R32[i:]...)
	rs.R32 = r32
	return nil
}

// Trim returns a string that is a copy of s with all characters that are not
// in the set removed.
func (rs *RuneSet) Trim(s string) string {
	return strings.Map(func(r rune) rune {
		if rs.Includes(r) {
			return r
		}
		return -1
	}, s)
}
