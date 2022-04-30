package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/litekube/LiteKube/cmd/leader/app"
	"github.com/spf13/pflag"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Init for global klog
	klog.InitFlags(nil)
	defer klog.Flush()

	// Init Cobra command
	cmd := app.NewLeaderCommand()
	pflag.CommandLine.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	// Run LiteKube
	if err := cmd.Execute(); err != nil {
		panic(fmt.Sprintf("LiteKube leader exit at %s, error info: %s", time.Now().Format("2006-01-02 15:04:05"), err.Error()))
	} else {
		klog.Infof("LiteKube leader goodby at %s", time.Now().Format("2006-01-02 15:04:05"))
	}

}
