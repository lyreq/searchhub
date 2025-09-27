package provider_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"lytemp/pkg/provider"

	"github.com/stretchr/testify/require"
)

const xmlFeed = `<?xml version="1.0" encoding="UTF-8"?>
<feed>
  <items>
    <item>
      <id>a1</id>
      <headline>Clean Architecture in Go</headline>
      <type>article</type>
      <stats>
        <reading_time>8</reading_time>
        <reactions>450</reactions>
        <comments>25</comments>
      </stats>
      <publication_date>2024-03-14</publication_date>
      <categories>
        <category>programming</category>
        <category>architecture</category>
      </categories>
    </item>
  </items>
  <meta>
    <total_count>75</total_count>
    <current_page>1</current_page>
    <items_per_page>10</items_per_page>
  </meta>
</feed>`

func TestXMLProvider_Fetch_Normalize(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		_, _ = w.Write([]byte(xmlFeed))
	}))
	defer srv.Close()

	p := provider.NewXMLProvider(srv.URL, provider.HTTPClient())
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	page, err := p.Fetch(ctx, 1, 10, "")
	require.NoError(t, err)
	require.Equal(t, 75, page.Total)
	require.Len(t, page.Items, 1)

	it := page.Items[0]
	require.Equal(t, "xml_provider", it.Provider)
	require.Equal(t, "a1", it.ProviderItemID)
	require.Equal(t, "article", it.NormalizedType)
	require.Equal(t, "Clean Architecture in Go", it.Title)
	require.Nil(t, it.Views)
	require.Nil(t, it.Likes)
	require.Nil(t, it.DurationSeconds)
	require.NotNil(t, it.ReadingTimeMinutes)
	require.NotNil(t, it.Reactions)
	require.NotNil(t, it.Comments)
	require.Equal(t, []string{"programming", "architecture"}, it.Tags)

	// publication_date → UTC midnight
	require.Equal(t, time.Date(2024, 3, 14, 0, 0, 0, 0, time.UTC), it.PublishedAt)
}
