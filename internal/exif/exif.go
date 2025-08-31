package exif

import (
	"fmt"
	"os"

	"github.com/dsoprea/go-exif"
)

func ParseExif(path string) (map[string]any, error) {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	exifData, err := exif.SearchAndExtractExif(fileData)
	if err != nil {
		return nil, err
	}

	exifTags, err := exif.GetFlatExifData(exifData)
	if err != nil {
		return nil, err
	}

	simplified := make(map[string]any)
	for i, tag := range exifTags {
		tagName := fmt.Sprintf("no_tag_name_%d", i)
		if tag.TagName != "" {
			tagName = tag.TagName
		}
		simplified[tagName] = tag.Value
	}

	return simplified, nil
}
