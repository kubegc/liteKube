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
	"os"

	"github.com/op/go-logging"
)

var logger *logging.Logger

func GetLogger() *logging.Logger {
	if logger == nil {
		InitLogger()
	}
	return logger
}

func InitLogger() {
	logger = logging.MustGetLogger("network-controller")
	fmt_string := "\r%{color}[%{time:06-01-02 15:04:05}][%{shortfile}][%{level:.6s}] %{shortfunc}%{color:reset} %{message}"
	format := logging.MustStringFormatter(fmt_string)
	logging.SetFormatter(format)
	logging.SetBackend(logging.NewLogBackend(os.Stdout, "", 0))
}

func SetLoggerLevel(debug bool) {
	if debug {
		logging.SetLevel(logging.DEBUG, "network-controller")
	} else {
		logging.SetLevel(logging.INFO, "network-controller")
	}
}
