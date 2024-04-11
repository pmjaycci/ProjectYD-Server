package game

import (
	"context"
	"database/sql"
	"encoding/json"
	"math/rand"
	"project_yd/server"
	packet "project_yd/server/server_packet"
	request "project_yd/server/server_packet/request_packet"
	response "project_yd/server/server_packet/response_packet"
	"project_yd/table"
	"project_yd/util"
	"time"
)

func RegistGameRpc() {
	server.RegistRpc("load_tables", LoadTables)
	server.RegistRpc("load_inventory", LoadInventory)
	server.RegistRpc("buy_item", BuyItem)
	server.RegistRpc("upgrade_item", UpgradeItem)
	server.RegistRpc("join_game", JoinGame)
	server.RegistRpc("load_ingame_shop", LoadIngameShop)
	server.RegistRpc("buy_ingame_item", BuyIngameItem)
	server.RegistRpc("user_name", ChangeUserName)
	server.RegistRpc("load_time_attack_rank_table", LoadTimeAttackRankTable)
	server.RegistRpc("update_time_attack_rank", UpdateTimeAttackRank)
	server.RegistRpc("game_over", GameOver)

	server.RegistRpc("user_name", ChangeUserName)
}

/*
서버에서 캐싱하고 있는 DB테이블들의 정보를 클라이언트에게 Return
*/
func LoadTables(UUID string, payload string) string {
	ctx := context.Background()
	responsePacket := response.GameDB{}

	item := response.Item{}
	responsePacket.ItemTable = make(map[int]response.Item)
	for key, val := range table.ItemTable {
		item.Id = val.Id
		item.ItemName = val.Name
		item.ItemType = val.Type
		item.Category = val.Category
		item.ImageName = val.ImageName
		item.IsStack = val.IsStack
		responsePacket.ItemTable[key] = item
	}

	itemWeapon := response.ItemWeapon{}
	responsePacket.ItemWeaponTable = make(map[int]response.ItemWeapon)
	for key, val := range table.ItemWeaponTable {
		itemWeapon.Damage = val.Damage
		itemWeapon.Speed = val.Speed
		itemWeapon.Range = val.Range
		responsePacket.ItemWeaponTable[key] = itemWeapon
	}
	itemEffect := response.ItemEffect{}
	responsePacket.ItemEffectTable = make(map[int]response.ItemEffect)
	for key, val := range table.ItemEffectTable {
		itemEffect.MaxHp = val.MaxHp
		itemEffect.RegenHp = val.RegenHp
		itemEffect.Speed = val.Speed
		itemEffect.Damage = val.Damage
		itemEffect.AttackSpeed = val.AttackSpeed
		responsePacket.ItemEffectTable[key] = itemEffect
	}
	weaponEnchant := response.WeaponEnchant{}
	responsePacket.WeaponEnchantTable = make(map[int]response.WeaponEnchant)
	for key, val := range table.WeaponEnchantProbabilityTable {
		weaponEnchant.Enchant = val.Enchant
		weaponEnchant.Probability = val.Probability
		weaponEnchant.Price = val.Price
		responsePacket.WeaponEnchantTable[key] = weaponEnchant
	}

	gameDB := server.DBManager.Game
	invenQuery := `SELECT item_id FROM inventory WHERE uid = ?`
	result, invenErr := gameDB.QueryContext(ctx, invenQuery, UUID)

	if invenErr != nil {
		return util.ResponseErrorMessage(util.ServerError, invenErr.Error())
	}
	inventory := map[int]bool{}
	inventory = make(map[int]bool)
	var itemId int
	for result.Next() {
		result.Scan(&itemId)
		inventory[itemId] = true
	}

	shopItem := response.ShopItem{}
	responsePacket.ShopTable = make(map[int]response.ShopItem)
	for key, val := range table.ShopTable {
		shopItem.Id = val.Id
		shopItem.MoneyType = val.MoneyType
		shopItem.Price = val.Price
		isBuy, isExist := inventory[val.Id]
		if isExist {
			shopItem.IsBuy = isBuy
		} else {
			shopItem.IsBuy = false
		}
		responsePacket.ShopTable[key] = shopItem
	}

	db := server.DBManager.Login
	moneyQuery := `SELECT money FROM account WHERE uid =?`
	err := db.QueryRowContext(ctx, moneyQuery, UUID).Scan(&responsePacket.Money)
	if err != nil {
		return util.ResponseErrorMessage(util.ServerError, err.Error())
	}

	responsePacket.Code = util.Success
	responsePacket.Message = "Success"

	return util.ResponseMessage(responsePacket)
}

/*
요청자의 UUID를 값으로 GameDB Inventory테이블에서 데이터를 불러온뒤 요청자에게 Return
*/
func LoadInventory(UUID string, payload string) string {

	responsePacket := response.Inventory{}

	ctx := context.Background()
	db := server.DBManager.Game

	query := `SELECT item_id, item_count, enchant_level FROM inventory WHERE uid = ?`
	result, err := db.QueryContext(ctx, query, UUID)

	item := response.InventoryItem{}
	items := map[int]response.InventoryItem{}
	items = make(map[int]response.InventoryItem)
	for result.Next() {
		if err == sql.ErrNoRows {
			responsePacket.Code = util.NotFound
			responsePacket.Message = "Have Not Items"
			return util.ResponseMessage(responsePacket)
		} else if err != nil {
			responsePacket.Code = util.ServerError
			responsePacket.Message = err.Error()
			return util.ResponseMessage(responsePacket)
		}

		result.Scan(&item.Id, &item.Count, &item.Enchant)
		items[item.Id] = item
	}

	responsePacket.Items = items
	responsePacket.Code = util.Success
	responsePacket.Message = "Success"

	return util.ResponseMessage(responsePacket)
}

/*
요청자는 구매할 아이템의 Id값을 요청하고
아이템 Id를 키값으로 캐싱되어 있는 Shop테이블에서 판매 가격을 불러오고
LoginDB account테이블에서 요청자의 재화 보유값과 판매 가격을 비교하여 구매 가능여부 판별후
구매 가능시 account테이블에서 Shop테이블 판매 가격만큼 차감후
요청자에게 return
*/
func BuyItem(UUID string, payload string) string {
	requestPacket := request.BuyItem{}
	responsePacket := response.BuyItem{}
	err := json.Unmarshal([]byte(payload), &requestPacket)
	if err != nil {
		return util.ResponseErrorMessage(util.BadRequest, err.Error())
	}

	ctx := context.Background()
	gameDB := server.DBManager.Game

	shopItem, exists := table.ShopTable[requestPacket.Id]
	if !exists {
		return util.ResponseErrorMessage(util.BadRequest, "Not Sell Item Id")
	}
	item, exists := table.ItemTable[requestPacket.Id]
	if !exists {
		return util.ResponseErrorMessage(util.BadRequest, "Have Not Item Id")
	}

	isUpdateMoney, updateMoney := CheckMoney(UUID, shopItem.Price)

	if !isUpdateMoney {
		return util.ResponseErrorMessage(util.BadRequest, "잔액 부족")
	}

	tx, err := gameDB.Begin()
	if err != nil {
		println("BuyItem Transaction QueryRow Error!!")
		println(err.Error())
		MoneyRollback(UUID, shopItem.Price)
		tx.Rollback()
		return util.ResponseErrorMessage(util.ServerError, err.Error())
	}

	var query string
	//-- 갯수가 겹치는 아이템인지
	if item.IsStack {
		//-- 이미 인벤토리에 있는 아이템의 경우 Count +1 아닐경우 Insert
		query = `INSERT INTO inventory (uid, item_id, item_count, enchant_level)
		VALUES (?, ?, 1, 0)
		ON DUPLICATE KEY UPDATE item_count = item_count + 1`
	} else {
		query = `INSERT INTO inventory (uid, item_id, item_count, enchant_level) VALUES (?, ?, 1, 0)`
	}

	_, err = tx.ExecContext(ctx, query, UUID, requestPacket.Id)
	if err != nil {
		println("BuyItem Transaction Exec Error!!")
		println(err.Error())
		MoneyRollback(UUID, shopItem.Price)
		tx.Rollback()
		return util.ResponseErrorMessage(util.ServerError, err.Error())
	}
	var count int
	var enchant int
	query = `SELECT item_count, enchant_level FROM inventory WHERE uid = ? AND item_id = ?`
	err = tx.QueryRowContext(ctx, query, UUID, requestPacket.Id).Scan(&count, &enchant)
	if err != nil {
		println("BuyItem Transaction QueryRow Error!!")
		println(err.Error())
		MoneyRollback(UUID, shopItem.Price)
		tx.Rollback()
		return util.ResponseErrorMessage(util.ServerError, err.Error())
	}

	err = tx.Commit()
	if err != nil {
		println("BuyItem Transaction Commit Error!!")
		println(err.Error())
		MoneyRollback(UUID, shopItem.Price)
		tx.Rollback()
		return util.ResponseErrorMessage(util.ServerError, err.Error())
	}

	responsePacket.Code = util.Success
	responsePacket.Message = "Success"

	responsePacket.Id = requestPacket.Id
	responsePacket.Count = count
	responsePacket.Enchant = enchant
	responsePacket.Money = updateMoney

	return util.ResponseMessage(responsePacket)
}

/*
요청자는 업그레이드할 아이템의 Id값을 요청하고
GameDB Inventory테이블에서 해당 아이템의 Enchant수치를 불러오고
LoginDB WeaponEnchantProbability테이블에서 인챈트 수치에 따른 확률정보를 불러온다.
확률은 소수점 4자리까지 표기하기 위해 0~100만까지의 랜덤값과 인챈트별 해당 확률정보 값을 비교하여
성공 여부를 판별하였고 성공 유무 상관없이 재화 또한 LoginDB account테이블에서 재화 감소후
업그레이드 결과 여부를 Return
*/
func UpgradeItem(UUID string, payload string) string {
	println("UpgradeItem")
	requestPacket := request.UpgradeItem{}
	responsePacket := response.UpgradeItem{}
	err := json.Unmarshal([]byte(payload), &requestPacket)
	if err != nil {
		return util.ResponseErrorMessage(util.BadRequest, err.Error())
	}

	ctx := context.Background()
	gameDB := server.DBManager.Game

	currentEnchant, err := GetWeaponEnchant(UUID, requestPacket.Id)
	if err == sql.ErrNoRows {
		return util.ResponseErrorMessage(util.BadRequest, "You Have Not This Item")
	} else if err != nil {
		return util.ResponseErrorMessage(util.BadRequest, err.Error())
	}

	maxEnchant := len(table.WeaponEnchantProbabilityTable)
	println("MaxEnchant:", maxEnchant, "/ CurrentEnchant:", currentEnchant)
	//-- 이미 현재 인챈트 상태가 최대치일경우
	if currentEnchant >= maxEnchant {
		var money int
		ctx := context.Background()
		loginDB := server.DBManager.Login
		query := `SELECT money FROM account WHERE uid = ?`
		err := loginDB.QueryRowContext(ctx, query, UUID).Scan(&money)
		if err != nil {
			println("UpgradeItem MaxEnchant QueryRow Error!!")
			println(err.Error())
			return util.ResponseErrorMessage(util.ServerError, err.Error())
		}

		responsePacket.Code = util.BadRequest
		responsePacket.Message = "Already Max Enchant"
		responsePacket.Id = requestPacket.Id
		responsePacket.Enchant = currentEnchant
		responsePacket.Money = money
		return util.ResponseMessage(responsePacket)
	}

	upgradeEnchantData, exists := table.WeaponEnchantProbabilityTable[currentEnchant+1]
	if !exists {
		return util.ResponseErrorMessage(util.BadRequest, "Have Not WeaponEnchantProbabilityTable Key")
	}
	println("EnchantData:: enchant:", upgradeEnchantData.Enchant, "/ price:", upgradeEnchantData.Price, "probability:", upgradeEnchantData.Probability)
	isUpdateMoney, updateMoney := CheckMoney(UUID, upgradeEnchantData.Price)
	if !isUpdateMoney {
		return util.ResponseErrorMessage(util.BadRequest, "잔액 부족")
	}

	// Seed 설정 (시드값 설정)
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	rand := rand.New(source)

	// 0 이상 100만 미만의 랜덤 정수 생성
	random := rand.Intn(1000000)
	println("성공확률:", upgradeEnchantData.Probability, "/ Random값:", random)
	//-- 업그레이드 성공했을경우
	if random <= upgradeEnchantData.Probability {
		query := `UPDATE inventory SET enchant_level = enchant_level + 1 WHERE uid = ? AND item_id = ?`
		_, err := gameDB.ExecContext(ctx, query, UUID, requestPacket.Id)
		if err != nil {
			println("UpgradeItem Error!! 244 Line")
			println(err.Error())
			return util.ResponseErrorMessage(util.ServerError, err.Error())
		}
		currentEnchant += 1
		responsePacket.Code = util.Success
		responsePacket.Message = "Success"
	} else {
		responsePacket.Code = util.Fail
		responsePacket.Message = "Fail"
	}
	responsePacket.Id = requestPacket.Id
	responsePacket.Enchant = currentEnchant
	responsePacket.Money = updateMoney

	return util.ResponseMessage(responsePacket)
}

/*
요청자는 인게임 플레이를 위해 착용할 장비 아이템의 Id값을 요청하고
GameDB Inventory테이블에서 해당 아이템의 정보와
요청자의 정보 (아이템 슬롯, 골드, 스테이지)를 초기화 및 캐싱하고
캐싱된 정보를 요청자에게 Return
*/
func JoinGame(UUID string, payload string) string {
	requestPacket := request.JoinGame{}
	responsePacket := response.JoinGame{}
	err := json.Unmarshal([]byte(payload), &requestPacket)
	if err != nil {
		return util.ResponseErrorMessage(util.BadRequest, err.Error())
	}

	weaponEnchant, err := GetWeaponEnchant(UUID, requestPacket.ItemId)
	if err != nil {
		return util.ResponseErrorMessage(util.ServerError, err.Error())
	}

	//-- 인게임 유저 정보 초기화
	_, exists := server.Users[UUID]
	if !exists {
		userData := server.User{}
		server.Users[UUID] = &userData
	}

	server.Users[UUID].Slot = []server.Weapon{}
	server.Users[UUID].Effect = []server.Effect{}
	server.Users[UUID].Gold = 0
	server.Users[UUID].CurrentStage = 1

	slot := server.Weapon{}
	slot.Id = requestPacket.ItemId
	slot.Enchant = weaponEnchant
	server.Users[UUID].Slot = append(server.Users[UUID].Slot, slot)

	//-- Response값 등록
	for _, val := range server.Users[UUID].Slot {
		slotData := response.Weapon{}
		slotData.Id = val.Id
		slotData.Enchant = val.Enchant
		responsePacket.Slot = append(responsePacket.Slot, slotData)
	}
	for _, val := range server.Users[UUID].Effect {
		effectData := response.Effect{}
		effectData.Id = val.Id
		effectData.Count = val.Count
		responsePacket.Effect = append(responsePacket.Effect, effectData)
	}
	responsePacket.Gold = server.Users[UUID].Gold
	responsePacket.CurrentStage = server.Users[UUID].CurrentStage

	_, isExist := server.Users[UUID]
	if !isExist {
		responsePacket.Code = util.Fail
		responsePacket.Message = "User Data Caching Fail"
		return util.ResponseMessage(responsePacket)
	}
	responsePacket.Code = util.Success
	responsePacket.Message = "Success"

	return util.ResponseMessage(responsePacket)
}

/*
요청자는 플레이한 스테이지와 획득한 골드값을 호출하고
플레이한 스테이지와 서버에서 캐싱하고 있는 스테이지 값과 비교하여
해당 유저가 스테이지 조작여부를 판별하고
획득한 골드의 재화값만큼 캐싱하고 있는 골드 값에 추가를 한다.

요청자에게 보여줄 판매할 인게임 아이템들은 랜덤으로 획득을 하게 된다.
RandomItemId함수를 통해 랜덤하게 뽑힌 아이템 리스트를 Return한다.

RandomItemId의 구조는
해당 함수 첫 실행시 매개변수인 ItemIds가 가지게 되는 데이터는
기본적으로 LoginDB ItemWeapon테이블, ItemEffect테이블에 등록된 ItemId들이며
요청자가 착용한 장비의 갯수가 2개 이상일시 ItemWeapon테이블에 등록된 ItemId들 대신
요청자가 착용한 장비의 ItemId들이 추가 된다.
함수가 실행되면서 0부터 매개변수로 받은 ItemIds의 크기의 범위까지 랜덤한 값을 뽑고
해당 값의 ItemIds의 인덱스를 가진 데이터를 비워가며 재귀함수로 연산하게 되며
ItemIds가 판매 아이템 목록 갯수만큼의 사이즈에 도달할 경우 결과를 Return
*/
func LoadIngameShop(UUID string, payload string) string {
	requestPacket := request.LoadIngameShop{}
	responsePacket := response.LoadIngameShop{}
	err := json.Unmarshal([]byte(payload), &requestPacket)
	if err != nil {
		return util.ResponseErrorMessage(util.BadRequest, err.Error())
	}

	//-- 현재 유저의 인게임 스테이지와 유저가 호출한 스테이지 값 비교 처리
	//-- 스테이지 값이 다르다면 클라이언트에서 스테이지 조작 의심
	currentStage := server.Users[UUID].CurrentStage
	if requestPacket.CurrentStage != currentStage {
		return util.ResponseBaseMessage(util.BadRequest, "Check Your CurrentStage")
	}

	var itemIds []int
	server.Users[UUID].Gold = requestPacket.Gold
	user := server.Users[UUID]

	//-- 무기 슬롯이 2개이상 착용중인지 비교
	//-- 착용중인 무기 슬롯이 2개 이상일경우 착용중인 무기 id만 담고 이하일 경우 모든 무기 id 담기
	if len(user.Slot) >= 2 {
		for _, item := range user.Slot {
			itemIds = append(itemIds, item.Id)
		}
	} else {
		for _, item := range table.ItemWeaponTable {
			itemIds = append(itemIds, item.Id)
		}
	}
	//-- 효과 아이템 id 담기
	for _, item := range table.ItemEffectTable {
		itemIds = append(itemIds, item.Id)
	}

	//-- 위에 담긴 아이템 Id 리스트에서 겹치지 않게 랜덤한 아이템 Id 4종 추출
	randomItemIds := RandomItemId(itemIds)
	ingameShopItem := response.IngameShopItem{}
	for _, itemId := range randomItemIds {
		ingameShopItem.Id = itemId
		ingameShopItem.Price = table.ShopIngameTable[itemId].Price
		responsePacket.Items = append(responsePacket.Items, ingameShopItem)
	}
	responsePacket.Code = util.Success
	responsePacket.Message = "Success"
	return util.ResponseMessage(responsePacket)
}

func BuyIngameItem(UUID string, payload string) string {
	requestPacket := request.BuyIngameItem{}
	responsePacket := response.BuyIngameItem{}
	err := json.Unmarshal([]byte(payload), &requestPacket)
	if err != nil {
		return util.ResponseErrorMessage(util.BadRequest, err.Error())
	}

	currentStage := server.Users[UUID].CurrentStage
	itemId := requestPacket.ItemId
	if requestPacket.CurrentStage != currentStage {
		return util.ResponseBaseMessage(util.BadRequest, "Check Your CurrentStage")
	}

	_, exists := server.Users[UUID]
	if !exists {
		userData := server.User{}
		server.Users[UUID] = &userData
	}
	/*
		if len(user.Slot) > 0 {
			server.Users[UUID].Slot = []server.Weapon{}
		}
	*/
	//-- 구매하려는 아이템의 가격보다 골드 보유량이 많은지 적은지 체크
	if server.Users[UUID].Gold < table.ShopIngameTable[itemId].Price {
		return util.ResponseBaseMessage(util.BadRequest, "Check Your Gold")
	}
	//-- 유저 골드 차감
	server.Users[UUID].Gold -= table.ShopIngameTable[itemId].Price
	//-- 구매하였으니 현재 스테이지 +1
	server.Users[UUID].CurrentStage += 1

	//-- 아이템 타입에 따라 처리
	itemType := table.ItemTable[itemId].Type
	switch itemType {
	case util.TYPE_WEAPON:
		//-- 이미 착용중인 슬롯에 있을 경우 인챈트 +1 아닐경우 새로운 슬롯에 추가
		isFound := false
		for idx, val := range server.Users[UUID].Slot {
			if itemId == val.Id {
				isFound = true
				server.Users[UUID].Slot[idx].Enchant++
				break
			}
		}
		if !isFound {
			weapon := server.Weapon{}
			weapon.Enchant = 0
			weapon.Id = itemId
			server.Users[UUID].Slot = append(server.Users[UUID].Slot, weapon)
		}
	case util.TYPE_EFFECT:
		//-- 이미 보유중인 효과 아이템일 경우 갯수 +1 아닐경우 리스트에 추가
		isFound := false
		for idx, val := range server.Users[UUID].Effect {
			if itemId == val.Id {
				isFound = true
				server.Users[UUID].Effect[idx].Count++
				break
			}
		}
		if !isFound {
			effect := server.Effect{}
			effect.Id = itemId
			effect.Count = 1
			server.Users[UUID].Effect = append(server.Users[UUID].Effect, effect)
		}
	}

	//-- Response값 등록
	responsePacket.Gold = server.Users[UUID].Gold
	responsePacket.CurrentStage = server.Users[UUID].CurrentStage
	slot := response.Weapon{}
	for _, item := range server.Users[UUID].Slot {
		slot.Id = item.Id
		slot.Enchant = item.Enchant
		responsePacket.Slot = append(responsePacket.Slot, slot)
	}
	effect := response.Effect{}
	for _, item := range server.Users[UUID].Effect {
		effect.Id = item.Id
		effect.Count = item.Count
		responsePacket.Effect = append(responsePacket.Effect, effect)
	}

	responsePacket.Code = util.Success
	responsePacket.Message = "Success"

	return util.ResponseMessage(responsePacket)
}

func ChangeUserName(UUID string, payload string) string {
	requestPacket := request.ChnageUserName{}
	responsePacket := packet.ResponsePacket{}
	err := json.Unmarshal([]byte(payload), &requestPacket)
	if err != nil {
		return util.ResponseErrorMessage(util.BadRequest, err.Error())
	}

	ctx := context.Background()
	db := server.DBManager.Login
	var count int
	query := `SELECT COUNT(*) FROM account WHERE user_name = ?`
	countErr := db.QueryRowContext(ctx, query, requestPacket.UserName).Scan(&count)
	if countErr != nil {
		return util.ResponseErrorMessage(util.ServerError, countErr.Error())
	}

	if count > 0 {
		responsePacket.Code = util.Fail
		responsePacket.Message = "Already UserName"
		return util.ResponseMessage(responsePacket)
	}

	query = `UPDATE account SET user_name = ? WHERE uid = ?`
	_, updateErr := db.ExecContext(ctx, query, requestPacket.UserName, UUID)
	if updateErr != nil {
		return util.ResponseErrorMessage(util.ServerError, updateErr.Error())
	}

	responsePacket.Code = util.Success
	responsePacket.Message = "Success"

	return util.ResponseMessage(responsePacket)
}

func LoadTimeAttackRankTable(UUID string, payload string) string {
	responsePacket := response.LoadTimeAttackRankTable{}

	ctx := context.Background()
	db := server.DBManager.Login
	query := `SELECT user_name, record_time FROM time_attack_rank ORDER BY record_time ASC LIMIT 10`
	result, err := db.QueryContext(ctx, query)
	defer result.Close()
	userData := response.TimeAttackUser{}

	userData.Rank = 1
	for result.Next() {
		if err != nil {
			if err == sql.ErrNoRows {
				responsePacket.Code = util.NotFound
				responsePacket.Message = "Empty Rank"
				return util.ResponseMessage(responsePacket)
			}
			return util.ResponseErrorMessage(util.ServerError, err.Error())
		}
		result.Scan(&userData.UserName, &userData.RecordTime)
		responsePacket.RankList = append(responsePacket.RankList, userData)
		userData.Rank++
	}
	responsePacket.Code = util.Success
	responsePacket.Message = "Success"

	return util.ResponseMessage(responsePacket)
}

/*
기존 타임어택 값과 새로운 타임어택 값 비교후
새로운 타임어택 값이 클경우 db갱신

신규일 경우 insert

랭킹 순위 유저중 해당될 될경우
재화 보상 2배

rank변수 0일경우 랭킹에 들지 못함
0이 아닐경우 랭킹에 해당
*/
//-- 테스트 위해서 체크 랭킹 범위 체크해야됨

func UpdateTimeAttackRank(UUID string, payload string) string {
	requestPacket := request.UpdateTimeAttackRank{}
	responsePacket := response.UpdateTimeAttackRank{}
	err := json.Unmarshal([]byte(payload), &requestPacket)
	if err != nil {
		return util.ResponseErrorMessage(util.BadRequest, err.Error())
	}
	userName := GetUserName(UUID)
	var recordTime float64
	var rank int
	rank = 1
	responsePacket.Rank = 0
	clearMoney := 200

	ctx := context.Background()
	db := server.DBManager.Login

	//-- 기존 타임어택 정보 가져오기
	query := `SELECT record_time FROM time_attack_rank WHERE uid = ?`
	rankErr := db.QueryRowContext(ctx, query, UUID).Scan(&recordTime)
	if rankErr != nil {
		//-- 랭킹 최초 등록시
		if rankErr == sql.ErrNoRows {
			query = `INSERT INTO time_attack_rank (uid, user_name, record_time) VALUES (?, ?, ?)`
			_, insertErr := db.ExecContext(ctx, query, UUID, userName, requestPacket.RecordTime)
			if insertErr != nil {
				return util.ResponseErrorMessage(util.ServerError, insertErr.Error())
			}

			query := `UPDATE account SET money = money + ? WHERE uid = ?`
			_, moneyErr := db.ExecContext(ctx, query, clearMoney, UUID)
			if moneyErr != nil {
				return util.ResponseErrorMessage(util.ServerError, moneyErr.Error())
			}

			var money int
			query = `SELECT money FROM account WHERE uid = ?`
			moneyErr = db.QueryRowContext(ctx, query, UUID).Scan(&money)
			if moneyErr != nil {
				return util.ResponseErrorMessage(util.ServerError, moneyErr.Error())
			}
			responsePacket.Money = money
			responsePacket.RewardMoney = clearMoney

			responsePacket.Code = util.Success
			responsePacket.Message = "Success"
			responsePacket.RecordTime = requestPacket.RecordTime
			return util.ResponseMessage(responsePacket)
		}
		return util.ResponseErrorMessage(util.ServerError, rankErr.Error())
	}

	responsePacket.RecordTime = recordTime

	//-- 기존 기록 갱신시 갱신된 값 저장
	if requestPacket.RecordTime < recordTime {
		query = `UPDATE time_attack_rank SET record_time = ? WHERE uid = ?`
		_, rankErr = db.ExecContext(ctx, query, requestPacket.RecordTime, UUID)
		if rankErr != nil {
			return util.ResponseErrorMessage(util.ServerError, rankErr.Error())
		}

		responsePacket.RecordTime = requestPacket.RecordTime
	}

	//-- 랭킹 순위에 들어가 있는지 (현재 10위 셋팅)
	query = `SELECT uid FROM time_attack_rank ORDER BY record_time ASC LIMIT 3`
	result, rankErr := db.QueryContext(ctx, query)
	var rankerId string

	if rankErr != nil {
		return util.ResponseErrorMessage(util.ServerError, rankErr.Error())
	}
	for result.Next() {
		result.Scan(&rankerId)
		if UUID == rankerId {
			responsePacket.Rank = rank
			clearMoney *= 2
			break
		}
		rank++
	}

	//-- 보상 재화 계정에 추가
	query = `UPDATE account SET money = money + ? WHERE uid = ?`
	_, moneyErr := db.ExecContext(ctx, query, clearMoney, UUID)
	if moneyErr != nil {
		return util.ResponseErrorMessage(util.ServerError, moneyErr.Error())
	}

	//-- 현재 재화 불러오기
	query = `SELECT money FROM account WHERE uid = ?`
	moneyErr = db.QueryRowContext(ctx, query, UUID).Scan(&responsePacket.Money)
	if moneyErr != nil {
		return util.ResponseErrorMessage(util.ServerError, moneyErr.Error())
	}
	responsePacket.Code = util.Success
	responsePacket.Message = "Success"
	responsePacket.RewardMoney = clearMoney

	return util.ResponseMessage(responsePacket)
}

func GameOver(UUID string, payload string) string {
	responsePacket := response.GameOver{}

	gameOverMoney := 100

	ctx := context.Background()
	db := server.DBManager.Login

	//-- 보상 재화 계정에 추가
	query := `UPDATE account SET money = money + ? WHERE uid = ?`
	_, moneyErr := db.ExecContext(ctx, query, gameOverMoney, UUID)
	if moneyErr != nil {
		return util.ResponseErrorMessage(util.ServerError, moneyErr.Error())
	}

	//-- 현재 재화 불러오기
	query = `SELECT money FROM account WHERE uid = ?`
	moneyErr = db.QueryRowContext(ctx, query, UUID).Scan(&responsePacket.Money)
	if moneyErr != nil {
		return util.ResponseErrorMessage(util.ServerError, moneyErr.Error())
	}
	responsePacket.Code = util.Success
	responsePacket.Message = "Success"
	responsePacket.RewardMoney = gameOverMoney
	return util.ResponseMessage(responsePacket)
}

/*
해당 유저가 가지고 있는 Money에서 구매하려는 아이템 가격이상 만큼 보유했는지 비교후 결과 Return
구매 가능시 True값과 구매 이후 잔액 Return
ㄴ구매 가격만큼 차감 Update
잔액 부족으로 인한 구매 실패시  False값과 0 Return
*/
func CheckMoney(UUID string, amount int) (bool, int) {
	db := server.DBManager.Login
	ctx := context.Background()
	var isCheckMoney bool
	var money int
	// 트랜잭션 시작
	tx, err := db.Begin()
	if err != nil {
		println("ChemMoney Transaction Begin Error!!")
		println(err.Error())
		return false, 0
	}

	// UPDATE 쿼리 작성 및 실행
	query := `UPDATE account SET money = money - ? WHERE uid = ? AND money >= ?`
	result, err := tx.ExecContext(ctx, query, amount, UUID, amount)
	if err != nil {
		println("CheckMoney Transaction Update Error!!")
		println(err.Error())
		tx.Rollback()
		return false, 0
	}

	// 영향 받은 행의 수 확인
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		println("CheckMoney Transaction RowsAffected Error!!")
		println(err.Error())
		tx.Rollback()
		return false, 0
	}

	query = `SELECT money FROM account WHERE uid =?`
	moneyErr := tx.QueryRowContext(ctx, query, UUID).Scan(&money)
	if moneyErr != nil {
		println("CheckMoney Transaction Select Error!!")
		println(err.Error())
		tx.Rollback()
		return false, 0
	}

	// 트랜잭션 커밋
	err = tx.Commit()
	if err != nil {
		println("CheckMoney Transaction Commit Error!!")
		println(err.Error())
		tx.Rollback()
		return false, 0
	}
	isCheckMoney = rowsAffected > 0

	// 영향 받은 행이 1개 이상이면 성공, 그렇지 않으면 실패
	return isCheckMoney, money
}

func MoneyRollback(UUID string, amount int) {
	db := server.DBManager.Login
	ctx := context.Background()
	query := `UPDATE money = money+? FROM account WHERE uid = ?`
	_, err := db.ExecContext(ctx, query, amount, UUID)
	if err != nil {
		println("MoneyRollback Exec Error!!")
		println("Error User UUID:", UUID, "/Amount:", amount)
		println(err.Error())
		return
	}

	println("MoneyRollBack Success!! UUID:", UUID, "/Amount:", amount)
}

func GetWeaponEnchant(UUID string, itemId int) (int, error) {
	db := server.DBManager.Game
	ctx := context.Background()
	var enchant int
	query := `SELECT enchant_level FROM inventory WHERE uid = ? AND item_id = ?`
	err := db.QueryRowContext(ctx, query, UUID, itemId).Scan(&enchant)

	if err != nil {
		if err != sql.ErrNoRows {
			println("GetWeaponEnchant QueryRow Error!!")
		}
		println(err.Error())
		return 0, err
	}

	return enchant, nil
}

func RandomItemId(itemIds []int) []int {
	//-- 매개변수 itemIds의 사이즈가 4보다 작거나 같을경우 종료
	if len(itemIds) <= 4 {
		return itemIds
	}

	// Seed 설정 (시드값 설정)
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	r := rand.New(source)

	//-- 매개변수 itemIds 사이즈만큼의 랜덤 범위중 무작위 값 추출
	random := r.Intn(len(itemIds))

	//-- 무작위 추출된 값을 인덱스로 대입하여 해당 인덱스 값 슬라이스에서 제거
	itemIds = append(itemIds[:random], itemIds[random+1:]...)

	//-- 매개변수 itemIds값이 4이하가 될때까지 재귀
	return RandomItemId(itemIds)
}

func GetUserName(UUID string) string {
	var userName string
	ctx := context.Background()
	db := server.DBManager.Login
	query := `SELECT user_name FROM account WHERE uid = ?`
	err := db.QueryRowContext(ctx, query, UUID).Scan(&userName)
	if err != nil {
		println("GetUserName Error!!::", err.Error())
		return ""
	}

	return userName
}
