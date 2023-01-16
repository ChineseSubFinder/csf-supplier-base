package current_progress

import (
	"github.com/ChineseSubFinder/csf-supplier/pkg"
	"github.com/ChineseSubFinder/csf-supplier/pkg/struct_json"
	"sync"
)

type CurrentProgress struct {
	StartIndexMovie int `json:"start_index_movie"`
	StartIndexTV    int `json:"start_index_tv"`
}

func NewCurrentProgress() CurrentProgress {
	var cp CurrentProgress
	cp.StartIndexTV = 1
	cp.StartIndexMovie = 1
	return cp
}

func GetCurrentProgress(startFromBegin bool) (CurrentProgress, error) {

	if startFromBegin == true {
		return CurrentProgress{
			StartIndexMovie: 1,
			StartIndexTV:    1,
		}, nil
	}

	cp := NewCurrentProgress()
	if pkg.IsFile(currentProgressName) == true {
		err := struct_json.ToStruct(currentProgressName, &cp)
		if err != nil {
			return NewCurrentProgress(), err
		}
	}
	return cp, nil
}

func UpdateStartIndexMovie(nowIndex int) error {

	currentProgressLock.Lock()
	defer currentProgressLock.Unlock()

	cp, err := GetCurrentProgress(false)
	if err != nil {
		return err
	}
	cp.StartIndexMovie = nowIndex

	return struct_json.ToFile(currentProgressName, cp)
}

func UpdateStartIndexTV(nowIndex int) error {

	currentProgressLock.Lock()
	defer currentProgressLock.Unlock()

	cp, err := GetCurrentProgress(false)
	if err != nil {
		return err
	}
	cp.StartIndexTV = nowIndex

	return struct_json.ToFile(currentProgressName, cp)
}

const currentProgressName = "current_progress.json"

var currentProgressLock sync.Mutex
