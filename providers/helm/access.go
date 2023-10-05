package helm

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/types"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
)

func (a *accessor) Type() v2.AccessType {
	return v2.AccessType("helmChart/v1")
}

func (a *accessor) MediaType() string {
	return MediaType
}

func (a *accessor) Labels() map[string]string {
	return a.labels
}

func (a *accessor) Data() (io.ReadCloser, error) {
	dl := &downloader.ChartDownloader{
		Out:     os.Stderr,
		Getters: getter.All(&cli.EnvSettings{}),
	}
	tmp, err := os.MkdirTemp(os.TempDir(), "ocm-test-v2-*")
	if err != nil {
		return nil, err
	}
	if _, _, err := dl.DownloadTo(a.chart, a.version, tmp); err != nil {
		return nil, err
	}
	return os.Open(tmp)
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

func (a *accessor) WithLocation(p string) {
	// a.filepath = p
}
