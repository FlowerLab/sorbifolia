package sitemap

import (
	"encoding/xml"
	"io"
)

const xmlNS = "http://www.sitemaps.org/schemas/sitemap/0.9"

// Sitemap
//
// https://www.sitemaps.org/protocol.html
type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	NS      string   `xml:"xmlns,attr"`
	URL     []*URL   `xml:"url"`
}

type URL struct {
	Loc        string      `xml:"loc"`
	LastMod    *string     `xml:"lastmod,omitempty"`
	ChangeFreq *ChangeFreq `xml:"changefreq,omitempty"`
	Priority   *float32    `xml:"priority,omitempty"` // Valid values range from 0.0 to 1.0
}

type ChangeFreq string

const (
	Always  ChangeFreq = "always"
	Hourly  ChangeFreq = "hourly"
	Daily   ChangeFreq = "daily"
	Weekly  ChangeFreq = "weekly"
	Monthly ChangeFreq = "monthly"
	Yearly  ChangeFreq = "yearly"
	Never   ChangeFreq = "never"
)

func NewURL(url string) *URL                            { return &URL{Loc: url} }
func (u *URL) SetLastMod(lastMod string) *URL           { u.LastMod = &lastMod; return u }
func (u *URL) SetChangeFreq(changeFreq ChangeFreq) *URL { u.ChangeFreq = &changeFreq; return u }
func (u *URL) SetPriority(priority float32) *URL        { u.Priority = &priority; return u }

func New() *Sitemap           { return &Sitemap{NS: xmlNS} }
func (s *Sitemap) Add(u *URL) { s.URL = append(s.URL, u) }
func (s *Sitemap) Write(writer io.Writer) error {
	xmlEncoder := xml.NewEncoder(writer)
	xmlEncoder.Indent("", "  ")

	if _, err := writer.Write([]byte(xml.Header)); err != nil {
		return err
	}

	if err := xmlEncoder.Encode(s); err != nil {
		return err
	}

	return nil
}

func Parse(reader io.Reader) (*Sitemap, error) {
	var arg Sitemap
	if err := xml.NewDecoder(reader).Decode(&arg); err != nil {
		return nil, err
	}
	return &arg, nil
}

func (s *Sitemap) AddOrUpdate(u *URL) {
	if u == nil {
		return
	}
	for _, v := range s.URL {
		if v.Loc != u.Loc {
			continue
		}
		if v.LastMod != nil && u.LastMod != nil {
			v.LastMod = u.LastMod
		}
		if v.ChangeFreq != nil && u.ChangeFreq != nil {
			v.ChangeFreq = u.ChangeFreq
		}
		if v.Priority != nil && u.Priority != nil {
			v.Priority = u.Priority
		}
		return
	}
	s.Add(u)
}

// Merge ns has priority less than s
func (s *Sitemap) Merge(ns *Sitemap) {
	if ns == nil {
		return
	}
	arr := make([]*URL, 0, len(ns.URL))
	for _, nv := range ns.URL {
		for _, v := range s.URL {
			if nv.Loc == v.Loc {
				goto CONTINUE
			}
		}
		arr = append(arr, nv)
	CONTINUE:
	}

	if len(arr) > 0 {
		s.URL = append(s.URL, arr...)
	}
}
