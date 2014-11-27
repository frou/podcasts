package podcasts

import (
	"encoding/xml"
	"io"
	"time"
)

const (
	rssXmlns   = "http://www.itunes.com/dtds/podcast-1.0.dtd"
	rssVersion = "2.0"
	rfc2822    = "Mon, 02 Jan 2006 15:04:05 -0700"
)

type PubDate struct {
	time.Time
}

func (p PubDate) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	e.EncodeToken(start)
	e.EncodeToken(xml.CharData([]byte(p.Format(rfc2822))))
	e.EncodeToken(xml.EndElement{start.Name})
	return nil
}

type Owner struct {
	XMLName xml.Name `xml:"itunes:owner"`
	Name    string   `xml:"itunes:name"`
	Email   string   `xml:"itunes:email"`
}

type Image struct {
	XMLName xml.Name `xml:"itunes:image"`
	Href    string   `xml:"href,attr"`
}

type Category struct {
	XMLName    xml.Name `xml:"itunes:category"`
	Href       string   `xml:"text,attr"`
	Categories []*Category
}

type Enclosure struct {
	XMLName xml.Name `xml:"enclosure"`
	URL     string   `xml:"url,attr"`
	Length  string   `xml:"length,attr,omitempty"`
	Type    string   `xml:"type,attr"`
}

type Item struct {
	XMLName         xml.Name `xml:"item"`
	Title           string   `xml:"title"`
	PubDate         *PubDate `xml:"pubDate"`
	Author          string   `xml:"itunes:author,omitempty""`
	Block           string   `xml:"itunes:block,omitempty"`
	Image           *Image
	Duration        time.Duration `xml:"itunes:duration,omitempty""`
	Explicit        string        `xml:"itunes:explicit,omitempty"`
	ClosedCaptioned string        `xml:"itunes:isClosedCaptioned,omitempty"`
	Order           int           `xml:"itunes:order,omitempty"`
	Subtitle        string        `xml:"itunes:subtitle,omitempty""`
	Summary         string        `xml:"itunes:summary,omitempty""`
	GUID            string        `xml:"guid"`
	Enclosure       *Enclosure
}

type Channel struct {
	XMLName     xml.Name `xml:"channel"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Copyright   string   `xml:"copyright"`
	Language    string   `xml:"language"`
	Description string   `xml:"description"`
	Author      string   `xml:"itunes:author,omitempty""`
	Block       string   `xml:"itunes:block,omitempty"`
	Categories  []*Category
	Image       *Image
	Explicit    string `xml:"itunes:explicit,omitempty"`
	Complete    string `xml:"itunes:complete,omitempty"`
	NewFeedURL  string `xml:"itunes:new-feed-url,omitempty"`
	Owner       *Owner
	Subtitle    string `xml:"itunes:subtitle,omitempty""`
	Summary     string `xml:"itunes:summary,omitempty""`
	Items       []*Item
}

type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Xmlns   string   `xml:"xmlns:itunes,attr"`
	Version string   `xml:"version,attr"`
	Channel *Channel
}

func (f *Feed) Xml() (string, error) {
	data, err := xml.MarshalIndent(f, "", "  ")
	if err != nil {
		return "", err
	}
	s := xml.Header + string(data)
	return s, nil
}

func (f *Feed) Write(w io.Writer) error {
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return err
	}
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	return enc.Encode(f)
}
