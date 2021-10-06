package bizrecord

//RechargeRecord ...
type RechargeRecord struct {
	UserID    uint64
	Category  int
	Quantity  uint64
	Date      string
	SaverType SaverType
}

func (r RechargeRecord) Save(saverIF SaverIF) bool {
	return false
}

func (r RechargeRecord) Load() {

}

func (r RechargeRecord) GetSaverType() SaverType {
	return r.SaverType
}
