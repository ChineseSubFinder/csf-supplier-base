package whisper_helper

import (
	"fmt"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg"
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
	return fmt.Sprintf("%d\n%s --> %s\n%s\n", s.Id, startTime, endTime, s.Text)
}

// WhisperJsonResult whisper transcribe 的结果
type WhisperJsonResult struct {
	Text     string    `json:"text"`
	Segments []Segment `json:"segments"`
	Language string    `json:"language"`
}
