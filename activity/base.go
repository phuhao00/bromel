package activity

import "github.com/phuhao00/bromel/trigger"

type Base struct {
	conf interface{}
}

func (b *Base) BuryLogBus(i interface{}) {
	//TODO implement me
	panic("implement me")
}

func (b *Base) Init() error {
	//TODO implement me
	panic("implement me")
}

func (b *Base) OnClientActive(param *ClientActiveParam) {
	//TODO implement me
	panic("implement me")
}

func (b *Base) OnLogin() {
	//TODO implement me
	panic("implement me")
}

func (b *Base) OnLogout() {
	//TODO implement me
	panic("implement me")
}

func (b *Base) Verify() {
	//TODO implement me
	panic("implement me")
}

func (b *Base) IsInTime() bool {
	//TODO implement me
	panic("implement me")
}

func (b *Base) OnNotify(event trigger.Event) {
	//TODO implement me
	panic("implement me")
}

func (b *Base) DailyRefresh() {
	//TODO implement me
	panic("implement me")
}

func (b *Base) Load() {
	//TODO implement me
	panic("implement me")
}

func (b *Base) Save() {
	//TODO implement me
	panic("implement me")
}
