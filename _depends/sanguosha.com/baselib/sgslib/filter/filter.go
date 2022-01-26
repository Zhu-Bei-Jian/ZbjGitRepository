package filter

import (
	//"encoding/binary"
	"errors"
	//"fmt"
	//"net"
	//"reflect"
	//"time"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
	//"github.com/golang/protobuf/proto"
	//"sanguosha.com/baselib/appframe"
	"sanguosha.com/baselib/network/msgprocessor"
	"sanguosha.com/baselib/sgslib/filter/protocol"
)

type Poster interface {
	Post(func())
}

func init() {
	msgprocessor.RegisterMessage((*protocol.Client_FilterReq)(nil))
	msgprocessor.RegisterMessage((*protocol.Server_FilterResp)(nil))
}

// Filter 敏感词过滤器
type Filter interface {
	Check(app Poster, content string, cb func(error))
	Filter(app Poster, content string, cb func(string, error))
	Close()
}

// New ...
/*func New(addr string) Filter {
	const timeout = 5 * time.Second
	msg2id := func(msg proto.Message) (id uint32, exist bool) {
		return msgprocessor.MessageID(reflect.TypeOf(msg))
	}
	id2msg := func(id uint32) (msg proto.Message, exist bool) {
		typ, ok := msgprocessor.MessageType(id)
		if !ok {
			return nil, false
		}
		return reflect.New(typ.Elem()).Interface().(proto.Message), true
	}
	return &filter{
		pool: pio.NewPool(func() (pio.ProtoIO, error) {
			conn, err := net.DialTimeout("tcp", addr, timeout)
			if err != nil {
				return nil, err
			}
			return pio.New(conn, msg2id, id2msg, binary.BigEndian, nil), nil
		}, 32),
		timeout: timeout,
	}
}

type filter struct {
	pool    *pio.Pool
	timeout time.Duration
}

func (f *filter) Check(content string, ) ([]string, error) {
	_, keys, err := f.do(content)
	return keys, err
}

func (f *filter) Filter(content string) (string, error) {
	result, _, err := f.do(content)
	return result, err
}

func (f *filter) do(content string) (result string, keys []string, err error) {
	pio := f.pool.Get()
	defer pio.Close()

	err = pio.Write(&protocol.Client_FilterReq{Content: content})
	if err != nil {
		return "", nil, err
	}
	iresp, err := pio.ReadTimeout(f.timeout)
	if err != nil {
		return "", nil, err
	}

	resp, ok := iresp.(*protocol.Server_FilterResp)
	if !ok {
		return "", nil, errors.New("resp msg not protocol.Server_FilterResp")
	}
	if resp.Data == nil || len(resp.Data.Keyword) == 0 {
		return content, nil, nil
	}
	return resp.Data.NewText, resp.Data.Keyword, nil
}

func (f *filter) Close() {
	f.pool.Close()
}
*/

// HttpFilter
func NewHttpFilter(url string) Filter {
	return httpFilter{
		url : url,
	}
}

type httpFilter struct{
	url string
}

func (f httpFilter) Check(app Poster, content string, cb func(error)){
	go func(){
		text, _, _ := f.do(content)
		app.Post(func(){
			if text != content{
				cb(errors.New("检测到屏蔽词"))
			}else{
				cb(nil)
			}
		})
	}()
}

func (f httpFilter) Filter(app Poster, content string, cb func(string, error)) {
	go func(){
		text, _, err := f.do(content)
		app.Post(func(){
			cb(text, err)
		})
	}()
}

func (f httpFilter) Close() {
}

type filterResult struct {
	Code 		int 		`json:"code"`
	Keywords 	[]string	`json:"keywords"`
	Text 		string 		`json:"text"`
}

func (f *httpFilter) do(content string) (result string, keys []string, err error) {
		v:=url.Values{}
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

// NewNullFilter 创建一个无用的 Filter, 用于测试.
func NewNullFilter() Filter {
	return nullFilter{}
}

type nullFilter struct{}

func (f nullFilter) Check(app Poster, content string, cb func(error)) {
	cb(nil)
}
func (f nullFilter) Filter(app Poster, content string, cb func(string, error)){
	cb(content, nil)
}
func (f nullFilter) Close() {
}
