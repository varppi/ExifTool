package pdf

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	extract = []string{"1 0 obj", "6 0 obj", "5 0 obj"}
)

type pdfMetaSection struct {
	Name string
	Tags []string
}

func ParseMeta(path string) (map[string]string, error) {
	handle, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer handle.Close()

	buffer := make([]byte, 1024*50)
	var validPdf bool
	for {
		n, err := handle.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			return nil, err
		}
		buffer = buffer[:n]
		if bytes.Contains(buffer, []byte("stream\n")) {
			validPdf = true
			break
		}
	}
	if !validPdf {
		return nil, nil
	}

	var pureSections []pdfMetaSection
	for _, section := range extract {
		if !bytes.Contains(buffer, []byte(section)) {
			continue
		}
		afterSectionStart := bytes.Split(buffer, []byte(section))[1]
		beforeSectionEnd := bytes.Split(afterSectionStart, []byte("stream"))[0]
		beforeSectionEnd = bytes.Split(beforeSectionEnd, []byte("endobj"))[0]
		finalMiddle := string(beforeSectionEnd)
		if strings.Contains(finalMiddle, ">>") {
			finalMiddle = strings.TrimPrefix(finalMiddle, "<<")
			finalMiddle = strings.TrimSuffix(finalMiddle, ">>")
		}

		finalMiddle = strings.ReplaceAll(finalMiddle, "\\)", "|curve_b|")
		pureSections = append(pureSections, pdfMetaSection{
			Name: section,
			Tags: strings.Split(strings.TrimSpace(finalMiddle), ")"),
		})
	}

	metadata := make(map[string]string)
	for _, section := range pureSections {
		for i, tag := range section.Tags {
			if len(strings.TrimSpace(tag)) < 3 {
				continue
			}
			var name string
			var value string
			if strings.Contains(tag, "(") {
				name = strings.ToLower(strings.Split(tag, "(")[0])
				name = strings.TrimSpace(name)
				nameSections := strings.Split(name, "/")
				name = nameSections[len(nameSections)-1]
				value = strings.Join(strings.Split(tag, "(")[1:], "(")
			}
			value = strings.Split(value, "\n")[0]
			value = strings.TrimPrefix(value, "(")
			value = strings.TrimSuffix(value, ")")
			value = strings.ReplaceAll(value, "\\(", "(")
			value = strings.ReplaceAll(value, "|curve_b|", ")")
			if name == "" && value == "" {
				continue
			}
			if name == "" {
				name = fmt.Sprintf("unnamed_tag_%d", i)
			}
			metadata[name] = value
		}
	}

	if len(metadata) == 0 {
		return nil, nil
	}
	return metadata, nil
}
