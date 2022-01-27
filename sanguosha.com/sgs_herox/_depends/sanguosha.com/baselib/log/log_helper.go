package log

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"sanguosha.com/baselib/log/asynchook"
	"sanguosha.com/baselib/log/iohook"
	"sanguosha.com/baselib/log/logstash"
)

// NewFileLogHook 异步记录本地文件日志插件 for logrus
func NewFileLogHook(dir string, filename string, useJSONFormat bool, rotate bool, maxSize int) (hook logrus.Hook, close func(), err error) {
	os.Mkdir(dir, os.ModePerm)

	dir, err = filepath.Abs(dir)
	if err != nil {
		return nil, nil, err
	}
	// Abs 会调用 Clean 方法, 因此会去除dir结尾的“/”
	dir += "/"

	var f io.WriteCloser

	if rotate {
		if maxSize > 0 {
			f = &lumberjack.Logger{
				Filename:   dir + filename + ".log",
				MaxSize:    maxSize,
				MaxBackups: 0,
				MaxAge:     0,
				Compress:   false,
				LocalTime:  true,
			}
		} else {
			f, err = rotatelogs.New(
				dir+filename+".%Y%m%d.log",
				rotatelogs.WithLinkName(dir+filename+".log"),
				rotatelogs.WithMaxAge(15*24*time.Hour),
				rotatelogs.WithRotationTime(24*time.Hour),
			)
			if err != nil {
				return nil, nil, err
			}
		}
	} else {
		f, err = os.OpenFile(dir+filename+".log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
		if err != nil {
			return nil, nil, err
		}
	}

	if useJSONFormat {
		hook = iohook.New(f, new(logrus.JSONFormatter))
	} else {
		hook = iohook.New(f, new(logrus.TextFormatter))
	}

	asyncHook := asynchook.NewWithHook(4096, hook)
	hook = asyncHook

	close = func() {
		asyncHook.Close()
		f.Close()
	}

	return hook, close, nil
}

// NewLogstashHook 异步输出 json 格式的日志到 logstash
func NewLogstashHook(addr string, typ string) (hook logrus.Hook, close func(), err error) {
	typ = strings.ToLower(typ)

	logstashhook, err := logstash.New(addr, logrus.Fields{"type": typ})
	if err != nil {
		return nil, nil, err
	}

	asyncHook := asynchook.NewWithHook(4096, logstashhook)
	hook = asyncHook

	close = func() {
		asyncHook.Close()
		logstashhook.Close()
	}

	return hook, close, nil
}
