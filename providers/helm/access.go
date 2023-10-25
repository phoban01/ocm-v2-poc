package helm

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"

	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/types"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
)

type accessor struct {
	initialized bool
	path        string
	chart       string
	version     string
	accessType  string
	mediaType   string
	labels      map[string]string
	data        io.ReadCloser
	digest      *types.Digest
	length      int64
}

var _ v2.Access = (*accessor)(nil)

func (a *accessor) compute() error {
	if a.initialized {
		return nil
	}

	cachedir, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	configdir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	dl := &downloader.ChartDownloader{
		Out:              os.Stderr,
		Getters:          getter.All(&cli.EnvSettings{}),
		RepositoryCache:  filepath.Join(cachedir, "helm/repository"),
		RepositoryConfig: filepath.Join(configdir, "helm/repositories.yaml"),
	}

	u, err := dl.ResolveChartVersion(a.chart, a.version)
	if err != nil {
		return err
	}

	g, err := dl.Getters.ByScheme(u.Scheme)
	if err != nil {
		return err
	}

	data, err := g.Get(u.String(), dl.Options...)
	if err != nil {
		return err
	}

	a.length = int64(data.Len())

	a.data = io.NopCloser(data)

	hash := sha256.New()

	_, err = io.Copy(hash, bytes.NewReader(data.Bytes()[:]))
	if err != nil {
		return err
	}

	a.digest = &types.Digest{
		HashAlgorithm:          "sha256",
		NormalisationAlgorithm: "json/v1",
		Value:                  fmt.Sprintf("%x", hash.Sum(nil)),
	}

	a.initialized = true

	return nil
}

func (a *accessor) Type() string {
	if a.accessType != "" {
		return a.accessType
	}
	return LocalAccessType
}

func (a *accessor) Reference() string {
	return a.chart
}

func (a *accessor) MediaType() string {
	if a.mediaType != "" {
		return a.mediaType
	}
	return MediaType
}

func (a *accessor) Length() (int64, error) {
	if err := a.compute(); err != nil {
		return 0, err
	}
	return a.length, nil
}

func (a *accessor) Labels() map[string]string {
	return a.labels
}

func (a *accessor) Data() (io.ReadCloser, error) {
	if err := a.compute(); err != nil {
		return nil, err
	}
	return a.data, nil
}

func (a *accessor) Digest() (*types.Digest, error) {
	if err := a.compute(); err != nil {
		return nil, err
	}
	return a.digest, nil
}
