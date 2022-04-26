package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/litekube/LiteKube/cmd/worker/app"
	"github.com/spf13/pflag"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Init for global klog
	klog.InitFlags(nil)
	defer klog.Flush()

	// add caller info
	klog.AddCallerName = true

	// Init Cobra command
	cmd := app.NewWorkerCommand()
	pflag.CommandLine.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	// Run LiteKube
	if err := cmd.Execute(); err != nil {
		year, month, day := time.Now().Date()
		panic(fmt.Sprintf("LiteKube worker exit at %d-%d-%d %d:%d:%d, error info: %s", year, month, day, time.Now().Hour(), time.Now().Minute(), time.Now().Second(), err.Error()))
	} else {
		year, month, day := time.Now().Date()
		klog.Infof("LiteKube worker goodby at %d-%d-%d %d:%d:%d", year, month, day, time.Now().Hour(), time.Now().Minute(), time.Now().Second())
	}

}

// func setDefaultLog() {
// 	klog.MaxSize = 10240
// 	if err := os.MkdirAll("litekube-logs/lite-apiserver", os.ModePerm); err != nil {
// 		panic(err)
// 	}

// 	flag.Set("logtostderr", "false")
// 	year, month, day := time.Now().Date()
// 	flag.Set("log_file", fmt.Sprintf("litekube-logs/lite-apiserver/log-%d-%d-%d_%d-%d.log", year, month, day, time.Now().Hour(), time.Now().Minute()))
// 	flag.Set("alsologtostderr", "true")
// }
