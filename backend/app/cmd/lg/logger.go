package lg

import (
	"github.com/cartmanis/logger"
	"log"
)

func Error(i ...interface{}) {
	l, err := initLogger()
	if err != nil {
		return
	}
	defer l.Close()
	l.ErrorDepth(3, i...)
}

func Warn(i ...interface{}) {
	l, err := initLogger()
	if err != nil {
		return
	}
	defer l.Close()
	l.WarnDepth(3, i...)
}

func Info(i ...interface{}) {
	l, err := initLogger()
	if err != nil {
		return
	}
	defer l.Close()
	l.InfoDepth(3, i...)
}

func Errorf(mes string, i ...interface{}) {
	l, err := initLogger()
	if err != nil {
		return
	}
	defer l.Close()
	l.ErrorDepthf(3, mes, i...)
}

func Warnf(mes string, i ...interface{}) {
	l, err := initLogger()
	if err != nil {
		return
	}
	defer l.Close()
	l.WarnDepthf(3, mes, i...)
}

func Infof(mes string, i ...interface{}) {
	l, err := initLogger()
	if err != nil {
		return
	}
	defer l.Close()
	l.InfoDepthf(3, mes, i...)
}

func initLogger() (*logger.Logger, error) {
	l, err := logger.NewLogger("./english.log", true, true)
	if err != nil {
		log.Println("[ERROR] не удалось произвести инициализацию системы логирования. Ошибка:", err)
		return nil, err
	}
	return l, nil
}
