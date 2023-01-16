package settings

type XrayPoolConfig struct {
	DefaultIndex        int        // 默认使用哪个代理池
	IntervalOfSwitching int64      // 多个 Xray Pool 时，切换的间隔时间，单位：秒
	Pools               []PoolInfo // Xray Pool 的信息
}

type PoolInfo struct {
	Url                string
	Port               string
	BrowserInstanceNum int
}
