package models

// 请假记录表
type SysLeave struct {
	Model
	UserId          uint    `gorm:"comment:'用户编号'" json:"userId"`
	User            SysUser `gorm:"foreignkey:UserId" json:"user"`
	Status          *uint   `gorm:"default:0;comment:'状态(0:提交 1:批准 2:拒绝 3:取消 4:重启 5:结束)'" json:"status"`
	ApprovalOpinion string  `gorm:"comment:'审批意见'" json:"approvalOpinion"`
	Desc            string  `gorm:"comment:'说明'" json:"desc"`
}

func (m SysLeave) TableName() string {
	return m.Model.TableName("sys_leave")
}
