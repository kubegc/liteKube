package verflag

import (
	"fmt"
	"os"
	"strconv"

	"github.com/litekube/LiteKube/pkg/version"
	flag "github.com/spf13/pflag"
)

type versionValue int

const (
	VersionFalse  versionValue = 0
	VersionTrue   versionValue = 1
	VersionRaw    versionValue = 2
	VersionSimple versionValue = 3
)

const strRawVersion string = "raw"
const strSimpleVersion string = "simple"

func (v *versionValue) IsBoolFlag() bool {
	return true
}

func (v *versionValue) Get() interface{} {
	return versionValue(*v)
}

func (v *versionValue) Set(s string) error {
	if s == strRawVersion {
		*v = VersionRaw
		return nil
	} else if s == strSimpleVersion {
		*v = VersionSimple
		return nil
	}

	boolVal, err := strconv.ParseBool(s)
	if boolVal {
		*v = VersionTrue
	} else {
		*v = VersionFalse
	}
	return err
}

func (v *versionValue) String() string {
	if *v == VersionRaw {
		return strRawVersion
	} else if *v == VersionSimple {
		return strSimpleVersion
	}
	return fmt.Sprintf("%v", bool(*v == VersionTrue))
}

// The type of the flag as required by the pflag.Value interface
func (v *versionValue) Type() string {
	return "versions"
}

func VersionVar(p *versionValue, name string, value versionValue, usage string) {
	*p = value
	if f := flag.Lookup(name); f != nil {
		f = &flag.Flag{
			Name:      name,
			Shorthand: "",
			Usage:     usage,
			Value:     p,
			DefValue:  value.String(),
		}
	} else {
		flag.Var(p, name, usage)
	}

	// "--version" will be treated as "--version=simple"
	flag.Lookup(name).NoOptDefVal = strSimpleVersion
}

func Version(name string, value versionValue, usage string) *versionValue {
	p := new(versionValue)
	VersionVar(p, name, value, usage)
	return p
}

const versionFlagName = "versions"

var (
	versionFlag = Version(versionFlagName, VersionFalse, "Print version information and quit, true/false/raw/simple (default: simple)")
)

// AddFlags registers this package's flags on arbitrary FlagSets, such that they point to the same value as the global flags.
func AddFlagsTo(fs *flag.FlagSet) {
	fs.AddFlag(flag.Lookup(versionFlagName)) // fs.IntVar(&versionFlag, "version", 0, "Print version information and quit")  //
}

// PrintAndExitIfRequested will check if the -version flag was passed and, if so, print the version and exit.
func PrintAndExitIfRequested() {
	if *versionFlag == VersionRaw {
		fmt.Printf("%#v\n", version.Get())
		os.Exit(0)
	} else if *versionFlag == VersionTrue {
		fmt.Printf("LiteKube %+v\n", version.Get())
		os.Exit(0)
	} else if *versionFlag == VersionSimple {
		fmt.Printf("%s\n", version.GetSimple())
		os.Exit(0)
	} else {
		// *versionFlag=="false"
		return
	}
}
