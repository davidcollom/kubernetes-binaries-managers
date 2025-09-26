package versions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"runtime"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/go-version"
	. "github.com/little-angry-clouds/kubernetes-binaries-managers/internal/helpers" // nolint:staticcheck
	"github.com/little-angry-clouds/kubernetes-binaries-managers/internal/logging"
	"github.com/mitchellh/go-homedir"
)

type Page struct {
	Release string `json:"tag_name"`
}

const (
	httpTimeout = 10 * time.Second
	httpRetries = 3
)

func SortVersions(versions []*version.Version, allReleases bool, allVersions bool) ([]*version.Version, error) {
	var (
		numberOfVersion int
		finalVersions   []*version.Version
	)

	sort.Sort(sort.Reverse(version.Collection(versions)))

	if allVersions {
		numberOfVersion = len(versions) - 1
	} else {
		numberOfVersion = 19
	}

	for i := 0; i <= numberOfVersion; i++ {
		if i == len(versions) {
			break
		}

		if allReleases {
			finalVersions = append(finalVersions, versions[i])
		} else {
			if !strings.ContainsAny(versions[i].String(), "beta") &&
				!strings.ContainsAny(versions[i].String(), "alpha") &&
				!strings.ContainsAny(versions[i].String(), "rc") {
				finalVersions = append(finalVersions, versions[i])
			} else {
				versions = append(versions[:i], versions[i+1:]...)
				i--
			}
		}
	}

	return finalVersions, nil
}

func PrintVersions(versions []*version.Version) {
	for _, element := range versions {
		fmt.Println(element)
	}
}

func GetLocalVersions(binary string) ([]*version.Version, error) {
	var versions []*version.Version // nolint:prealloc

	logging.Debug("GetLocalVersions called", "binary", binary)

	home, _ := homedir.Dir()
	binDir, _ := filepath.Abs(fmt.Sprintf("%s/.bin/%s-v*", home, binary))
	logging.Debug("searching for local versions", "binDir", binDir)
	matches, _ := filepath.Glob(binDir)
	logging.Debug("found local versions", "matches", matches)

	for _, match := range matches {
		v := strings.Split(match, string(os.PathSeparator))
		vs := strings.Replace(v[len(v)-1], binary+"-v", "", 1)

		if runtime.GOOS == "windows" {
			vs = strings.Replace(vs, ".exe", "", 1)
		}

		ver, err := version.NewVersion(vs)
		if err != nil {
			return versions, err
		}

		versions = append(versions, ver)
	}

	return versions, nil
}

func processPage(body io.Reader) ([]*version.Version, error) {
	var ( //nolint:prealloc
		pageVersions []*version.Version
		rel          []Page
	)

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &rel); err != nil {
		return nil, err
	}

	for _, element := range rel {
		v, err := version.NewVersion(element.Release)
		if err != nil {
			return nil, err
		}

		pageVersions = append(pageVersions, v)
	}

	return pageVersions, nil
}

func GetRemoteVersions(endpoint string) ([]*version.Version, error) {
	var versions []*version.Version

	client := retryablehttp.NewClient()
	client.RetryMax = httpRetries
	client.HTTPClient.Timeout = httpTimeout
	client.Backoff = backoffHandler
	client.CheckRetry = retryPolicy
	client.HTTPClient.Transport = &AuthRoundTripper{
		token:            os.Getenv("GITHUB_TOKEN"),
		nextRoundTripper: http.DefaultTransport,
	}

	client.Logger = nil

	// Fetch the first page
	logging.Debug("fetching first page", "endpoint", endpoint+"1")

	resp, err := client.Get(endpoint + "1")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		fmt.Println("The request to Github's API failed, sorry.")
		fmt.Println("You may still install the version you want, if you know it. It will always go as X.Y.Z.")
		fmt.Println("The complete request response is ", resp)

		return nil, errors.New("request to Github's API failed with 403 Forbidden")
	}

	lastPage, err := GetLastPage(resp.Header.Get("Link"))
	if err != nil {
		return nil, err
	}

	logging.Debug("last page determined", "lastPage", lastPage)

	// Process first page
	pageVersions, err := processPage(resp.Body)
	if err != nil {
		return nil, err
	}

	versions = append(versions, pageVersions...)
	logging.Debug("Processed page\n", "versions", len(pageVersions), "page", 1, "of", lastPage)

	// Process remaining pages
	for page := 2; page <= lastPage; page++ {
		logging.Debug("fetching page", "endpoint", endpoint+"1")

		resp, err := client.Get(endpoint + strconv.Itoa(page))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		pageVersions, err := processPage(resp.Body)
		if err != nil {
			return nil, err
		}

		versions = append(versions, pageVersions...)
		logging.Debug("Processed page", "page", page, "of", lastPage, "with", len(pageVersions))
	}

	logging.Debug("Found releases", "count", len(versions), "over pages", lastPage)

	return versions, nil
}
