package nlp

import (
	"strings"

	"github.com/ikawaha/kagome/tokenizer"
)

func TrimRepeatString(s string, w int) string {
	in := s
	for i := w; i > 1; i-- {
		ng := NewNgram(i)
		result := ng.Split(in)

		for {
			isDelete := false
			for i, e := range result.body {
				ae := result.GetAdjacentElements(i)
				if ae == nil {
					break
				}
				if e == *ae {
					result.DeleteAjacentElement(i)
					isDelete = true
					break
				}
			}
			if !isDelete {
				break
			}
		}

		if result.Build() != s {
			return TrimRepeatString(result.Build(), w-1)
		}
	}
	return in
}

func filterHeadCaseParticle(s string) string {
	tokens := t.Tokenize(s)
	builder := new(strings.Builder)

	filtered := false
	for i := 1; i < len(tokens); i++ {
		token := tokens[i]
		if token.Class == tokenizer.DUMMY {
			continue
		}
		if i == 1 {
			switch token.Pos() {
			case "格助詞", "助詞":
				filtered = true
				continue
			}
		}
		builder.WriteString(token.Surface)
	}
	if filtered {
		return filterHeadCaseParticle(builder.String())
	}
	return builder.String()
}

func TrimEmotionalVerb(s string) string {
	tokens := t.Tokenize(s)
	ret := ""

	for _, token := range tokens {
		if token.Class == tokenizer.DUMMY {
			continue
		}
		switch token.Pos() {
		case "感動詞":
			continue
		}
		switch token.Features()[0] {
		case "フィラー":
			continue
		}

		ret += token.Surface
	}

	return ret
}

func SplitClause(s string) []string {
	tokens := t.Tokenize(s)

	ret := []string{}

	buf := ""
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		if token.Class == tokenizer.DUMMY {
			continue
		}
		buf += token.Surface
		switch token.Pos() {
		case "助動詞":
			if i+1 < len(tokens) {
				if tokens[i+1].Pos() == "助詞" && tokens[i+1].Surface == "か" {
					buf += tokens[i+1].Surface
					i++
				} else if len(tokens[i+1].Features()) >= 2 && tokens[i+1].Features()[1] == "接続助詞" {
					//continue
					buf += tokens[i+1].Surface
					i++
				}

				for n := 1; n+i < len(tokens); n++ {
					if len(tokens[i+n].Features()) >= 2 && strings.Contains(tokens[i+n].Features()[1], "終助詞") {
						buf += tokens[i+n].Surface
					} else {
						i += (n - 1)
						break
					}
				}
			}

			if token.Surface == "です" {
				ret = append(ret, buf)
				buf = ""
				continue
			} else if token.Surface == "ます" {
				ret = append(ret, buf)
				buf = ""
				continue
			}
		}
	}

	if buf != "" {
		ret = append(ret, buf)
	}

	for i, v := range ret {
		ret[i] = filterHeadCaseParticle(v)
	}
	return ret
}
