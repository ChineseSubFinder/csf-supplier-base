package task_system

import (
	"github.com/ChineseSubFinder/csf-supplier-base/pkg"
	"testing"
)

func TestTaskSystemClient_AddMachineTranslationTask(t *testing.T) {

	subtitleAPIKEY := ""
	client := NewTaskSystemClient(subtitleAPIKEY)

	//getTranslatedResult, err := client.GetTranslatedResult("hkWySGaWCCfFcYDYRXkmbrfhrGeQaNwNVmFbmkurTMgEPvPuTIZhpwuORTXGDArc")
	//if err != nil {
	//	t.Errorf("GetTranslatedResult error: %v", err)
	//	return
	//}
	//println("getTranslatedResult:", getTranslatedResult.Status, getTranslatedResult.Message)

	imdbID := "tt2861424"
	isMovie := false
	Season := 1
	Episode := 1

	isAudioOrSRT := true
	fileFullPath := "D:\\tmp\\test_audio\\tt2861424\\output-full.wav"

	//isAudioOrSRT := false
	//fileFullPath := "D:\\tmp\\test_srt\\tt2861424\\123.srt"

	AudioSrcLanguage := "en"
	TranslatedLanguage := "CN"
	taskPackageResp, err := client.AddMachineTranslationTask(imdbID, isMovie, Season, Episode,
		isAudioOrSRT, fileFullPath, AudioSrcLanguage, TranslatedLanguage)
	if err != nil {
		t.Errorf("AddMachineTranslationTask error: %v", err)
		return
	}

	if taskPackageResp.Status != 1 {
		t.Errorf("AddMachineTranslationTask error: %v", taskPackageResp.Message)
		return
	}

	err = pkg.UploadFile2R2(taskPackageResp.UploadURL, fileFullPath)
	if err != nil {
		t.Errorf("UploadFile2R2 error: %v", err)
		return
	}

	token := taskPackageResp.Token
	taskPackageID := taskPackageResp.TaskPackageId

	firstPackageTaskDone, err := client.SetFirstPackageTaskDone(taskPackageID, token)
	if err != nil {
		t.Errorf("SetFirstPackageTaskDone error: %v", err)
		return
	}
	println("firstPackageTaskDone:", firstPackageTaskDone.Status, firstPackageTaskDone.Message)

	taskPackageStatus, err := client.GetTaskPackageStatus(taskPackageID)
	if err != nil {
		t.Errorf("GetTaskPackageStatus error: %v", err)
		return
	}
	println("taskPackageStatus:", taskPackageStatus.Status, taskPackageStatus.Message)

	//println("TaskPackage ID: ", taskPackageResp.TaskPackageId)
}
