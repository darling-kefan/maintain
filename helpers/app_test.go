package helpers_test

import (
	"sort"
	"testing"

	"maintain/helpers"
)

func TestCheckTag(t *testing.T) {
	t.Logf("%#v", helpers.CheckTag("v1.2.32"))
	t.Logf("%#v", helpers.CheckTag("v0.1.3-457-ga8f22b9"))
}

func TestSortTags(t *testing.T) {
	// tags := []string{"v0.1.32", "v0.1.22", "v0.1.5", "v1.2.0", "v3.4.21", "v1.2.6", "v3.4"}
	tags := []string{"onlinefix0.1","onlinefix0.2","onlinefix0.3","onlinefix0.4","onlinefix0.5","onlinefix0.6","onlinefix0.7","onlinefix0.8","v0.1","v0.1.1","v0.1.10","v0.1.11","v0.1.12","v0.1.13","v0.1.14","v0.1.15","v0.1.16","v0.1.17","v0.1.18","v0.1.19","v0.1.2","v0.1.20","v0.1.21","v0.1.22","v0.1.23"}

	t.Logf("%v", tags)

	sort.Sort(helpers.Tags(tags))

	t.Logf("%v", tags)
}
