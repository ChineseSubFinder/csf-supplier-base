package orc_helper

import (
	"testing"
)

func TestORCHelper_File(t *testing.T) {

	orcHelper := NewORCHelper("http://127.0.0.1:18080")

	s, err := orcHelper.File("C:\\Users\\allan716\\SynologyDrive\\Tmp\\ca_target\\13531_37bd0f7275e6ca6587f9deae44b3ce01_realesrgan-x4plus.png")
	if err != nil {
		t.Fatal(err)
	}
	println("code:", s)
}

func TestORCHelper_GetStatus(t *testing.T) {

	orcHelper := NewORCHelper("http://192.168.50.135:18080")
	if orcHelper.GetStatus() == false {
		t.Fatal("orcHelper.GetStatus() == false")
	}
}
