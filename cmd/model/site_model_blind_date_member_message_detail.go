package model

import ()

type BlindDateMemberMessageDetail struct {
	Model
	SendUserId     int    `json:"sendUserId" gorm:"column:send_user_id;comment:发送人"`
	DetailType     string `json:"detailType" gorm:"column:detail_type;comment:信息类型"`
	ReceiverUserId int    `json:"receiverUserId" gorm:"column:receiver_user_id;comment:接收人"`
	Content        string `json:"content" gorm:"column:content;type:text;comment:内容"`
	MessageId      int    `json:"messageId" gorm:"column:message_id;comment:信息ID"`
}

func (BlindDateMemberMessageDetail) TableName() string {
	return "blind_date_member_message_detail"
}
