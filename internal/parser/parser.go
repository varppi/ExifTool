package parser

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/varppi/ExifTool/internal/exif"
	"github.com/varppi/ExifTool/internal/pdf"
)

type Metadata struct {
	Pdf      map[string]string
	Exif     map[string]any
	ExifJson string
	PdfJson  string
}

func ParseFileMetadata(file string) (*Metadata, error) {
	info, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	if info.Size() >= 1073741824 {
		return nil, errors.New("too big file")
	}

	exifTags, imgError := exif.ParseExif(file)
	pdfMeta, pdfError := pdf.ParseMeta(file)
	if imgError != nil && pdfError != nil {
		return nil, errors.New("no metadata found")
	}

	resultMetadata := Metadata{}
	if imgError == nil {
		exifJsonOut, err := json.Marshal(exifTags)
		if err != nil {
			return nil, err
		}
		resultMetadata.Exif = exifTags
		resultMetadata.ExifJson = string(exifJsonOut)
	}
	if pdfError == nil {
		pdfJsonOut, err := json.Marshal(pdfMeta)
		if err != nil {
			return nil, err
		}
		resultMetadata.Pdf = pdfMeta
		resultMetadata.PdfJson = strings.TrimPrefix(string(pdfJsonOut), "null")
	}

	return &resultMetadata, nil
}
