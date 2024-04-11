package request_packet

type HeartBeat struct {
	HeartBeat string `json:"heartBeat"`
}

type DuplicateLoginFromNotificationServer struct {
	UUID string `json:"uuid"`
}

type BuyItem struct {
	Id int `json:"id"`
}

type UpgradeItem struct {
	Id int `json:"id"`
}

type JoinGame struct {
	ItemId int `json:"itemId"`
}

type LoadIngameShop struct {
	Gold         int `json:"gold"`
	CurrentStage int `json:"currentStage"`
}

type BuyIngameItem struct {
	CurrentStage int `json:"currentStage"`
	ItemId       int `json:"itemId"`
}

type ChnageUserName struct {
	UserName string `json:"userName"`
}

type UpdateTimeAttackRank struct {
	RecordTime float64 `json:"recordTime"`
}
