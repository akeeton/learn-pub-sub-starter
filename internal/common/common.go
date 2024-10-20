package common

import (
	"os"
	"os/signal"
)

func WaitForCtrlC() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh
}
