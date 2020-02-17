package render

import (
	"encoding/csv"
	"strconv"
	"strings"

	"github.com/zladovan/luster/fb"
)

// Csv renders collection of Fans as comma separated values with header
func Csv(fans fb.Fans) string {
	sb := &strings.Builder{}
	w := csv.NewWriter(sb)

	// header
	w.Write([]string{"TIME", "KIND", "ID", "NAME", "LINK"})

	// rows
	for _, f := range fans {
		w.Write([]string{
			strconv.Itoa(int(f.Time)),
			f.Kind.String(),
			f.Profile.ID,
			f.Profile.Name,
			f.Profile.Link(),
		})
	}

	w.Flush()

	return sb.String()
}
