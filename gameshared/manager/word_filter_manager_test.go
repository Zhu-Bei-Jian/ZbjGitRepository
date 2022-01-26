package manager

import (
	"testing"
)

func TestWordFilterManager_FilterSync(t *testing.T) {
	var wordTests = []struct {
		input  string
		expect string
	}{
		{"SB只会吃", "*只会吃"},
		{"我是良民", "我是良民"},
		{"", ""},
	}

	url := "http://filterword.sanguosha.com:50000/v1/query?"
	wordFilterMgr := NewWordFilterManager(url)

	for _, v := range wordTests {
		text, err := wordFilterMgr.FilterSync(v.input)
		if err != nil {
			t.Fail()
		}
		if text != v.expect {
			t.Errorf("%v filter to %v,expected:%v", v.input, text, v.expect)
		}
	}
}

func TestWordFilterManager_CheckIsLegalSync(t *testing.T) {
	var wordTests = []struct {
		input  string
		expect bool
	}{
		{"SB只会吃", false},
		{"我是良民", true},
	}

	url := "http://filterword.sanguosha.com:50000/v1/query?"
	wordFilterMgr := NewWordFilterManager(url)

	for _, v := range wordTests {
		isLegal := wordFilterMgr.CheckIsLegalSync(v.input)
		if isLegal != v.expect {
			t.Errorf("%v isLegal %v,expected:%v", v.input, isLegal, v.expect)
		}
	}
}
