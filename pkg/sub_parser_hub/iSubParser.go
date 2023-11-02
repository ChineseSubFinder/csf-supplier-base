package sub_parser_hub

import subparser "github.com/ChineseSubFinder/csf-supplier-base/pkg/sub_parser_hub/sub_parser"

type ISubParser interface {
	GetParserName() string

	DetermineFileTypeFromFile(filePath string) (bool, *subparser.FileInfo, error)

	DetermineFileTypeFromBytes(inBytes []byte) (bool, *subparser.FileInfo, error)
}
