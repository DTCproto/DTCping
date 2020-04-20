package colo

//fl=4f197
//h=104.24.71.22
//ip=160.120.134.135
//ts=1584887600.227
//visit_scheme=http
//uag=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36
//colo=SJC
//http=http/1.1
//loc=CN
//tls=off
//sni=off
//warp=off

type HttpColoByte struct {
	Ip       string
	respBody []byte
	Error    error
}

type IpColo struct {
	Ip    string
	Colo  string
	Error error
}

const (
	defaultLimiterNumber = 256
	defaultBodyCache     = 1024
)
