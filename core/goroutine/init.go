package goroutine

import (
	"fmt"

	"github.com/abulo/ratel/v2/core/logger"
	"github.com/abulo/ratel/v2/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

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
			ret = errors.Wrap(ret, fmt.Sprintf("%s", util.FunctionName(fn)))
		}
	}()
	return fn()
}

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
