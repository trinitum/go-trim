package trim

import (
	"fmt"
	"unicode"
)

const (
	maxRune      = '\U0010FFFF'
	surrogateMin = 0xD800
	surrogateMax = 0xDFFF
)

type RuneSet unicode.RangeTable

func (rs *RuneSet) Includes(r rune) bool {
	return unicode.Is((*unicode.RangeTable)(rs), r)
}

func MakeRuneSet(chars string) (*RuneSet, error) {
	rs := &RuneSet{}
	var state int
	var lr rune
	for _, r := range chars {
		if state == 2 {
			if r < lr {
				return nil, fmt.Errorf("invalid range %s-%s in rune set %s", string(lr), string(r), chars)
			}
			if lr < surrogateMin && r > surrogateMax {
				return nil, fmt.Errorf("range %s-%s includes surrogates", string(lr), string(r))
			}
			if err := rs.addRange(lr, r); err != nil {
				return nil, err
			}
			state = 0
			continue
		}
		if state == 1 {
			if r == '-' {
				state = 2
				continue
			}
			if err := rs.addRange(lr, lr); err != nil {
				return nil, err
			}
			state = 0
		}
		lr = r
		state = 1
	}
	if state > 0 {
		if err := rs.addRange(lr, lr); err != nil {
			return nil, err
		}
	}
	if state == 2 {
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
