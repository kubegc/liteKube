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
		panic(fmt.Sprintf("LiteKube worker exit at %s, error info: %s", time.Now().Format("2006-01-02 15:04:05"), err.Error()))
	} else {
		klog.Infof("LiteKube worker goodby at %s", time.Now().Format("2006-01-02 15:04:05"))
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
