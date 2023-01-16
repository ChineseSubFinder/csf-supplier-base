package zimuku

import (
	"errors"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/settings"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/rod_helper"
)

func GetImg(baseUrl string, times int, saveImgRootPath string) error {

	var timeConfig rod_helper.TimeConfig
	timeConfig.OneProxyNodeUseInternalMinTime = 10
	timeConfig.OneProxyNodeUseInternalMaxTime = 15
	timeConfig.ProxyNodeSkipAccessTime = 1200
	poolOptions := rod_helper.NewPoolOptions(logger.GetLogger(), true, true, timeConfig)
	poolOptions.SetCacheRootDirPath(settings.Get().CacheRootDirPath)
	poolOptions.SetXrayPoolUrl(settings.Get().XrayPoolConfig.Pools[0].Url)
	poolOptions.SetXrayPoolPort(settings.Get().XrayPoolConfig.Pools[0].Port)

	pool := rod_helper.NewPool(poolOptions)
	if pool == nil {
		return errors.New("pool is nil, xray_pool not running")
	}
	defer func() {
		pool.Close()
	}()

	//getOneImg := func(pool *rod_helper.Pool) (string, string, error) {

	//	page, _, err := rod_helper.NewPageNavigate(pool.GetLBBrowser(), pool.LbHttpUrl, baseUrl, 15*time.Second)
	//	defer func() {
	//		if page != nil {
	//			_ = page.Close()
	//		}
	//	}()
	//	if err != nil {
	//		if errors.Is(err, context.DeadlineExceeded) == false {
	//			// 不是超时错误，那么就返回错误，跳过
	//			return "", "", err
	//		}
	//	}
	//
	//	element, err := page.Element("body > div > div:nth-child(4) > table > tbody > tr:nth-child(1) > td:nth-child(3) > img")
	//	if err != nil {
	//		return "", "", err
	//	}
	//	attribute, err := element.Attribute("src")
	//	if err != nil {
	//		println(page.String())
	//		return "", "", err
	//	}
	//
	//	imgAttList := strings.Split(*attribute, ";")
	//	if len(imgAttList) != 2 {
	//		return "", "", errors.New("imgAttList len != 2")
	//	}
	//
	//	imgTypeAtt := imgAttList[0]
	//	imgTypes := strings.Split(imgTypeAtt, "/")
	//	if len(imgTypes) != 2 {
	//		return "", "", errors.New("imgTypes len != 2")
	//	}
	//	imgType := imgTypes[1]
	//	base64ImgString := imgAttList[1]
	//	base64ImgString = strings.Replace(base64ImgString, "base64,", "", 1)
	//
	//	return imgType, base64ImgString, nil
	//}
	//
	//for i := 0; i < times; i++ {
	//
	//	println("times:", i)
	//	//time.Sleep(15 * time.Second)
	//	imgType, base64ImgString, err := getOneImg(pool)
	//	if err != nil {
	//		println("getOneImg err:", err)
	//		continue
	//	}
	//
	//	_, err = pkg.Base642IMGFile(imgType, base64ImgString, saveImgRootPath)
	//	if err != nil {
	//		println("Base642IMGFile err:", err)
	//		continue
	//	}
	//}

	return nil
}
