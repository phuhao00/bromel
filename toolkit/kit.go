package toolkit

import (
	"os/user"
	"strings"
)

// GetUser 返回当前运行的用户
func GetUser() string {
	var userName string
	u, err := user.Current()
	if err != nil {
		userName = "unknow"
	} else {
		userName = u.Username
	}
	sl := strings.Split(userName, "\\")
	userName = sl[len(sl)-1]
	return userName
}
