package types

type ChatInfo struct {
	FromUser int `json:"from_user"`
	ToUser   int `json:"to_user"`
}

type DeleteChatInfo struct {
	ChatId int `json:"chat_id"`
}

type RegisterUser struct {
	PassWord string `json:"password"`
	Nickname string `json:"login"`
}

type GetUsersByPrefixReq struct {
	CurrUserId int    `json:"this_user_id"`
	Prefix     string `json:"search_prefix"`
}

type UserInfo struct {
	UserId   int    `json:"id"`
	Nickname string `json:"nickname"`
}

type UsersList struct {
	Users []UserInfo `json:"users"`
}

type CompleteInfoJson struct {
	CompleteTime string `json:"complete_time"`
	CourierId    int    `json:"courier_id"`
	OrderId      int    `json:"order_id"`
}
type CompleteInfosJson struct {
	InfoList []CompleteInfoJson `json:"complete_info"`
}
