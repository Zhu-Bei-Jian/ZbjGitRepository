package log

import "github.com/sirupsen/logrus"

// 此文件内为将废弃的辅助函数, 不建议再使用.

// Emergency 已废弃
func Emergency(format string, params ...interface{}) { logrus.Errorf(format, params...) }

// Alert 已废弃
func Alert(format string, params ...interface{}) { logrus.Errorf(format, params...) }

// Critical 已废弃
func Critical(format string, params ...interface{}) { logrus.Errorf(format, params...) }

// Error 已废弃
func Error(format string, params ...interface{}) { logrus.Errorf(format, params...) }

// Warning 已废弃
func Warning(format string, params ...interface{}) { logrus.Warningf(format, params...) }



// Debug 已废弃
func Debug(format string, params ...interface{}) { logrus.Debugf(format, params...) }
