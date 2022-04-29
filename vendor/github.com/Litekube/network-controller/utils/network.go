package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func QueryPublicIp() string {
	resp, err := http.Get("http://ip.dhcp.cn/?ip") // 获取外网 IP
	if err != nil {
		logger.Errorf("fail to get public ip err:%+v", err)
		return ""
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	ip := fmt.Sprintf("%s", string(body))
	return ip
}
