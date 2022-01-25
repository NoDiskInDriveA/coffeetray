package coffeetray

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
)

func GetPngIconBuffer(path string) (*bytes.Buffer, error) {
	iconFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer iconFile.Close()

	fileType, err := GetFileContentType(iconFile)
	if err != nil {
		return nil, err
	}

	if fileType != "image/png" {
		return nil, fmt.Errorf("given image not a PNG file, found %s", fileType)
	}

	imageData, err := png.Decode(iconFile)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, imageData, nil)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func GetFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}
	defer out.Seek(0, 0)
	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
