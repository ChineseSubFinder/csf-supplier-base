package settings

import "github.com/ChineseSubFinder/csf-supplier-base/pkg/common"

type ZiMuKuConfig struct {
	SearchFormatUrlPrefix  string                // 搜索的格式化url前缀
	SiteRootUrl            string                // https://zimuku.org
	SubsSaveRootDirPath    string                // 字幕存储的根目录，后续在数据库中只会存储相对路径
	Interval               int                   // 每隔多少分钟检查一次，这个是理想的间隔
	HotMovieEnable         bool                  // 下载热门电影的
	HotTVEnable            bool                  // 下载热门连续剧
	HotMovieLimit          int                   // 热门电影的数量限制
	HotTVLimit             int                   // 热门连续剧的数量限制
	OneProxyUseSearchCount int                   // 一个代理使用多少次搜索后，就换一个代理
	IntervalConfig         common.IntervalConfig // 每隔多少分钟检查一次
	TopMovieEnable         bool                  // 下载Top电影的
	TopTVEnable            bool                  // 下载Top连续剧
	FinishMovieEnable      bool                  // 下载已经上映电影的
	FinishTVEnable         bool                  // 下载已经完结连续剧
}
