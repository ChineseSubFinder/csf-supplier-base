package settings

type SubsHuo720Config struct {
	StartIndexMovie int
	EndIndexMovie   int
	StartIndexTV    int
	EndIndexTV      int
	Thread          int
	StartFromBegin  bool // 是否从头开始，不管记录的进度
}
