package captcha_client

import (
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/settings"
	"testing"
)

func TestCaptchaClient_Get(t *testing.T) {

	ca := NewCaptchaClient(settings.Get().CaptchaConfig.ZiMuKu.Url)
	result, err := ca.GetFromImgFile("C:\\WorkSpace\\TrainData\\zimuku\\Validate_mark\\03055_1b3e6b439e2645f4bf02f59eb388906e.bmp")
	if err != nil {
		t.Fatal(err)
	}

	println("result:", result)
}
