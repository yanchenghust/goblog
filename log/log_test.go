package log_test

import (
	"github.com/yanchenghust/goblog/log"
	"testing"
)

func TestLog(t *testing.T) {
	fp := log.NewFileProvider("test")
	lg := log.NewLog(log.LvlDebug, fp)
	lg.Start()
	for i := 0; i < 10; i++ {
		lg.Debugf("test%d", 1)
	}
}
