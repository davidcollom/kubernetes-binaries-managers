package binary

import (
	"fmt"
	"io"
	"math/big"

	// "math/rand/v2"
	"crypto/rand"
	"net/http"
	"os"
	"strings"

	"path/filepath"

	. "github.com/little-angry-clouds/kubernetes-binaries-managers/internal/helpers" // nolint:staticcheck
	"github.com/mholt/archiver/v3"
)

const (
	zip   = ".zip"
	targz = ".tar.gz"
	exe   = ".exe"
)

type DownloadError struct {
	Err  string
	URL  string
	Body string
}

func (e *DownloadError) Error() string {
	var error string
	if e.Body == "" {
		error = fmt.Sprintf("%s\nurl: %s", e.Err, e.URL)
	} else {
		error = fmt.Sprintf("%s\nurl: %s\nbody: %s", e.Err, e.URL, e.Body)
	}

	return error
}

func Download(version string, url string) ([]byte, error) { // nolint: funlen
	var (
		osArch string
		err    error
		body   []byte

		errorCodeFail = 404
		errorCodePass = 200
	)

	// Don't control the error since at this point it should be controlled
	osArch, _ = GetOSArch()
	os := strings.Split(osArch, "/")[0]
	arch := strings.Split(osArch, "/")[1]

	if strings.Contains(url, "openshift") {
		// They use mac instead of darwin in the url
		switch os {
		case "darwin":
			os = "mac"
		case "windows":
			url = strings.Replace(url, ".tar.gz", ".zip", 1)
		}

		url = fmt.Sprintf(url, version, os, version)
	} else {
		url = fmt.Sprintf(url, version, os, arch)
	}

	if strings.Contains(url, "helm") {
		if strings.Contains(osArch, "windows") {
			url += zip
		} else {
			url += targz
		}
	} else if strings.Contains(url, "kubectl") {
		if strings.Contains(osArch, "windows") {
			url += exe
		}
	}

	fmt.Println("Downloading binary...")

	resp, err := http.Get(url) // nolint
	if err != nil {
		return body, err
	}

	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}

	if resp.StatusCode == errorCodeFail {
		return body, &DownloadError{"binary not found", url, string(body)}
	} else if resp.StatusCode != errorCodePass {
		return body, &DownloadError{"unhandled error", url, string(body)}
	}

	return body, nil
}

func Save(fileName string, body []byte) error { // nolint: funlen
	var err error

	// helm returns a compressed file, so save it somewhere and decompress it
	if strings.Contains(fileName, "helm") {
		var fileExt string

		randomNumbers := int64(5000) // nolint:mnd

		randomInt, err := rand.Int(rand.Reader, big.NewInt(randomNumbers))
		if err != nil {
			return err
		}

		tempDir, err := os.MkdirTemp("", "helm")
		if err != nil {
			return err
		}
		// clean temp dir
		defer os.RemoveAll(tempDir)

		osArch, _ := GetOSArch()
		file := fmt.Sprintf("%s/helm-%d", tempDir, randomInt.Int64())
		file, _ = filepath.Abs(file)

		if strings.Contains(osArch, "windows") {
			fileExt = zip
		} else {
			fileExt = targz
		}

		err = os.WriteFile(file+fileExt, body, 0750) // nolint: gosec,mnd
		if err != nil {
			return err
		}

		err = archiver.Unarchive(file+fileExt, file)
		if err != nil {
			return err
		}

		OS := strings.Split(osArch, "/")[0]
		arch := strings.Split(osArch, "/")[1]
		path, _ := filepath.Abs(file + fmt.Sprintf("/%s-%s/helm", OS, arch))

		if strings.Contains(osArch, "windows") {
			path += exe
		}

		body, err = os.ReadFile(path)
		if err != nil {
			return err
		}
	} else if strings.Contains(fileName, "okd") {
		// oc returns a compressed file, so save it somewhere and decompress it
		var fileExt string

		randomNumbers := int64(5000) // nolint:mnd

		randomInt, err := rand.Int(rand.Reader, big.NewInt(randomNumbers))
		if err != nil {
			return err
		}

		tempDir, err := os.MkdirTemp("", "oc")
		if err != nil {
			return err
		}
		// clean temp dir
		defer os.RemoveAll(tempDir)

		osArch, _ := GetOSArch()
		file := fmt.Sprintf("%s/oc-%d", tempDir, randomInt.Int64())
		file, _ = filepath.Abs(file)

		if strings.Contains(osArch, "windows") {
			fileExt = zip
		} else {
			fileExt = targz
		}

		err = os.WriteFile(file+fileExt, body, 0750) // nolint: gosec,mnd
		if err != nil {
			return err
		}

		err = archiver.Unarchive(file+fileExt, file)
		if err != nil {
			return err
		}

		path, _ := filepath.Abs(file + "/oc")

		if strings.Contains(osArch, "windows") {
			path += exe
		}

		body, err = os.ReadFile(path)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(fileName, body, 0750) // nolint: gosec,mnd
	if err != nil {
		return err
	}

	return nil
}
