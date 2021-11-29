package util

import (
	"fmt"
	"strings"
)

func GetBaseUrl(host string, production bool) string {
	var baseUrl string
	baseDomain := GetBaseDomain(host)
	if production {
		baseUrl = fmt.Sprintf("https://%v", baseDomain)
	} else {
		baseUrl = fmt.Sprintf("http://%v:3000", baseDomain)
	}
	return baseUrl
}

func GetBaseDomain(host string) string {
	addr, err := SplitHostPortDefault(host, "", "")
	if err != nil {
		return "byte.builders"
	}
	parts := strings.Split(addr.Host, ".")
	base := strings.Join(parts[1:], ".")
	return base
}
