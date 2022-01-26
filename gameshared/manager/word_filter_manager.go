package manager

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"sanguosha.com/baselib/util"
)

type WordFilterManager struct {
	f Filter
}

func NewWordFilterManager(url string) *WordFilterManager {
	var f Filter
	if url == "" || url == "null" {
		f = NewNullFilter()
	} else {
		f = newHttpFilter(url)
	}
	return &WordFilterManager{
		f: f,
	}
}

func (p *WordFilterManager) CheckIsLegalAsync(content string, cb func(bool)) {
	util.SafeGo(func() {
		isLegal := p.f.CheckIsLegal(content)
		cb(isLegal)
	})
}

func (p *WordFilterManager) FilterAsync(content string, cb func(string, error)) {
	util.SafeGo(func() {
		text, err := p.f.Filter(content)
		cb(text, err)
	})
}

func (p *WordFilterManager) CheckIsLegalSync(content string) bool {
	return p.f.CheckIsLegal(content)
}

func (p *WordFilterManager) FilterSync(content string) (string, error) {
	text, err := p.f.Filter(content)
	return text, err
}

func (p *WordFilterManager) Close() {
	p.f.Close()
}

type Filter interface {
	CheckIsLegal(content string) bool
	Filter(content string) (string, error)
	Close()
}

func newHttpFilter(url string) httpFilter {
	return httpFilter{
		url: url,
	}
}

type httpFilter struct {
	url string
}

func (f httpFilter) CheckIsLegal(content string) bool {
	text, _, _ := f.do(content)

	if text != content {
		return false
	}

	return true
}

func (f httpFilter) Filter(content string) (string, error) {
	text, _, err := f.do(content)
	return text, err
}

func (f httpFilter) Close() {
}

type filterResult struct {
	Code     int      `json:"code"`
	Keywords []string `json:"keywords"`
	Text     string   `json:"text"`
}

func (f *httpFilter) do(content string) (result string, keys []string, err error) {
	v := url.Values{}
	v.Add("q", content)
	filterUrl := f.url + v.Encode()

	resp, err := http.Get(filterUrl)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	res := &filterResult{}
	err = json.Unmarshal(body, res)
	return res.Text, res.Keywords, err
}

func NewNullFilter() Filter {
	return nullFilter{}
}

type nullFilter struct{}

func (f nullFilter) CheckIsLegal(content string) bool {
	return true
}

func (f nullFilter) Filter(content string) (string, error) {
	return content, nil
}

func (f nullFilter) Close() {
}
