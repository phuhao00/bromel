package bizrecord

//ItemRecord 道具记录
type ItemRecord struct {
	Count         uint32
	ID            uint64
	Date          string
	OperationType string
	ServerTag     string
	SaverType     SaverType
}

//Save 存储
func (r *ItemRecord) Save(Saver SaverIF) bool {
	return Saver.Save(r)
}

//Load 加载
func (r *ItemRecord) Load() {

}

func (r *ItemRecord) GetSaverType() SaverType {
	return r.SaverType
}
