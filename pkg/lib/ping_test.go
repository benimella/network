package lib

import "testing"

func TestPing(t *testing.T) {
	host := "www.baidu.com"
	Ping(host)
}