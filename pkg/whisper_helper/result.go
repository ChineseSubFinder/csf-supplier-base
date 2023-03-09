package whisper_helper

import (
	"fmt"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg"
	"strings"
)

type Segment struct {
	Id               int     `json:"id"`
	Seek             int     `json:"seek"`
	Start            float64 `json:"start"`
	End              float64 `json:"end"`
	Text             string  `json:"text"`
	Tokens           []int   `json:"tokens"`
	Temperature      float64 `json:"temperature"`
	AvgLogprob       float64 `json:"avg_logprob"`
	CompressionRatio float64 `json:"compression_ratio"`
	NoSpeechProb     float64 `json:"no_speech_prob"`
}

// ToSRTContent 将 segment 转换为 srt 格式的内容
func (s Segment) ToSRTContent() string {
	// SRT 字幕的一句话的模板
	// 1
	// 00:00:00,000 --> 00:00:00,000
	// 你好
	startTime := pkg.SecondsToHMS(s.Start)
	endTime := pkg.SecondsToHMS(s.End)
	return fmt.Sprintf("%d\n%s --> %s\n%s\n", s.Id, startTime, endTime, strings.TrimLeft(s.Text, " "))
}

// WhisperJsonResult whisper transcribe 的结果
type WhisperJsonResult struct {
	Text     string    `json:"text"`
	Segments []Segment `json:"segments"`
	Language string    `json:"language"`
}

// SplitSegment 分割句子，将连续的 segment 分割为多个句子
func (w *WhisperJsonResult) SplitSegment() [][]int {

	// 将 Segment 进行分组，时间轴连续的 Segment 为一组
	var segmentsGroup [][]int
	var tmpGroup []int
	for i, segment := range w.Segments {

		// 有可能 segment 的 index 不对，所以需要重新设置一下
		w.Segments[i].Id = i + 1

		if i == 0 {
			tmpGroup = append(tmpGroup, i)
			continue
		}
		// 如果当前的 segment 的开始时间和上一个 segment 的结束时间相差在 0.5 秒以内，则认为是连续的。
		if segment.Start-w.Segments[i-1].End <= 0.5 {
			tmpGroup = append(tmpGroup, i)
		} else {
			segmentsGroup = append(segmentsGroup, tmpGroup)
			tmpGroup = []int{i}
		}
	}
	segmentsGroup = append(segmentsGroup, tmpGroup)

	return segmentsGroup
}

// SplitSegmentGroup 分割 Segment Group 为多个 Segment Group
func (w *WhisperJsonResult) SplitSegmentGroup(group []int, maxGroupLen int) [][]int {

	groupText := ""
	for _, segmentIndex := range group {
		groupText += w.Segments[segmentIndex].ToSRTContent()
	}
	// 如果统计出来的最大长度小于 5000，则不需要分割
	if len(groupText) <= maxGroupLen {
		return [][]int{group}
	}
	// 如果统计出来的最大长度大于 5000，则需要分割
	// 将总大小除以 5000，得到需要分割的组数，分割最小单位是一个 segment，不能拆分 segment
	groupCount := len(groupText) / maxGroupLen
	groupCount += 3
	// 如果不能整除，则需要多分割一组
	if len(groupText)%maxGroupLen != 0 {
		groupCount++
	}
	// 计算每组的大小
	groupSize := len(group) / groupCount
	// 如果不能整除，则需要多分割一组
	if len(group)%groupCount != 0 {
		groupSize++
	}
	// 分割
	var segmentGroups [][]int
	for i := 0; i < len(group); i += groupSize {
		end := i + groupSize
		if end > len(group) {
			end = len(group)
		}
		segmentGroups = append(segmentGroups, group[i:end])
	}
	return segmentGroups
}
