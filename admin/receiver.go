package admin

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"net/http"
)

// 命令接收者
type commandCallback interface {
	// 返回
	writeMsg(interface{}, error)
}

// http接收者
type receiver struct {
	res map[string]interface{}
	w   http.ResponseWriter
	r   *http.Request
}

func (p *receiver) write() {
	//
	//jsonpb.Marshaler{
	//	EnumsAsInts:  true,
	//	EmitDefaults: false,
	//	Indent:       "",
	//	OrigName:     true,
	//	AnyResolver:  nil,
	//}

	result, _ := json.Marshal(p.res)
	logrus.WithFields(logrus.Fields{}).Debug("return ", string(result))
	p.w.Write([]byte(result))
}

func (p *receiver) writeMsg(data interface{}, err error) {
	p.res = make(map[string]interface{})
	if err != nil {
		p.res[resCode] = 1
		p.res[resMsg] = err.Error()
		p.res[resItems] = data
		p.write()
		return
	}
	p.res[resMsg] = "成功"
	switch data.(type) {
	case proto.Message:
		p.res[resMsg] = "成功:"
	case string:
		p.res[resMsg] = "成功:" + data.(string)
	}

	p.res[resCode] = codeSuccess
	p.res["data"] = data
	p.write()
}
