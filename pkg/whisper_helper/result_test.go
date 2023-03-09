package whisper_helper

import (
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/struct_json"
	"testing"
)

func TestWhisperJsonResult_SplitSegment(t *testing.T) {

	var result WhisperJsonResult
	err := struct_json.ToStruct("No Responders Left Behind (2021) WEBRip-1080p.json", &result)
	if err != nil {
		t.Fatal(err)
	}

	max := 5000
	maxGroupLen := 0
	segmentsGroup := result.SplitSegment()
	for index, group := range segmentsGroup {

		groupText := ""
		for _, segmentIndex := range group {
			groupText += result.Segments[segmentIndex].ToSRTContent()
		}
		println("==================================")
		nowTextLen := len(groupText)

		if nowTextLen > maxGroupLen {
			maxGroupLen = nowTextLen
		}
		if nowTextLen >= max {
			println(index, "groupText length:", len(groupText))
			newGroup := result.SplitSegmentGroup(group, max)
			for _, newGroupItem := range newGroup {
				newGroupText := ""
				for _, segmentIndex := range newGroupItem {
					newGroupText += result.Segments[segmentIndex].ToSRTContent()
				}
				if len(newGroupText) > max {
					println("newGroupText length:", len(newGroupText))
				}
				println(newGroupText)
			}
		} else {
			println(groupText)
			println("==================================")
		}
	}

	println("maxGroupLen:", maxGroupLen)
}
