package epub

import (
	"archive/zip"
	"bytes"
	"testing"

	"github.com/readium/go-toolkit/pkg/archive"
	"github.com/readium/go-toolkit/pkg/fetcher"
	"github.com/readium/go-toolkit/pkg/manifest"
	"github.com/stretchr/testify/require"
)

func TestRootFileLiteralPath(t *testing.T) {
	const path = "OEBPS/package %20#?.opf"
	var original bytes.Buffer
	w := zip.NewWriter(&original)
	entry, err := w.Create("META-INF/container.xml")
	require.NoError(t, err)
	_, err = entry.Write([]byte(`<container xmlns="urn:oasis:names:tc:opendocument:xmlns:container"><rootfiles><rootfile full-path="` + path + `"/></rootfiles></container>`))
	require.NoError(t, err)
	entry, err = w.Create(path)
	require.NoError(t, err)
	_, err = entry.Write([]byte("original package"))
	require.NoError(t, err)
	require.NoError(t, w.Close())
	z, err := zip.NewReader(bytes.NewReader(original.Bytes()), int64(original.Len()))
	require.NoError(t, err)
	f := fetcher.NewArchiveFetcher(archive.NewGoZIPArchive(z, func() error { return nil }, false))
	t.Cleanup(f.Close)
	u, err := GetRootFilePath(t.Context(), f)
	require.NoError(t, err)
	require.Equal(t, path, u.Path())
	require.Empty(t, u.Raw().RawQuery)
	require.Empty(t, u.Fragment())
	body, rerr := f.Get(t.Context(), manifest.Link{Href: manifest.NewHREF(u)}).Read(t.Context(), 0, 0)
	require.Nil(t, rerr)
	require.Equal(t, "original package", string(body))
}
