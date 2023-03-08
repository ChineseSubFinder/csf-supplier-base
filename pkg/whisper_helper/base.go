package whisper_helper

type SendTask struct {
	TaskID     int    `json:"task_id"`
	InputAudio string `json:"input_audio"`
	Language   string `json:"language"`
	TaskStatus int    `json:"task_status"` // 0 "pending", 1 "running", 2 "finished", 3 "error"
}

type SendTaskReply struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Status int    `json:"status"`
}

func (s SendTaskReply) GetTaskStatus() TaskStatus {
	return TaskStatus(s.Status)
}

type WhisperResult struct {
	Text     string `json:"text"`
	Segments []struct {
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
	} `json:"segments"`
	Language string `json:"language"`
}

// TaskStatus 任务状态的枚举类
type TaskStatus int

const (
	TaskStatusPending TaskStatus = iota + 1
	TaskStatusRunning
	TaskStatusFinished
	TaskStatusError
)

func (t TaskStatus) String() string {
	switch t {
	case TaskStatusPending:
		return "pending"
	case TaskStatusRunning:
		return "running"
	case TaskStatusFinished:
		return "finished"
	case TaskStatusError:
		return "error"
	default:
		return "unknown"
	}
}

/*
设置 whisper transcribe 的目标语言
LANGUAGES = {
	"en": "english",
	"zh": "chinese",
	"de": "german",
	"es": "spanish",
	"ru": "russian",
	"ko": "korean",
	"fr": "french",
	"ja": "japanese",
	"pt": "portuguese",
	"tr": "turkish",
	"pl": "polish",
	"ca": "catalan",
	"nl": "dutch",
	"ar": "arabic",
	"sv": "swedish",
	"it": "italian",
	"id": "indonesian",
	"hi": "hindi",
	"fi": "finnish",
	"vi": "vietnamese",
	"he": "hebrew",
	"uk": "ukrainian",
	"el": "greek",
	"ms": "malay",
	"cs": "czech",
	"ro": "romanian",
	"da": "danish",
	"hu": "hungarian",
	"ta": "tamil",
	"no": "norwegian",
	"th": "thai",
	"ur": "urdu",
	"hr": "croatian",
	"bg": "bulgarian",
	"lt": "lithuanian",
	"la": "latin",
	"mi": "maori",
	"ml": "malayalam",
	"cy": "welsh",
	"sk": "slovak",
	"te": "telugu",
	"fa": "persian",
	"lv": "latvian",
	"bn": "bengali",
	"sr": "serbian",
	"az": "azerbaijani",
	"sl": "slovenian",
	"kn": "kannada",
	"et": "estonian",
	"mk": "macedonian",
	"br": "breton",
	"eu": "basque",
	"is": "icelandic",
	"hy": "armenian",
	"ne": "nepali",
	"mn": "mongolian",
	"bs": "bosnian",
	"kk": "kazakh",
	"sq": "albanian",
	"sw": "swahili",
	"gl": "galician",
	"mr": "marathi",
	"pa": "punjabi",
	"si": "sinhala",
	"km": "khmer",
	"sn": "shona",
	"yo": "yoruba",
	"so": "somali",
	"af": "afrikaans",
	"oc": "occitan",
	"ka": "georgian",
	"be": "belarusian",
	"tg": "tajik",
	"sd": "sindhi",
	"gu": "gujarati",
	"am": "amharic",
	"yi": "yiddish",
	"lo": "lao",
	"uz": "uzbek",
	"fo": "faroese",
	"ht": "haitian creole",
	"ps": "pashto",
	"tk": "turkmen",
	"nn": "nynorsk",
	"mt": "maltese",
	"sa": "sanskrit",
	"lb": "luxembourgish",
	"my": "myanmar",
	"bo": "tibetan",
	"tl": "tagalog",
	"mg": "malagasy",
	"as": "assamese",
	"tt": "tatar",
	"haw": "hawaiian",
	"ln": "lingala",
	"ha": "hausa",
	"ba": "bashkir",
	"jw": "javanese",
	"su": "sundanese",
}
*/
