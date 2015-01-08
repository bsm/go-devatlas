package devatlas

import (
	"log"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/martinolsen/go-pcre"
)

func rxCompile(expr string) (*Regexp, error) {
	// TODO - pcre.Study

	if re, err := pcre.Compile(expr, pcre.Utf8|pcre.Dupnames, nil); err != nil {
		return nil, err
	} else {
		return &Regexp{expr: expr, pcre: re}, nil
	}
}

func rxMustCompile(str string) *Regexp {
	if re, err := rxCompile(str); err != nil {
		panic(err)
	} else {
		return re
	}
}

type Regexp struct {
	expr string
	pcre *pcre.PCRE
}

func (re *Regexp) MatchString(s string) bool {
	if err := re.pcre.Exec(nil, s, 0, 0, nil); err == pcre.ErrNomatch {
		return false
	} else if err < 0 {
		panic("dont know what to do")
	}

	return true
}

func (re *Regexp) FindAllSubmatchIndex(b []byte, n int) [][]int {
	var locs [][]int
	var options pcre.Option
	for start := 0; start <= len(b) && n != 0; n-- {
		ovector := make([]int, (1+re.pcre.Capturecount())*3)
		if e := re.pcre.Exec(nil, string(b), start, options, ovector); e == pcre.ErrNomatch {
			break
		} else if e < 0 {
			log.Panicf("while matching %q[%d:]: %d", b, start, e)
		} else {
			locs = append(locs, ovector[:(1+re.pcre.Capturecount())*2])
		}
		options |= pcre.NotemptyAtstart
		start = ovector[1]
	}

	return locs
}

func (re *Regexp) FindSubmatch(b []byte) [][]byte {
	var bs [][]byte
	locs := re.FindSubmatchIndex(b)
	for i := 0; i < len(locs); i += 2 {
		if locs[i] == -1 || locs[i+1] == -1 {
			bs = append(bs, []byte{})
		} else {
			bs = append(bs, b[locs[i]:locs[i+1]])
		}
	}
	return bs
}

func (re *Regexp) FindStringSubmatch(s string) []string {
	var subs []string
	for _, sub := range re.FindSubmatch([]byte(s)) {
		subs = append(subs, string(sub))
	}
	return subs
}

func (re *Regexp) FindSubmatchIndex(b []byte) []int {
	var t = string(b) == "aacc" || re.expr == "(a){0}"
	ovector := make([]int, (1+re.pcre.Capturecount())*3)
	if e := re.pcre.Exec(nil, string(b), 0, 0, ovector); e == pcre.ErrNomatch {
		return nil
	} else if e < 0 {
		log.Panicf("e: %s", e)
	} else if t {
		log.Printf("expr: %q, b: %q, e: %d, ovector: %#v, %#v", re.expr, b, e, ovector, ovector[:(1+re.pcre.Capturecount())*2])
	}
	return ovector[:(1+re.pcre.Capturecount())*2]
}

func (re *Regexp) ReplaceAllString(src, repl string) string {
	return string(re.ReplaceAll([]byte(src), []byte(repl)))
}

func (re *Regexp) ReplaceAll(src, repl []byte) []byte {
	return re.replaceAll(src, func(dst []byte, match []int) []byte {
		return re.expand(dst, string(repl), src, match)
	})
}

func (re *Regexp) SubexpNames() []string { return re.pcre.Nametable() }

func (re *Regexp) expand(dst []byte, template string, src []byte, match []int) []byte {
	for len(template) > 0 {
		i := strings.Index(template, "$")
		if i < 0 {
			break
		}
		dst = append(dst, template[:i]...)
		template = template[i:]
		if len(template) > 1 && template[1] == '$' {
			// Treat $$ as $.
			dst = append(dst, '$')
			template = template[2:]
			continue
		}
		name, num, rest, ok := re.extract(template)
		if !ok {
			// Malformed; treat $ as raw text.
			dst = append(dst, '$')
			template = template[1:]
			continue
		}
		template = rest
		if num >= 0 {
			if 2*num+1 < len(match) && match[2*num] >= 0 {
				dst = append(dst, src[match[2*num]:match[2*num+1]]...)
			}
		} else {
			for i, namei := range re.SubexpNames() {
				if name == namei && 2*i+1 < len(match) && match[2*i] >= 0 {
					dst = append(dst, src[match[2*i]:match[2*i+1]]...)
					break
				}
			}
		}
	}
	return append(dst, template...)
}

func (re *Regexp) extract(str string) (name string, num int, rest string, ok bool) {
	if len(str) < 2 || str[0] != '$' {
		return
	}
	brace := false
	if str[1] == '{' {
		brace = true
		str = str[2:]
	} else {
		str = str[1:]
	}
	i := 0
	for i < len(str) {
		rune, size := utf8.DecodeRuneInString(str[i:])
		if !unicode.IsLetter(rune) && !unicode.IsDigit(rune) && rune != '_' {
			break
		}
		i += size
	}
	if i == 0 {
		// empty name is not okay
		return
	}
	name = str[:i]
	if brace {
		if i >= len(str) || str[i] != '}' {
			// missing closing brace
			return
		}
		i++
	}

	// Parse number.
	num = 0
	for i := 0; i < len(name); i++ {
		if name[i] < '0' || '9' < name[i] || num >= 1e8 {
			num = -1
			break
		}
		num = num*10 + int(name[i]) - '0'
	}
	// Disallow leading zeros.
	if name[0] == '0' && len(name) > 1 {
		num = -1
	}

	rest = str[i:]
	ok = true
	return
}

func (re *Regexp) replaceAll(src []byte, repl func([]byte, []int) []byte) []byte {
	var dst []byte
	locs := re.FindAllSubmatchIndex(src, -1)
	var srci int
	for i := 0; i < len(locs); i++ {
		dst = append(dst, src[srci:locs[i][0]]...)
		dst = repl(dst, locs[i])
		srci = locs[i][1]
	}
	if srci < len(src) {
		dst = append(dst, src[srci:]...)
	}
	return dst
}
