package util

import (
	"os"
	"path"

	"github.com/cavaliercoder/grab"
	"github.com/go-vela/archiver/v3"
)

// Download archive from URL & decompress
func DownloadArchive(url string, outPath string) error {
	// Get the data
	resp, err := grab.Get(path.Join(TmpPath), url)
	if err != nil {
		return err
	}

	err = archiver.Unarchive(resp.Filename, TmpSrcPath)
	if err != nil {
		return err
	}

	return os.Remove(resp.Filename)
}
