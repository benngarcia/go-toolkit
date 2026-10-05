package fetcher

import (
	"archive/zip"
	"bytes"
	"testing"

	"github.com/readium/go-toolkit/pkg/archive"
	"github.com/stretchr/testify/require"
)

func TestArchiveFetcherLiteralPaths(t *testing.T) {
	names := []string{"space chapter.xhtml", "café😀.xhtml", "100% true.xhtml", "literal%20.xhtml", "hash#query?.xhtml"}
	var original bytes.Buffer
	w := zip.NewWriter(&original)
	for _, name := range names {
		entry, err := w.Create(name)
		require.NoError(t, err)
		_, err = entry.Write([]byte(name))
		require.NoError(t, err)
	}
	require.NoError(t, w.Close())
	z, err := zip.NewReader(bytes.NewReader(original.Bytes()), int64(original.Len()))
	require.NoError(t, err)
	f := NewArchiveFetcher(archive.NewGoZIPArchive(z, func() error { return nil }, false))
	t.Cleanup(f.Close)
	links, err := f.Links(t.Context())
	require.NoError(t, err)
	require.Len(t, links, len(names))
	for _, link := range links {
		u := link.Href.Resolve(nil, nil)
		require.Empty(t, u.Raw().RawQuery)
		require.Empty(t, u.Fragment())
		body, err := f.Get(t.Context(), link).Read(t.Context(), 0, 0)
		require.Nil(t, err)
		require.Equal(t, u.Path(), string(body))
	}
}
