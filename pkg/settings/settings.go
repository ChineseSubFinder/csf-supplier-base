package settings

import (
	"github.com/WQGroup/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"sync"
)

func Get() *Config {

	_updateConfigLocker.Lock()
	defer _updateConfigLocker.Unlock()

	if _config == nil {
		_initViperInstanceOnce.Do(func() {
			// 读取Yaml配置文件，并转换成 struct 结构
			_viper = viper.New()
			_viper.AddConfigPath(".")                   //设置读取的文件路径
			_viper.SetConfigName("csf-supplier-config") //设置读取的文件名
			_viper.SetConfigType("yaml")                //设置文件的类型
			//尝试进行配置读取
			if err := _viper.ReadInConfig(); err != nil {
				logger.Panicln(err)
			}
			// 监听配置更改
			_viper.WatchConfig()
			// 反序列化配置
			err := _viper.Unmarshal(&_config)
			if err != nil {
				logger.Panicln(err)
			}
			_viper.OnConfigChange(func(e fsnotify.Event) {
				logger.Printf("config.yaml:%s Op:%s\n", e.Name, e.Op)
				// 反序列化配置
				err := _viper.Unmarshal(&_config)
				if err != nil {
					logger.Panicln(err)
				}
			})
		})
	}

	return _config
}

func Save(key string, value interface{}) {
	_updateConfigLocker.Lock()
	defer _updateConfigLocker.Unlock()

	if _viper == nil {
		logger.Warningln("Save config failed, _viper is nil")
		return
	}
	_viper.Set(key, value)
	err := _viper.WriteConfig()
	if err != nil {
		logger.Errorln("Save config failed, err:", err)
		return
	}
}

type Config struct {
	CacheRootDirPath     string
	XrayPoolConfig       XrayPoolConfig
	DBConnectConfig      DBConnectConfig
	SuccessWordsConfig   SuccessWordsConfig
	FailWordsConfig      FailWordsConfig
	TimeConfig           TimeConfig
	TMDBConfig           TMDBConfig
	SubsHuo720Config     SubsHuo720Config
	CaptchaConfig        CaptchaConfig
	ZiMuKuConfig         ZiMuKuConfig
	SearchZiMuKuConfig   SearchZiMuKuConfig
	SubHDConfig          SubHDConfig
	HouseKeepingConfig   HouseKeepingConfig
	ZipConfig            ZipConfig
	ImdbInfoCenterConfig ImdbInfoCenterConfig
	CloudFlareConfig     CloudFlareConfig
	AuthConfig           AuthConfig
	TitleFilterConfig    TitleFilterConfig
}

var (
	_config                *Config
	_viper                 *viper.Viper
	_initViperInstanceOnce sync.Once
	_updateConfigLocker    sync.Mutex
)
