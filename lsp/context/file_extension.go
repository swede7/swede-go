package context

import (
	"strings"
)

type FileExtension string

func (e FileExtension) String() string {
	return string(e)
}

const (
	Go    FileExtension = ".go"
	Swede FileExtension = ".feature"
)

func GetFileExtensionByURL(URL string) FileExtension {
	if strings.HasSuffix(URL, Go.String()) {
		return Go
	}
	if strings.HasSuffix(URL, Swede.String()) {
		return Swede
	}
	panic("unsupported file extension")
}
