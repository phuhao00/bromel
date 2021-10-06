package bizrecord

type RecordIF interface {
	Save(saverIF SaverIF) bool
	Load()
	GetSaverType() SaverType
}

type SaverIF interface {
	Save(recordIF RecordIF) bool
	Load(f func(saverIF SaverIF) RecordIF) RecordIF
}
