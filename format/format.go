package format

import (
	"fmt"
)

const (
	B  = 1
	KB = B * 1024
	MB = KB * 1024
	GB = MB * 1024
	TB = GB * 1024
)

var (
	Sizes []int64  = []int64{B, KB, MB, GB, TB}
	Units []string = []string{"B", "kB", "MB", "GB", "TB"}
)

func Rate(r int64) string {
	return fmt.Sprintf("%d kB/s", r/1000)
}

func FileSize(sz int64) string {
	for i := len(Sizes) - 1; i >= 0; i-- {
		n := sz / Sizes[i]
		if n != 0 {
			s := fmt.Sprintf("%d", n)
			if n < 10 && i > 0 {
				f := (sz - n*Sizes[i]) / Sizes[i-1]
				if f > 0 {
					d := float64(f) * float64(Sizes[i-1]) / float64(Sizes[i])
					fs := fmt.Sprintf("%.3f", d)[2:3]
					if fs != "0" {
						s += "." + fs
					}
				}
			} else if n >= 1000 {
				s = "1"
				i += 1
			}
			s += " " + Units[i]

			return s
		}
	}

	return "0 " + Units[0]
}
