package util

import (
	"io"
	"mime"
	"net/http"
	"os"
	"path"

	"github.com/mholt/archiver/v3"
)

// Download archive from URL & decompress
func DownloadArchive(url string, outPath string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	if err != nil {
		return err
	}
	tarName := params["filename"]

	// Create the file
	out, err := os.Create(path.Join(TmpPath, tarName))
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	err = archiver.Unarchive(out.Name(), TmpSrcPath)
	if err != nil {
		return err
	}

	return os.Remove(out.Name())
}
