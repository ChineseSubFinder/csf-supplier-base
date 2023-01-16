package settings

import (
	"time"
)

type TimeConfig struct {
	OnePageTimeOut                 int   // 一页请求的超时时间，单位是秒
	OneProxyNodeUseInternalMinTime int32 // 一个代理节点，两次使用最短间隔，单位是秒
	OneProxyNodeUseInternalMaxTime int32 // 一个代理节点，两次使用最长间隔，单位是秒
	ProxyNodeSkipAccessTime        int64 // 设置一个代理节点可被再次访问的时间间隔（然后需要再加上现在时间为基准来算），单位是秒
}

// GetOnePageTimeOut 获取一页请求的超时时间，单位是秒
func (t *TimeConfig) GetOnePageTimeOut() time.Duration {
	return time.Duration(t.OnePageTimeOut) * time.Second
}
