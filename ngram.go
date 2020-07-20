package nlp

import (
	"fmt"
	"strings"

	"github.com/ikawaha/kagome/tokenizer"
)

var t = tokenizer.New()

func genCorpus(s string) Corpus {
	tokens := t.Tokenize(s)
	ret := Corpus{}

	for _, token := range tokens {
		if token.Class == tokenizer.DUMMY {
			continue
		}

		ret = append(ret, token.Surface)
	}

	return ret
}

func ngramFromStringModeCharacterUnit(s string, n int) *NgramResult {
	if len([]rune(s)) < n {
		return &NgramResult{body: []string{s}, n: n}
	}
	ret := &NgramResult{n: n}

	for i := 0; i+n <= len([]rune(s)); i++ {
		buf := []rune(s)[i : i+n]
		ret.body = append(ret.body, string(buf))
	}

	return ret
}

func ngramFromCorpusModeCorpusUnit(c Corpus, n int) *NgramResult {
	if len(c) < n {
		return &NgramResult{body: []string{strings.Join(c, "")}, n: n}
	}
	ret := &NgramResult{n: n}
	for i := 0; i+n <= len(c); i++ {
		ret.body = append(ret.body, strings.Join(c[i:i+n], ""))
	}

	return ret
}

type Mode int

const (
	ModeCharacterUnit Mode = iota
	ModeCorpusUnit
)

type Ngram struct {
	N    int
	Mode Mode
}

func (n *Ngram) Split(s string) *NgramResult {
	switch n.Mode {
	case ModeCharacterUnit:
		return ngramFromStringModeCharacterUnit(s, n.N)
	case ModeCorpusUnit:
		return ngramFromCorpusModeCorpusUnit(genCorpus(s), n.N)
	}
	return nil
}

func (n *Ngram) SplitFromCorpus(c Corpus) *NgramResult {
	switch n.Mode {
	case ModeCharacterUnit:
		return ngramFromStringModeCharacterUnit(c.String(), n.N)
	case ModeCorpusUnit:
		return ngramFromCorpusModeCorpusUnit(c, n.N)
	}
	return nil
}

func NewNgram(n int) *Ngram {
	return &Ngram{
		N: n,
	}
}

type Corpus []string

func (c Corpus) String() string {
	if len(c) == 0 {
		return ""
	}
	return strings.Join(c, "")
}

type NgramResult struct {
	n    int
	body Corpus
}

func (r *NgramResult) String() string {
	return fmt.Sprint([]string(r.body))
}

func (r *NgramResult) Corpus() Corpus {
	return r.body
}

func (r *NgramResult) GetAdjacentElements(i int) *string {
	if i < 0 {
		return nil
	}
	if len(r.body) <= i+r.n {
		return nil
	}
	return &r.body[i+r.n]
}

func (r *NgramResult) GetElements(i int) *string {
	if i < 0 {
		return nil
	}
	if len(r.body) <= i {
		return nil
	}
	return &r.body[i]
}

func (r *NgramResult) DeleteElement(i int) {
	if i > len(r.body) {
		return
	}
	r.body = append(r.body[:i], r.body[i+1:]...)
}

func (r *NgramResult) DeleteAjacentElement(i int) {
	if i < 0 {
		return
	}
	if len(r.body) <= i+r.n {
		return
	}

	for n := i + r.n; i+1 <= n; n-- {
		r.DeleteElement(n)
	}
}

func (r *NgramResult) Build() string {
	if len(r.body) == 0 {
		return ""
	}
	ret := r.body[0]
	if len(r.body) == 1 {
		return ret
	}

	for i := 1; i < len(r.body); i++ {
		tail := []rune(r.body[i])
		tail = tail[len(tail)-1:]
		ret += string(tail)
	}

	return ret
}
