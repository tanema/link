package link

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"golang.org/x/exp/slices"
)

type (
	HeaderConfig struct {
		First string
		Last  string
		Next  string
		Prev  string
	}
	// Header contains all of the links in the header with ways of accessing them
	Header struct {
		Links []*Link
	}
	// Link contains a single link with it's associated data.
	Link struct {
		URL    *url.URL
		Params map[string]string
	}
)

// NewHeader will format a new link header with a map of rel = link
func NewHeader(links map[string]*url.URL) *Header {
	header := &Header{}
	for rel, link := range links {
		if link == nil {
			continue
		}
		header.Links = append(header.Links, &Link{
			URL:    link,
			Params: map[string]string{"rel": rel},
		})
	}
	return header
}

// Parse will retrieve the link header from a response and parse it.
func Parse(req *http.Response) (*Header, error) {
	return ParseString(req.Header.Get("Link"))
}

// ParseString will parse a string from the link header and find all appropriate
// infomation from link headers
func ParseString(header string) (*Header, error) {
	linkHeader := &Header{Links: make([]*Link, 0, strings.Count(header, ",")+1)}
	header = strings.ReplaceAll(header, "\n", "")
	for _, part := range strings.Split(header, ",") {
		parts := strings.Split(part, ";")
		u, err := url.Parse(strings.Trim(parts[0], " <>"))
		if err != nil {
			return nil, err
		}
		params := map[string]string{}
		for _, param := range parts[1:] {
			kv := strings.SplitN(strings.TrimSpace(param), "=", 2)
			params[strings.Trim(kv[0], `"`)] = strings.Trim(kv[1], `"`)
		}
		link := &Link{
			URL:    u,
			Params: params,
		}
		linkHeader.Links = append(linkHeader.Links, link)
	}
	return linkHeader, nil
}

// First will find the link for the link to the first page with rel="first" if
// it exists. It will return nil if it does not exist.
func (header Header) First() *Link {
	return header.Find("first")
}

// Last will find the link for the link to the last page with rel="last" if
// it exists. It will return nil if it does not exist.
func (header Header) Last() *Link {
	return header.Find("last")
}

// Next will find the link for the link to the next page with rel="next" if
// it exists. It will return nil if it does not exist.
func (header Header) Next() *Link {
	return header.Find("next")
}

// Prev will find the link for the link to the previous with rel="prev" or
// rel="previous"" if it exists. It will return nil if it does not exist.
func (header Header) Prev() *Link {
	return header.Find("prev", "previous")
}

// Find will search for a link with the rel matching one of the values passed.
// If it does not exist then nil will be returned.
func (header Header) Find(rel ...string) *Link {
	for _, link := range header.Links {
		if link.Params["rel"] != "" && slices.Contains(rel, link.Params["rel"]) {
			return link
		}
	}
	return nil
}

func (header Header) String() string {
	links := make([]string, len(header.Links))
	for i, link := range header.Links {
		links[i] = link.String()
	}
	return strings.Join(links, ", ")
}

func (link Link) String() string {
	parts := make([]string, 0, len(link.Params)+1)
	parts = append(parts, fmt.Sprintf("<%s>", link.URL))
	// Get a list of sorted keys so that the output is deterministic
	keys := make([]string, 0, len(link.Params))
	for key := range link.Params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf(`%v="%v"`, key, link.Params[key]))
	}
	return strings.Join(parts, "; ")
}
