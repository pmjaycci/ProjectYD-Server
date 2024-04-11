package response_packet

type HeartBeat struct {
	Code      uint   `json:"code"`
	Message   string `json:"message"`
	HeartBeat string `json:"heartBeat"`
}

type GameDB struct {
	Code               uint                  `json:"code"`
	Message            string                `json:"message"`
	ItemTable          map[int]Item          `json:"itemTable"`
	ItemWeaponTable    map[int]ItemWeapon    `json:"itemWeaponTable"`
	ItemEffectTable    map[int]ItemEffect    `json:"itemEffectTable"`
	ShopTable          map[int]ShopItem      `json:"shopTable"`
	WeaponEnchantTable map[int]WeaponEnchant `json:"weaponEnchantTable"`

	Money int `json:"money"`
}

type Item struct {
	Id        int    `json:"id"`
	ItemName  string `json:"itemName"`
	ItemType  int    `json:"itemType"`
	Category  int    `json:"category"`
	ImageName string `json:"imageName"`
	IsStack   bool   `json:"isStack"`
}
type ItemWeapon struct {
	Damage int     `json:"damage"`
	Speed  float64 `json:"speed"`
	Range  int     `json:"range"`
}
type ItemEffect struct {
	MaxHp       float64 `json:"maxHp"`
	RegenHp     float64 `json:"regenHp"`
	Speed       float64 `json:"shortDamage"`
	Damage      float64 `json:"damage"`
	AttackSpeed float64 `json:"attackSpeed"`
}
type WeaponEnchant struct {
	Enchant     int `json:"enchant"`
	Probability int `json:"probability"`
	Price       int `json:"price"`
}
type ShopItem struct {
	Id        int  `json:"id"`
	MoneyType int  `json:"moneyType"`
	Price     int  `json:"price"`
	IsBuy     bool `json:"isBuy"`
}

type InventoryItem struct {
	Id      int `json:"id"`
	Count   int `json:"count"`
	Enchant int `json:"enchant"`
}

type Inventory struct {
	Code    uint                  `json:"code"`
	Message string                `json:"message"`
	Items   map[int]InventoryItem `json:"items"`
}

type BuyItem struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
	Id      int    `json:"id"`
	Count   int    `json:"count"`
	Enchant int    `json:"enchant"`
	Money   int    `json:"money"`
}

type UpgradeItem struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
	Id      int    `json:"id"`
	Enchant int    `json:"enchant"`
	Money   int    `json:"money"`
}

type Weapon struct {
	Id      int `json:"id"`
	Enchant int `json:"enchant"`
}
type Effect struct {
	Id    int `json:"id"`
	Count int `json:"count"`
}
type JoinGame struct {
	Code         uint     `json:"code"`
	Message      string   `json:"message"`
	Gold         int      `json:"gold"`
	CurrentStage int      `json:"currentStage"`
	Slot         []Weapon `json:"slot"`
	Effect       []Effect `json:"effect"`
}

type IngameShopItem struct {
	Id    int `json:"id"`
	Price int `json:"price"`
}
type LoadIngameShop struct {
	Code    uint             `json:"code"`
	Message string           `json:"message"`
	Items   []IngameShopItem `json:"items"`
}

type BuyIngameItem struct {
	Code         uint     `json:"code"`
	Message      string   `json:"message"`
	Gold         int      `json:"gold"`
	CurrentStage int      `json:"currentStage"`
	Slot         []Weapon `json:"slot"`
	Effect       []Effect `json:"effect"`
}

type TimeAttackUser struct {
	Rank       int     `json:"rank"`
	UserName   string  `json:"userName"`
	RecordTime float64 `json:"recordTime"`
}

type LoadTimeAttackRankTable struct {
	Code     uint             `json:"code"`
	Message  string           `json:"message"`
	RankList []TimeAttackUser `json:"rankList"`
}

type UpdateTimeAttackRank struct {
	Code        uint    `json:"code"`
	Message     string  `json:"message"`
	RecordTime  float64 `json:"recordTime"`
	Rank        int     `json:"rank"`
	Money       int     `json:"money"`
	RewardMoney int     `json:"rewardMoney"`
}

type GameOver struct {
	Code        uint   `json:"code"`
	Message     string `json:"message"`
	Money       int    `json:"money"`
	RewardMoney int    `json:"rewardMoney"`
}
