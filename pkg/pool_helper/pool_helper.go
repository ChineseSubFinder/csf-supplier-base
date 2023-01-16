package pool_helper

import (
	"github.com/ChineseSubFinder/csf-supplier-base/pkg"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/settings"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/struct_json"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/rod_helper"
	"sync"
	"time"
)

type Pools struct {
	index  int
	pools  []*rod_helper.Pool
	locker sync.Mutex
}

func NewPool() *Pools {

	pool := &Pools{
		index: 0,
		pools: make([]*rod_helper.Pool, 0),
	}
	for i := range settings.Get().XrayPoolConfig.Pools {
		pool.pools = append(pool.pools, rod_helper.NewPool(GetPoolOption(i)))
	}
	if len(pool.pools) == 0 {
		logger.Panicln("NewPool:", len(pool.pools))
	}

	if settings.Get().XrayPoolConfig.DefaultIndex >= len(pool.pools) {
		logger.Panicln("NewPool:", settings.Get().XrayPoolConfig.DefaultIndex)
	}

	return pool
}

func (p *Pools) Close() {
	for _, pool := range p.pools {
		pool.Close()
	}
}

func (p *Pools) GetOnePool() *rod_helper.Pool {

	p.locker.Lock()
	defer p.locker.Unlock()

	switchPoolInfo := NewSwitchPoolInfo(time.Now().Unix(), settings.Get().XrayPoolConfig.DefaultIndex)
	if pkg.IsFile(currentSwitchPoolName) == true {
		err := struct_json.ToStruct(currentSwitchPoolName, &switchPoolInfo)
		if err != nil {
			logger.Panicln("GetOnePool.ToStruct:", err)
		}
	} else {
		// 不存在缓存文件，就新建一个
		err := struct_json.ToFile(currentSwitchPoolName, switchPoolInfo)
		if err != nil {
			logger.Panicln("GetOnePool.ToFile:", err)
		}
	}
	// 判断缓存的信息是否过期，需要更新了
	if time.Now().Unix()-switchPoolInfo.StartTime > settings.Get().XrayPoolConfig.IntervalOfSwitching {

		switchPoolInfo.StartTime = time.Now().Unix()
		switchPoolInfo.Index += 1
		if switchPoolInfo.Index >= len(p.pools) {
			switchPoolInfo.Index = 0
		}
		// 写入缓存文件
		err := struct_json.ToFile(currentSwitchPoolName, switchPoolInfo)
		if err != nil {
			logger.Panicln("GetOnePool.ToFile:", err)
		}
	}

	return p.pools[switchPoolInfo.Index]
}

func GetPoolOption(poolIndex int) *rod_helper.PoolOptions {
	var timeConfig rod_helper.TimeConfig
	timeConfig.OneProxyNodeUseInternalMinTime = settings.Get().TimeConfig.OneProxyNodeUseInternalMinTime
	timeConfig.OneProxyNodeUseInternalMaxTime = settings.Get().TimeConfig.OneProxyNodeUseInternalMaxTime
	timeConfig.ProxyNodeSkipAccessTime = settings.Get().TimeConfig.ProxyNodeSkipAccessTime

	poolOptions := rod_helper.NewPoolOptions(
		logger.GetLogger(),
		true,
		false,
		timeConfig,
	)
	poolOptions.SetCacheRootDirPath(settings.Get().CacheRootDirPath)
	poolOptions.SetXrayPoolUrl(settings.Get().XrayPoolConfig.Pools[poolIndex].Url)
	poolOptions.SetXrayPoolPort(settings.Get().XrayPoolConfig.Pools[poolIndex].Port)
	poolOptions.SetSuccessWordsConfig(rod_helper.SuccessWordsConfig{
		WordsConfig: rod_helper.WordsConfig{
			Enable: settings.Get().SuccessWordsConfig.Enable,
			Words:  settings.Get().SuccessWordsConfig.Words,
		},
	})
	poolOptions.SetFailWordsConfig(rod_helper.FailWordsConfig{
		WordsConfig: rod_helper.WordsConfig{
			Enable: settings.Get().FailWordsConfig.Enable,
			Words:  settings.Get().FailWordsConfig.Words,
		},
	})

	return poolOptions
}

type SwitchPoolInfo struct {
	StartTime int64 `json:"start_time"` // 启用这个 Xray Pools 的时间
	Index     int   `json:"index"`      // Xray Pools 的索引
}

func NewSwitchPoolInfo(startTime int64, index int) *SwitchPoolInfo {
	return &SwitchPoolInfo{StartTime: startTime, Index: index}
}

const currentSwitchPoolName = "current_pool.json"
