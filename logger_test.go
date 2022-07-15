package xxxlogger

import (
	"go.uber.org/zap"
	"testing"
)

type User struct {
	Name string
}

func TestMain(m *testing.M) {
	InitDevLogger("./tmp/dev.log")

	user := &User{Name: "lzc"}

	Info("Info test log", zap.Any("user", user))
	Debug("Debug test log", zap.Any("user", user))


	InitProdLogger("./tmp/prod.log")

	Info("Info test log", zap.Any("user", user))
	Debug("Debug test log", zap.Any("user", user))
}
