package locale

import "golang.org/x/text/message"

func LocaliseInt(num int64) string {
	p := message.NewPrinter(message.MatchLanguage("en"))
	return p.Sprintf("%d", num)
}

func LocaliseFloat(num float64) string {
	p := message.NewPrinter(message.MatchLanguage("en"))
	res := p.Sprintf("%.2f", num)
	return res
}
