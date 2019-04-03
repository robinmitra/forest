package disk

import "fmt"

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

func Humanise(bytes int64) string {
	if bytes < GB {
		if bytes < MB {
			if bytes < KB {
				return fmt.Sprintf("%i B", bytes)
			}
			return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
		}
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	} else {
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	}
}
