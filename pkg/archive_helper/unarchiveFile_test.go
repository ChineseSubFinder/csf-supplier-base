package archive_helper

import (
	"github.com/ChineseSubFinder/csf-supplier-base/pkg"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/unit_test_helper"
	"path/filepath"
	"strings"
	"testing"
)

func TestUnArchiveFile(t *testing.T) {

	testRootDir := unit_test_helper.GetTestDataResourceRootPath([]string{"zips"}, 4, true)

	testUnArchive(t, testRootDir, "zip.zip")
	testUnArchive(t, testRootDir, "tar.tar")
	testUnArchive(t, testRootDir, "rar.rar")
	testUnArchive(t, testRootDir, "7z.7z")
}

func testUnArchive(t *testing.T, testRootDir string, missionName string) {
	fileFPath := filepath.Join(testRootDir, missionName)
	desPath := filepath.Join(testRootDir, strings.ReplaceAll(filepath.Ext(missionName), ".", ""))
	err := UnArchiveFile(fileFPath, desPath)
	if err != nil {
		t.Fatal(err)
	}
	if pkg.IsFile(filepath.Join(desPath, subASS)) == false {
		t.Fatal(missionName, " unArchive failed")
	}
	if pkg.IsFile(filepath.Join(desPath, subSRT)) == false {
		t.Fatal(missionName, " unArchive failed")
	}
}

const subASS = "oslo.2021.1080p.web.h264-naisu.繁体&英文.ass"
const subSRT = "oslo.2021.1080p.web.h264-naisu.繁体&英文.srt"

func TestUnArchiveFileEx(t *testing.T) {

	testRootDir := unit_test_helper.GetTestDataResourceRootPath([]string{"zips"}, 1, true)

	testRootDirAbs, err := filepath.Abs(testRootDir)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		fileFullPath string
		desRootPath  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[zimuku]_0_Inside No 9_S7E1.zip", args: args{
				fileFullPath: filepath.Join(testRootDirAbs, "[zimuku]_0_Inside No 9_S7E1.zip"),
				desRootPath:  filepath.Join(testRootDirAbs, "[zimuku]_0_Inside No 9_S7E1"),
			}, wantErr: false,
		},
		{
			name: "[zmk.pw]Suits.S04E01-E16.1080p.WEB-DL.DD5.1.H.264-NTb.rar", args: args{
				fileFullPath: filepath.Join(testRootDirAbs, "[zmk.pw]Suits.S04E01-E16.1080p.WEB-DL.DD5.1.H.264-NTb.rar"),
				desRootPath:  filepath.Join(testRootDirAbs, "[zmk.pw]Suits.S04E01-E16.1080p.WEB-DL.DD5.1.H.264-NTb"),
			}, wantErr: false,
		},
		{
			name: "[zmk.pw]The.Walking.Dead.S07E10.720p.HDTV.x264-AVS.rar", args: args{
				fileFullPath: filepath.Join(testRootDirAbs, "[zmk.pw]The.Walking.Dead.S07E10.720p.HDTV.x264-AVS.rar"),
				desRootPath:  filepath.Join(testRootDirAbs, "[zmk.pw]The.Walking.Dead.S07E10.720p.HDTV.x264-AVS"),
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UnArchiveFileEx(tt.args.fileFullPath, tt.args.desRootPath, true, true); (err != nil) != tt.wantErr {
				t.Errorf("UnArchiveFileEx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
