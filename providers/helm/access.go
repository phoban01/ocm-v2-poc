package helm

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"github.com/phoban01/ocm-v2/api/v2/types"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
)

func (a *accessor) Type() string {
	if a.accessType != "" {
		return a.accessType
	}
	return LocalAccessType
}

func (a *accessor) MediaType() string {
	if a.mediaType != "" {
		return a.mediaType
	}
	return MediaType
}

func (a *accessor) Labels() map[string]string {
	return a.labels
}

func (a *accessor) Data() (io.ReadCloser, error) {
	if a.path != "" {
		return os.Open(a.path)
	}
	dl := &downloader.ChartDownloader{
		Out:              os.Stderr,
		Getters:          getter.All(&cli.EnvSettings{}),
		RepositoryCache:  "/home/piaras/.cache/helm/repository",
		RepositoryConfig: "/home/piaras/.config/helm/repositories.yaml",
	}

	tmp, err := os.MkdirTemp(os.TempDir(), "ocm-test-v2-*")
	if err != nil {
		return nil, err
	}

	dst, _, err := dl.DownloadTo(a.chart, a.version, tmp)
	if err != nil {
		return nil, err
	}

	return os.Open(dst)
}

func (a *accessor) Digest() (*types.Digest, error) {
	data, err := a.Data()
	if err != nil {
		return nil, err
	}
	defer data.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, data)
	if err != nil {
		return nil, err
	}

	return &types.Digest{
		HashAlgorithm:          "sha256",
		NormalisationAlgorithm: "json/v1",
		Value:                  fmt.Sprintf("%x", hash.Sum(nil)),
	}, nil
}
