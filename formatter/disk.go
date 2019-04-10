package formatter

import (
	"fmt"
	"github.com/robinmitra/forest/locale"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

func HumaniseStorage(bytes int64) string {
	if bytes < GB {
		if bytes < MB {
			if bytes < KB {
				return fmt.Sprintf("%s B", locale.LocaliseInt(bytes))
			}
			return fmt.Sprintf("%s KB", locale.LocaliseFloat(float64(bytes)/float64(KB)))
		}
		return fmt.Sprintf("%s MB", locale.LocaliseFloat(float64(bytes)/float64(MB)))
	}
	return fmt.Sprintf("%s GB", locale.LocaliseFloat(float64(bytes)/float64(GB)))
}
