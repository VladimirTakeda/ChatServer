package types

import "time"

type ChatInfo struct {
	FromUser int `json:"from_user"`
	ToUser   int `json:"to_user"`
}

type DeleteChatInfo struct {
	ChatId int `json:"chat_id"`
}

type DeviceInfo struct {
	DeviceId   string `json:"device_id"`
	DeviceType string `json:"device_type"`
}

type RegisterUser struct {
	PassWord string `json:"password"`
	Nickname string `json:"login"`
	DeviceInfo
}

type GetUsersByPrefixReq struct {
	CurrUserId int    `json:"this_user_id"`
	Prefix     string `json:"search_prefix"`
}

type GetFileReq struct {
	FileId string `json:"file_id"`
}

type UserInfo struct {
	UserId   int    `json:"id"`
	Nickname string `json:"nickname"`
}

type UsersList struct {
	Users []UserInfo `json:"users"`
}

type ChatMetaInfo struct {
	ChatId int       `json:"chat_id"`
	Time   time.Time `json:"last_message_time"`
}

type WsMessage struct {
	Content     string   `json:"content"`
	Attachments []string `json:"attachments"`
	ChatName    string   `json:"chat_name"`
	UserFrom    int      `json:"user_from_id"`
	ChatTo      int      `json:"chat_to_id"`
}

type WsMessageOut struct {
	WsMessage
	Time time.Time `json:"time"`
}
