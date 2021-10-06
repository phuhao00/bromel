package bizrecord

type RecordMgr struct {
	ChRecord chan RecordIF
	savers   map[SaverType]SaverIF
	ChStop   chan struct{}
}

func (m *RecordMgr) Run() {
	defer func() {
		if err := recover(); err != nil {
			//todo err
			go m.Run()
		}
	}()

	for {
		select {
		case record := <-m.ChRecord:
			succeed := m.Save(record)
			if succeed {
				//todo
			}
		case <-m.ChStop:
			m.Stop()
			return
		}
	}
}

//Stop ...
func (m *RecordMgr) Stop() {
	close(m.ChRecord)
	for record := range m.ChRecord {
		m.Save(record)
	}
	//todo log
}

//Save ...
func (m *RecordMgr) Save(record RecordIF) bool {
	if saver, exist := m.savers[record.GetSaverType()]; exist {
		return saver.Save(record)
	}
	return false
}

//AsyncSave ...
func (m *RecordMgr) AsyncSave(record RecordIF) {
	if saver, exist := m.savers[record.GetSaverType()]; exist {
		go saver.Save(record)
	}
}

//AddSaver add saver
func (m *RecordMgr) AddSaver(t SaverType, saver SaverIF) bool {
	if m.savers == nil {
		return false
	}
	if _, exist := m.savers[t]; exist {
		return false
	}
	m.savers[t] = saver
	return true
}

//DelSaver ...
func (m *RecordMgr) DelSaver(t SaverType) {
	delete(m.savers, t)
}
