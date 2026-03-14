package utilities

import (
	"strings"

)

func GetIP(id string) string {
    parts := strings.SplitN(id, "-", 2)
    if len(parts) < 2 {
        return ""
    }
    return parts[1]
}
