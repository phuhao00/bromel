package robot

func (r *Robot) registerAction() {
	r.actionList = append(r.actionList, r.Chat)
}
