package robot

type Robot struct {
	robotID uint64
	//gatewayClient
	actionList []func() bool // 机器人行为
	Name       string
	Pwd        string
	timer      interface{} // 机器人定时器
	conf       *Config
	isLogin    bool // 是否登录完成
	isLogining bool // 登录中...
	isOnline   bool // 是否登录完成
}

func newRobot(server string) *Robot {
	r := &Robot{}
	r.actionList = make([]func() bool, 0)
	r.registerAction()
	return r
}
