package archive_helper

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"errors"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/cmdutil"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/settings"
	"github.com/WQGroup/logger"
	"github.com/bodgit/sevenzip"
	"github.com/mholt/archiver/v3"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

func CheckArchiveByBandizip(fileFullPath string) bool {
	var command *exec.Cmd
	cmdArgs := []string{
		"t",
		fileFullPath,
	}
	if pkg.IsFile(settings.Get().ZipConfig.ExeFPath) == false {
		logger.Panicln("Bandizip.exe not found")
	}

	command = exec.Command(settings.Get().ZipConfig.ExeFPath, cmdArgs...)
	stdoutBytes, err := cmdutil.ExecAndGetStdoutBytes(command)
	if err != nil {
		logger.Errorln("cmdutil.ExecAndGetStdoutBytes Error", err)
		return false
	}
	logger.Infoln("Bandizip stdout: ", string(stdoutBytes))
	if strings.Contains(string(stdoutBytes), "All OK") == true {
		return true
	}
	return false
}

func UnArchiveByBandizip(fileFullPath, desRootPath string) error {

	var command *exec.Cmd
	cmdArgs := []string{
		"x",
		//"-r",
		"-y",
		"-o:" + desRootPath,
		fileFullPath,
	}
	if pkg.IsFile(settings.Get().ZipConfig.ExeFPath) == false {
		logger.Panicln("Bandizip.exe not found")
	}

	command = exec.Command(settings.Get().ZipConfig.ExeFPath, cmdArgs...)
	stdoutBytes, err := cmdutil.ExecAndGetStdoutBytes(command)
	if err != nil {
		return err
	}
	logger.Infoln("Bandizip stdout: ", string(stdoutBytes))
	return nil
}

// UnArchiveFileEx 发现打包的字幕内部还有一层压缩包···所以···
func UnArchiveFileEx(fileFullPath, desRootPath string, useBandizip bool, ignoreUnZipErr bool) error {

	if pkg.IsDir(desRootPath) == false {
		err := os.MkdirAll(desRootPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	var err error
	// 先进行一次解压
	if useBandizip == true {

		isPassed := CheckArchiveByBandizip(fileFullPath)
		if ignoreUnZipErr == false && isPassed == false {
			return errors.New("CheckArchiveByBandizip == false")
		}

		if isPassed == true {
			err = UnArchiveByBandizip(fileFullPath, desRootPath)
			if ignoreUnZipErr == false && err != nil {
				return err
			}
		}

	} else {
		err = UnArchiveFile(fileFullPath, desRootPath)
		if ignoreUnZipErr == false && err != nil {
			return err
		}
	}
	// --------------------------------------------------
	doUnzipFun := func() error {
		// 判断一次
		needUnzipFileFPaths := make([]string, 0)
		err = filepath.Walk(desRootPath, func(path string, info fs.FileInfo, err error) error {

			if info.IsDir() == true {
				return nil
			}
			nowExt := filepath.Ext(path)
			// 然后对于解压的内容再次进行判断
			if nowExt != Zip.String() && nowExt != Tar.String() && nowExt != Rar.String() && nowExt != SevenZ.String() {
				return nil
			} else {
				needUnzipFileFPaths = append(needUnzipFileFPaths, path)
			}
			return nil
		})
		if err != nil {
			return err
		}
		// 如果有压缩包，那么就继续解压，然后删除压缩包
		for _, needUnzipFileFPath := range needUnzipFileFPaths {

			if useBandizip == true {

				isPassed := CheckArchiveByBandizip(needUnzipFileFPath)
				if isPassed == true {
					err = UnArchiveByBandizip(needUnzipFileFPath, desRootPath)
					if ignoreUnZipErr == false {
						if err != nil {
							return err
						}
					}
				}
			} else {
				err = UnArchiveFile(needUnzipFileFPath, desRootPath)
				if ignoreUnZipErr == true {
					if err != nil {
						return err
					}
				}
			}
			err = os.Remove(needUnzipFileFPath)
			if err != nil {
				return err
			}
		}

		return nil
	}
	// 第二次解压
	err = doUnzipFun()
	if ignoreUnZipErr == false && err != nil {
		return err
	}
	// 第三次解压
	err = doUnzipFun()
	if ignoreUnZipErr == false && err != nil {
		return err
	}
	// 第四次解压
	err = doUnzipFun()
	if ignoreUnZipErr == false && err != nil {
		return err
	}

	return nil
}

// UnArchiveFile 7z 以外的都能搞定中文编码的问题，但是 7z 有梗，需要单独的库去解析，且编码是解决不了的，以后他们搞定了再测试
// 所以效果就是，7z 外的压缩包文件解压ok，字幕可以正常从名称解析出是简体还是繁体，但是7z就没办法了，一定乱码
func UnArchiveFile(fileFullPath, desRootPath string) error {
	switch filepath.Ext(strings.ToLower(fileFullPath)) {
	case Zip.String():
		z := archiver.Zip{
			CompressionLevel:       flate.DefaultCompression,
			MkdirAll:               true,
			SelectiveCompression:   true,
			ContinueOnError:        false,
			OverwriteExisting:      true,
			ImplicitTopLevelFolder: false,
		}
		err := z.Walk(fileFullPath, func(f archiver.File) error {
			if f.IsDir() == true {
				return nil
			}

			zfh, ok := f.Header.(zip.FileHeader)
			if ok == true {
				err := processOneFile(f, zfh.NonUTF8, desRootPath)
				if err != nil {
					return err
				}
			} else {
				// 需要检测文件名是否是乱码
				err := processOneFile(f, !utf8.ValidString(f.Name()), desRootPath)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	case Tar.String():
		z := archiver.Tar{
			MkdirAll:               true,
			ContinueOnError:        false,
			OverwriteExisting:      true,
			ImplicitTopLevelFolder: false,
			StripComponents:        1,
		}
		err := z.Walk(fileFullPath, func(f archiver.File) error {
			if f.IsDir() == true {
				return nil
			}
			err := processOneFile(f, false, desRootPath)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	case Rar.String():
		z := archiver.Rar{
			MkdirAll:               true,
			ContinueOnError:        false,
			OverwriteExisting:      true,
			ImplicitTopLevelFolder: false,
			StripComponents:        1,
		}
		err := z.Walk(fileFullPath, func(f archiver.File) error {
			if f.IsDir() == true {
				return nil
			}
			err := processOneFile(f, false, desRootPath)
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return err
		}
	case SevenZ.String():
		return unArr7z(fileFullPath, desRootPath)
	default:
		return errors.New("not support un archive file ext")
	}

	return nil
}

func processOneFile(f archiver.File, notUTF8 bool, desRootPath string) error {
	decodeName := f.Name()
	if notUTF8 == true {

		//ouBytes, err := ChangeFileCoding2UTF8([]byte(f.Name()))
		//if err != nil {
		//	return err
		//}
		i := bytes.NewReader([]byte(f.Name()))
		decoder := transform.NewReader(i, simplifiedchinese.GB18030.NewDecoder())
		content, _ := io.ReadAll(decoder)
		decodeName = string(content)
		//decodeName = string(ouBytes)
	}
	var chunk []byte
	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		//说明读取结束
		if n == 0 {
			break
		}
		//读取到最终的缓冲区中
		chunk = append(chunk, buf[:n]...)
	}
	err := pkg.WriteFile(filepath.Join(desRootPath, decodeName), chunk)
	if err != nil {
		return err
	}
	return nil
}

func unArr7z(fileFullPath, desRootPath string) error {

	r, err := sevenzip.OpenReader(fileFullPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, file := range r.File {
		err = un7zOneFile(file, desRootPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func un7zOneFile(file *sevenzip.File, desRootPath string) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return err
	}
	decodeName := file.Name
	decodeName = filepath.Base(decodeName)
	err = pkg.WriteFile(filepath.Join(desRootPath, decodeName), data)
	if err != nil {
		return err
	}
	return nil
}

func IsWantedArchiveExtName(fileName string) bool {

	switch strings.ToLower(filepath.Ext(fileName)) {
	case Zip.String(), Tar.String(), Rar.String(), SevenZ.String():
		return true
	default:
		return false
	}
}

type ExtType int

const (
	Zip ExtType = iota + 1
	Tar
	Rar
	SevenZ
)

func (t ExtType) String() string {
	switch t {
	case Zip:
		return ".zip"
	case Tar:
		return ".tar"
	case Rar:
		return ".rar"
	case SevenZ:
		return ".7z"
	default:
		logger.Panicln("not support ext type")
		return ""
	}
}
