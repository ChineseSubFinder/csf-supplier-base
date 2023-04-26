package subtitle_mark

// MarkType 定义一个枚举类型：人工翻译、机器翻译
type MarkType int

const (
	MarkTypeManual  MarkType = iota + 1 // 人工翻译
	MarkTypeMachine                     // 机器翻译
)

// Score 定义一个枚举类型：满意、一般、差、不匹配
type Score int

const (
	Good     Score = iota + 1 // 满意
	Normal                    // 一般
	Bad                       // 差
	NotMatch                  // 不匹配
)
