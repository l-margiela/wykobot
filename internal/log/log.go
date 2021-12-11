package log

import (
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// errJSON is a failover method of logging zap-like JSON with error.
//
// Useful if zap initialization fails.
func errJSON(err error) error {
	if err == nil {
		return nil
	}

	errJSON, err := json.Marshal(struct {
		Err string `json:"error"`
	}{err.Error()})
	if err != nil {
		return err
	}
	fmt.Println(string(errJSON))

	return nil
}

func NewZap(mode string, debug bool) *zap.Logger {
	var l *zap.Logger

	switch mode {
	case "prod":
		var err error
		logCfg := zap.NewProductionConfig()

		if debug {
			logCfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		}
		l, err = logCfg.Build()
		if err != nil {
			if err := errJSON(err); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		// https://github.com/uber-go/zap/issues/328
		//nolint
		defer l.Sync()
	case "dev":
		var err error
		l, err = zap.NewDevelopment()
		if err := errJSON(err); err != nil {
			panic(err)
		}

		// https://github.com/uber-go/zap/issues/328
		//nolint
		defer l.Sync()
	default:
		err := fmt.Errorf("invalid mode: %s", mode)
		if err := errJSON(err); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	return l
}
