package helpers

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
)

// Checksum check for a file
//
// Adds a if statement to exit if checksum passes
func Checksum(data []byte, path string) []string {
	checksum := sha1.Sum(data)
	return []string{
		fmt.Sprintf(`echo "%x %s" | sha1sum -c -`, checksum, path),
		`if [ $? = 0 ]; then exit(40); if`,
	}
}

// WriteEncodedFile will encode an array of bytes
// and create the command lines to write to a file
func WriteEncodedFile(data []byte, path string) []string {
	encodedData := base64.StdEncoding.EncodeToString(data)
	return []string{
		fmt.Sprintf("filebase64encoded=%s", encodedData),
		fmt.Sprintf("echo $filebase64encoded | base64 -d > %s", path),
	}
}
