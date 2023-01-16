package zimuku

import (
	"testing"
)

func TestGetImg(t *testing.T) {

	err := GetImg("https://zimuku.org/", 50, "C:\\WorkSpace\\TrainData\\zimuku\\Validate")
	if err != nil {
		t.Fatal(err)
	}
}
