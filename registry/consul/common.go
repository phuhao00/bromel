package consul

func GetIntranetIP() string {
	//addrs, err := net.InterfaceAddrs()
	//if err != nil {
	//	return "127.0.0.1"
	//}
	//
	//for _, address := range addrs {
	//	// 检查ip地址判断是否回环地址
	//	if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
	//		if ipnet.IP.To4() != nil {
	//			return ipnet.IP.String()
	//		}
	//	}
	//}
	return "127.0.0.1"
}
