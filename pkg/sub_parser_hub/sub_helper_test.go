package sub_parser_hub

import (
	"testing"
)

func TestSearchMatchedSubFileByDir(t *testing.T) {

	ssDir := "C:\\WorkSpace\\Go2Hell\\src\\github.com\\ChineseSubFinder\\ChineseSubFinder-TestData\\FixTimeline\\org"
	got, err := SearchMatchedSubFileByDir(ssDir)
	if err != nil {
		t.Errorf("SearchMatchedSubFileByDir() error = %v", err)
		return
	}
	println(got)
}
