package bizrecord

type SaverType uint

const (
	SaverType_Mongo SaverType = iota + 1
	SaverType_Mysql
)
