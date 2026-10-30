package goroutine

import (
	"fmt"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/abulo/ratel/v3/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// try 在恢复 panic 的前提下执行 fn，cleaner 在退出时无条件执行；cleaner 为 nil 时跳过清理。
//
// panic 被转换为附带 fn 函数名的错误并作为返回值，非 error 类型的 panic 值会先格式化为错误。
func try(fn func() error, cleaner func()) (ret error) {
	if cleaner != nil {
		defer cleaner()
	}
	defer func() {
		if err := recover(); err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"err": err,
			}).Error("recover")
			if _, ok := err.(error); ok {
				ret = err.(error)
			} else {
				ret = fmt.Errorf("%+v", err)
			}
			ret = errors.Wrap(ret, util.FunctionName(fn))
		}
	}()
	return fn()
}

// try2 在恢复 panic 的前提下执行无返回值的 fn，cleaner 在退出时无条件执行；cleaner 为 nil
// 时跳过清理。panic 被记录日志后转换为错误返回，非 error 类型的 panic 值会先格式化。
func try2(fn func(), cleaner func()) (ret error) {
	if cleaner != nil {
		defer cleaner()
	}
	defer func() {
		if err := recover(); err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"err": err,
			}).Error("recover")
			if _, ok := err.(error); ok {
				ret = err.(error)
			} else {
				ret = fmt.Errorf("%+v", err)
			}
		}
	}()
	fn()
	return nil
}
