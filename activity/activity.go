package activity

import "github.com/phuhao00/bromel/trigger"

type If interface {
	Init() error
	IsInTime() bool
	OnNotify(event trigger.Event)
	DailyRefresh()
	Load()
	Save()
	Verify()
	OnLogin()
	OnLogout()
	OnClientActive(param *ClientActiveParam)
	BuryLogBus(interface{})
}
