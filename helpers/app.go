package helpers

import (
	"regexp"
	"strings"
	"strconv"
)

func CheckTag(tag string) bool {
	var validTag = regexp.MustCompile(`^(?:v|V)\d(\.[0-9a-zA-Z\-]+)+$`)
	if validTag.MatchString(tag) && !strings.Contains(strings.ToLower(tag), "fatal") && !strings.Contains(strings.ToLower(tag), "error") {
		return true
	}
	return false
}

type Tags []string
func (t Tags) Len() int { return len(t) }
func (t Tags) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t Tags) Less(i, j int) bool {
	ti, tj := strings.Split(t[i][1:], "."), strings.Split(t[j][1:], ".")
	var length int
	if len(ti) > len(tj) {
		length = len(tj)
	} else {
		length = len(ti)
	}
	for k := 0; k < length; k++ {
		m, _ := strconv.Atoi(ti[k])
		n, _ := strconv.Atoi(tj[k])
		if m == n {
			continue
		}
		return m > n
	}
	// 如果ti,tj的前半部分都一样，则认为长度长的tag大
	return len(ti) > len(tj)
}
