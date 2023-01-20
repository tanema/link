package link

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testHeader = `<http://red.com>; other="foo"; rel="first", <http://purple.com>; rel="last", <http://green.com>; rel="next", <http://yellow.com>; rel="prev", <http://cyan.com>; rel="home"`

func TestNewHeader(t *testing.T) {
	header := NewHeader(map[string]*url.URL{
		"first":   u("http://red.com"),
		"last":    u("http://purple.com"),
		"home":    u("http://cyan.com"),
		"next":    u("http://green.com"),
		"prev":    u("http://yellow.com"),
		"another": nil,
	})
	assert.Equal(t, &Link{URL: u("http://red.com"), Params: map[string]string{"rel": "first"}}, header.First())
	assert.Equal(t, &Link{URL: u("http://purple.com"), Params: map[string]string{"rel": "last"}}, header.Last())
	assert.Equal(t, &Link{URL: u("http://green.com"), Params: map[string]string{"rel": "next"}}, header.Next())
	assert.Equal(t, &Link{URL: u("http://yellow.com"), Params: map[string]string{"rel": "prev"}}, header.Prev())
	assert.Equal(t, &Link{URL: u("http://cyan.com"), Params: map[string]string{"rel": "home"}}, header.Find("home"))
	assert.Len(t, header.Links, 5)
}

func TestParse(t *testing.T) {
	req := &http.Response{Header: http.Header{"Link": {" " + testHeader}}}
	header, err := Parse(req)
	assert.Nil(t, err)
	assert.Equal(t, &Link{URL: u("http://red.com"), Params: map[string]string{"rel": "first", "other": "foo"}}, header.First())
	assert.Equal(t, &Link{URL: u("http://purple.com"), Params: map[string]string{"rel": "last"}}, header.Last())
	assert.Equal(t, &Link{URL: u("http://green.com"), Params: map[string]string{"rel": "next"}}, header.Next())
	assert.Equal(t, &Link{URL: u("http://yellow.com"), Params: map[string]string{"rel": "prev"}}, header.Prev())
	assert.Equal(t, &Link{URL: u("http://cyan.com"), Params: map[string]string{"rel": "home"}}, header.Find("home"))
	assert.Len(t, header.Links, 5)
}

func TestParseString(t *testing.T) {
	t.Run("without error", func(t *testing.T) {
		header, err := ParseString(" " + testHeader)
		assert.Nil(t, err)
		assert.Equal(t, &Link{URL: u("http://red.com"), Params: map[string]string{"rel": "first", "other": "foo"}}, header.First())
		assert.Equal(t, &Link{URL: u("http://purple.com"), Params: map[string]string{"rel": "last"}}, header.Last())
		assert.Equal(t, &Link{URL: u("http://green.com"), Params: map[string]string{"rel": "next"}}, header.Next())
		assert.Equal(t, &Link{URL: u("http://yellow.com"), Params: map[string]string{"rel": "prev"}}, header.Prev())
		assert.Equal(t, &Link{URL: u("http://cyan.com"), Params: map[string]string{"rel": "home"}}, header.Find("home"))
		assert.Len(t, header.Links, 5)
	})

	t.Run("expecting error", func(t *testing.T) {
		_, err := ParseString(" <:/fooboar>; rel=first")
		assert.NotNil(t, err)
	})

	t.Run("handles chunk newlines", func(t *testing.T) {
		const newlineHeader = `<http://red.com>; rel=first; other=foo,
<http://purple.com>; rel=last,
<http://green.com>; rel=next,
<http://yellow.com>; rel=prev,
<http://cyan.com>; rel=home`

		header, err := ParseString(newlineHeader)
		assert.Nil(t, err)
		assert.Len(t, header.Links, 5)
		assert.Equal(t, &Link{URL: u("http://red.com"), Params: map[string]string{"rel": "first", "other": "foo"}}, header.First())
		assert.Equal(t, &Link{URL: u("http://purple.com"), Params: map[string]string{"rel": "last"}}, header.Last())
		assert.Equal(t, &Link{URL: u("http://green.com"), Params: map[string]string{"rel": "next"}}, header.Next())
		assert.Equal(t, &Link{URL: u("http://yellow.com"), Params: map[string]string{"rel": "prev"}}, header.Prev())
		assert.Equal(t, &Link{URL: u("http://cyan.com"), Params: map[string]string{"rel": "home"}}, header.Find("home"))
	})
}

func TestHeader_Find(t *testing.T) {
	req := &http.Response{Header: http.Header{"Link": {" " + testHeader}}}
	header, err := Parse(req)
	assert.Nil(t, err)
	assert.Equal(t, &Link{URL: u("http://red.com"), Params: map[string]string{"rel": "first", "other": "foo"}}, header.Find("first"))
	assert.Equal(t, &Link{URL: u("http://purple.com"), Params: map[string]string{"rel": "last"}}, header.Find("last"))
	assert.Equal(t, &Link{URL: u("http://green.com"), Params: map[string]string{"rel": "next"}}, header.Find("next"))
	assert.Equal(t, &Link{URL: u("http://yellow.com"), Params: map[string]string{"rel": "prev"}}, header.Find("prev"))
	assert.Equal(t, &Link{URL: u("http://cyan.com"), Params: map[string]string{"rel": "home"}}, header.Find("home"))
	assert.Nil(t, header.Find("notThere"))
}

func TestHeader_String(t *testing.T) {
	header := &Header{
		Links: []*Link{
			{URL: u("http://red.com"), Params: map[string]string{"rel": "first", "other": "foo"}},
			{URL: u("http://purple.com"), Params: map[string]string{"rel": "last"}},
			{URL: u("http://green.com"), Params: map[string]string{"rel": "next"}},
			{URL: u("http://yellow.com"), Params: map[string]string{"rel": "prev"}},
			{URL: u("http://cyan.com"), Params: map[string]string{"rel": "home"}},
		},
	}
	assert.Equal(t, testHeader, header.String())
}

func TestLink_String(t *testing.T) {
	link := &Link{URL: u("http://red.com"), Params: map[string]string{"rel": "first", "other": "foo"}}
	expected := `<http://red.com>; other="foo"; rel="first"`
	assert.Equal(t, expected, link.String())
}

func u(path string) *url.URL {
	parsed, err := url.Parse(path)
	if err != nil {
		panic(err)
	}
	return parsed
}
