package sub_parser_hub

import (
	"github.com/ChineseSubFinder/csf-supplier/pkg/sub_parser_hub/sub_parser"
)

type ISubParser interface {
	GetParserName() string

	DetermineFileTypeFromFile(filePath string) (bool, *subparser.FileInfo, error)

	DetermineFileTypeFromBytes(inBytes []byte, nowExt string) (bool, *subparser.FileInfo, error)
}
