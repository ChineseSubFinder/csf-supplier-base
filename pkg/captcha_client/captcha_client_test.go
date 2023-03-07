package captcha_client

import (
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/settings"
	"github.com/allanpk716/rod_helper"
	"testing"
)

func TestCaptchaClient_Get(t *testing.T) {

	rod_helper.InitFakeUA(settings.Get().CacheRootDirPath, settings.Get().TMDBConfig.TMDBHttpProxy)
	ca := NewCaptchaClient(settings.Get().CaptchaConfig.ZiMuKu.Url)
	dd := ""
	result, err := ca.GetFromImgBase64(dd)
	//result, err := ca.GetFromImgFile("123.bmp")
	if err != nil {
		t.Fatal(err)
	}

	println("result:", result)
}
