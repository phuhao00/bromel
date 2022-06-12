package activity

import "sync"

type Mgr struct {
	Acts sync.Map
}

func (m *Mgr) OnLogin() {

}

func (m *Mgr) OnLogout() {

}
