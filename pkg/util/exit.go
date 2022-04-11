package util

import (
	"os"

	"k8s.io/klog/v2"
)

// exit the main process elegantly
func Exit(code int) {
	klog.Info("----------------------------------------------")
	klog.Infof("| ==> process exits voluntarily with code=%d", code)
	klog.Info("----------------------------------------------")
	shutdownHandler <- os.Interrupt
}
