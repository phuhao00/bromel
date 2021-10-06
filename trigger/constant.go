package trigger

type TriggerType uint32

const (
	TriggerType_Unit       TriggerType = iota + 1 // 触发的玩家指针
	TriggerType_Skill                             // 技能ID
	TriggerType_Target                            // 目标单位ID
	TriggerType_MainAction                        // 主要动作事件
)
