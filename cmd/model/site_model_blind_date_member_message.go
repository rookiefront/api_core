package model

import ()

type BlindDateMemberMessage struct {
	Model
	MessageType      string                         `json:"messageType" gorm:"column:message_type;comment:类型"`
	ParticipantUsers DataJSON                       `json:"participantUsers" gorm:"column:participant_users;type:longtext;comment:参与人"`
	MessageList      []BlindDateMemberMessageDetail `json:"messageList" gorm:"many2many:m2m_blind_date_member_message_blind_date_member_message_detail;comment:聊天信息"`
}

func (BlindDateMemberMessage) TableName() string {
	return "blind_date_member_message"
}
