package table

import (
	"context"
	"project_yd/server"
)

type item struct {
	Id        int
	Name      string
	Type      int
	Category  int
	ImageName string
	IsStack   bool
}
type shop struct {
	Id        int
	MoneyType int
	Price     int
}
type weaponEnchantProbability struct {
	Enchant     int
	Probability int
	Price       int
}
type itemWeapon struct {
	Id     int
	Name   string
	Damage int
	Speed  float64
	Range  int
}
type itemEffect struct {
	Id          int
	Name        string
	MaxHp       float64
	RegenHp     float64
	Speed       float64
	Damage      float64
	AttackSpeed float64
}
type shopIngame struct {
	Id    int
	Price int
}

var ItemTable map[int]item
var ShopTable map[int]shop
var WeaponEnchantProbabilityTable map[int]weaponEnchantProbability
var ItemWeaponTable map[int]itemWeapon
var ItemEffectTable map[int]itemEffect
var ShopIngameTable map[int]shopIngame

func LoadLoginDatabaseTables() {
	println("LoadLoginDatabaseTables -- Start")
	ItemTable = make(map[int]item)
	ItemWeaponTable = make(map[int]itemWeapon)
	ItemEffectTable = make(map[int]itemEffect)
	ShopTable = make(map[int]shop)
	ShopIngameTable = make(map[int]shopIngame)

	WeaponEnchantProbabilityTable = make(map[int]weaponEnchantProbability)

	LoadItemTable()
	LoadItemWeaponTable()
	LoadItemEffectTable()
	LoadShopTable()
	LoadShopIngameTable()
	LoadWeaponEnchantProbabilityTable()
	println("LoadLoginDatabaseTables -- End")
}

func LoadItemTable() {
	println("LoadItemTable -- Start")
	query := `SELECT * FROM items`
	db := server.DBManager.Login
	ctx := context.Background()
	result, err := db.QueryContext(ctx, query)

	item := item{}
	for result.Next() {
		if err != nil {
			println("LoadItemTable Error!! QueryContext Sql::", query)
			println(err.Error())
		}
		result.Scan(&item.Id, &item.Name, &item.Type, &item.Category, &item.ImageName, &item.IsStack)
		ItemTable[item.Id] = item
	}
	println("LoadItemTable Success!! Size::", len(ItemTable))
}

func LoadShopTable() {
	println("LoadShopTable -- Start")

	query := `SELECT * FROM shop`
	db := server.DBManager.Login
	ctx := context.Background()
	result, err := db.QueryContext(ctx, query)

	item := shop{}
	for result.Next() {
		if err != nil {
			println("LoadShopTable Error!! QueryContext Sql::", query)
			println(err.Error())
		}
		result.Scan(&item.Id, &item.MoneyType, &item.Price)
		ShopTable[item.Id] = item
	}

	println("LoadShopTable Success!! Size::", len(ShopTable))
}

func LoadWeaponEnchantProbabilityTable() {
	println("LoadWeaponEnchantProbabilityTable -- Start")

	query := `SELECT * FROM weapon_enchant_probability`
	db := server.DBManager.Login
	ctx := context.Background()
	result, err := db.QueryContext(ctx, query)

	weapon := weaponEnchantProbability{}
	for result.Next() {
		if err != nil {
			println("LoadWeaponEnchantProbabilityTable Error!! QueryContext Sql::", query)
			println(err.Error())
		}
		result.Scan(&weapon.Enchant, &weapon.Probability, &weapon.Price)
		WeaponEnchantProbabilityTable[weapon.Enchant] = weapon
	}

	println("LoadWeaponEnchantProbabilityTable Success!! Size::", len(WeaponEnchantProbabilityTable))
}

func LoadItemWeaponTable() {
	println("LoadItemWeaponTable -- Start")

	query := `SELECT * FROM item_weapon`
	db := server.DBManager.Login
	ctx := context.Background()
	result, err := db.QueryContext(ctx, query)

	item := itemWeapon{}
	for result.Next() {
		if err != nil {
			println("LoadItemWeaponTable Error!! QueryContext Sql::", query)
			println(err.Error())
		}
		result.Scan(&item.Id, &item.Name, &item.Damage, &item.Speed, &item.Range)
		ItemWeaponTable[item.Id] = item
	}
	println("LoadItemWeaponTable Success!! Size::", len(ItemWeaponTable))
}

func LoadItemEffectTable() {
	println("LoadItemEffectTable -- Start")

	query := `SELECT * FROM item_effect`
	db := server.DBManager.Login
	ctx := context.Background()
	result, err := db.QueryContext(ctx, query)

	item := itemEffect{}
	for result.Next() {
		if err != nil {
			println("LoadItemEffectTable Error!! QueryContext Sql::", query)
			println(err.Error())
		}
		result.Scan(&item.Id, &item.Name, &item.MaxHp, &item.RegenHp, &item.Speed, &item.Damage, &item.AttackSpeed)
		ItemEffectTable[item.Id] = item
	}
	println("LoadItemEffectTable Success!! Size::", len(ItemEffectTable))
}

func LoadShopIngameTable() {
	println("LoadShopIngameTable -- Start")

	query := `SELECT * FROM shop_ingame`
	db := server.DBManager.Login
	ctx := context.Background()
	result, err := db.QueryContext(ctx, query)

	item := shopIngame{}
	for result.Next() {
		if err != nil {
			println("LoadShopIngameTable Error!! QueryContext Sql::", query)
			println(err.Error())
		}
		result.Scan(&item.Id, &item.Price)
		ShopIngameTable[item.Id] = item
	}
	println("LoadShopIngameTable Success!! Size::", len(ShopIngameTable))
}
