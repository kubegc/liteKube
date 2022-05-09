/*
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * Author: wanna <wananzjx@163.com>
 *
 */
package utils

import (
	"github.com/Litekube/network-controller/contant"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/op/go-logging"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var logger *logging.Logger

func GetLogger() *logging.Logger {
	return logger
}

func InitLogger(logDir, logName string, debug bool) {
	if logName == "" {
		logName = contant.DefualtLogName
	}

	// set rotate log
	latestPath := filepath.Join(logDir, logName)
	prefix := strings.Split(latestPath, ".log")[0]

	content, _ := rotatelogs.New(
		// retate log format
		prefix+"_%Y-%m-%d.log",
		// ref to latest log file
		rotatelogs.WithLinkName(latestPath),
		//MaxAge and RotationCount cannot be both set
		rotatelogs.WithMaxAge(time.Duration(168)*time.Hour),
		//rotate each day
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)

	// set logging format
	logger = logging.MustGetLogger("network-controller")
	loggerFmt := "\r%{color}[%{time:06-01-02 15:04:05}][%{shortfile}][%{level:.6s}] %{shortfunc}%{color:reset} %{message}"
	format := logging.MustStringFormatter(loggerFmt)
	logging.SetFormatter(format)

	// set output: stdout & file
	fback := logging.NewLogBackend(content, "", 0)
	stdback := logging.NewLogBackend(os.Stdout, "", 0)
	logging.SetBackend(fback, stdback)

	// set log level
	SetLoggerLevel(debug)
}

func SetLoggerLevel(debug bool) {
	if debug {
		logging.SetLevel(logging.DEBUG, "network-controller")
	} else {
		logging.SetLevel(logging.INFO, "network-controller")
	}
}
