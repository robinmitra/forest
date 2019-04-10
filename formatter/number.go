package formatter

import "github.com/robinmitra/forest/locale"

func HumaniseNumber(num int64) string {
	return locale.LocaliseInt(num)
}
