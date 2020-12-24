package gorm_test

import (
	"log"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shopspring/decimal"
	"strings"
)

var (
	db *gorm.DB
)

func init() {
	var err error

	db, err = gorm.Open("mysql", "root:@(0.0.0.0:3306)/peatio_dev?charset=utf8&parseTime=True&loc=Asia%2FShanghai")
	if err != nil {
		log.Fatal(err.Error())
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(10)
	db.DB().SetConnMaxLifetime(14400 * time.Second)
	db.Set("gorm:table_options", "ENGINE=InnoDB")

	db.LogMode(true)
}

// testing association
type Market struct {
	ID              string
	AskUnit         string
	BidUnit         string
	AskFee          decimal.Decimal
	BidFee          decimal.Decimal
	AskPrecision    int8
	BidPrecision    int8
	AskDecimalLimit int
	BidMin          decimal.Decimal
	PriceUpLimit    int
	PriceDownLimit  int
	Position        int
	Enabled         int8
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (m *Market) TableName() string {
	return "markets"
}

func TestTableName(t *testing.T) {
	t.Log(db.HasTable(&Market{}))
	t.Log(db.HasTable("markets"))
}

type IDentity struct {
	ID             int
	Email          string
	PhoneNumber    string
	PasswordDigest string
	PasswordHash   string
	IsActive       bool
	RetryCount     int
	IsLocked       bool
	LockedAt       time.Time
	LastVerifyAt   time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Member struct {
	ID                int
	Level             int8
	Sn                string
	Email             string
	PhoneNumber       string
	JwtToken          string
	ExpireIn          int
	ClientPublicKey   string
	PinCode           string
	PinSalt           string
	IDentityID        int
	Disabled          int8
	APIDisabled       bool
	AppDeviceType     string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	MaskPhoneNumber   string
	HashedPhoneNumber string
	RegistrationToken string
	AppDeviceLang     string

	// used for preload
	Orders []Order
}

type IDDocument struct {
	ID                int
	FirstName         string
	LastName          string
	MemberID          int
	BirthDate         time.Time
	Country           string
	City              string
	Address           string
	Zipcode           string
	Last4ssn          string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	EncryptedUsername string
	HashedUsername    string
	State             int8
	Region            string
	Gdpr              int
	Occupation        string
	IDmState          int
}

func TestBelongsTo(t *testing.T) {
	var member = Member{ID: 1}
	var document IDDocument

	db.Model(&member).Related(&document)
	t.Log(document.ID)
}

func TestBelongsToReverse(t *testing.T) {
	var document = IDDocument{MemberID: 1}
	var member Member

	db.Model(&document).Related(&member)
	t.Log(member.ID)
}

func TestHasOne(t *testing.T) {
	var identity = IDentity{ID: 860}
	var member Member

	db.Model(&identity).Related(&member)
	t.Log(member.ID)
}

func TestHasOneReverse(t *testing.T) {
	var member = Member{IDentityID: 860}
	var identity IDentity

	db.Model(&member).Related(&identity)
	t.Log(identity.ID)
}

type Order struct {
	ID            int
	Bid           string
	Ask           string
	MarketID      string
	Price         decimal.Decimal
	Volume        decimal.Decimal
	OriginVolume  decimal.Decimal
	Fee           decimal.Decimal
	State         int
	Type          string
	MemberID      int
	OrdType       string
	Locked        decimal.Decimal
	OriginLocked  decimal.Decimal
	FundsReceived decimal.Decimal
	TradesCount   int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Trade struct {
	ID          int
	Price       decimal.Decimal
	Volume      decimal.Decimal
	AskID       int
	BidID       int
	Trend       int
	MarketID    string
	AskMemberID int
	BidMemberID int
	Funds       decimal.Decimal
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PriceUsd    decimal.Decimal

	Bills []Bill `gorm:"ForeignKey:BillableID;polymorphic:Billable;polymorphic_value:Trade"`
}

func TestHasMany(t *testing.T) {
	var member = Member{ID: 1}
	var orders []Order

	db.Model(&member).Related(&orders)
	t.Log(len(orders))
}

func TestHasManyReverse(t *testing.T) {
	var order = Order{MemberID: 1}
	var member Member

	db.Model(&order).Related(&member)
	t.Log(member.ID)
}

type AdminRole struct {
	ID         int
	Name       string
	AdminItems []AdminItem `gorm:"many2many:admin_items_admin_roless;association_jointable_foreignkey:admin_items_id;jointable_foreignkey:admin_roles_id;"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type AdminItemsAdminRoles struct {
	ID           int
	AdminItemsID int
	AdminRolesID int
}

type AdminItem struct {
	ID              int
	Desc            string
	Key             string
	Route           string
	AdminCategoryID int
	AdminRoles      []AdminRole `gorm:"many2many:admin_items_admin_roless;association_jointable_foreignkey:admin_roles_id;jointable_foreignkey:admin_items_id;"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func TestManyToMany(t *testing.T) {
	var role AdminRole
	var items []AdminItem

	db.First(&role)
	t.Log(role.ID)

	db.Model(&role).Related(&items, "AdminItems")
	t.Log(len(items))
}

func TestManyToManyReverse(t *testing.T) {
	var item AdminItem
	var roles []AdminRole

	db.First(&item)
	t.Log(item.ID)

	db.Model(&item).Related(&roles, "AdminRoles")
	t.Log(len(roles))
}

func TestPreload(t *testing.T) {
	var members []Member
	db.Where([]int{1}).Preload("Orders").Find(&members)
}

func TestHasManyOrder2Trade(t *testing.T) {
	var order = Order{ID: 1486}
	var trades []Trade

	db.First(&order)
	db.Model(&order).Related(&trades, strings.ToLower(order.Type[5:])+"_id")
	t.Log(len(trades))
}

type Bill struct {
	ID               int
	MemberID         int
	CurrencyID       string
	Price            decimal.Decimal
	Volume           decimal.Decimal
	PriceUsd         decimal.Decimal
	BillType         string
	Description      string
	BillableType     string
	BillableID       int
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Currencies       string
	Fee              decimal.Decimal
	FeeUnit          string
	OriginFee        decimal.Decimal
	OriginFeeUnit    string
	DiscountPriceUsd decimal.Decimal
}

func TestHasManyPolymorphism(t *testing.T) {
	var trade = Trade{ID: 1490826}

	var bills []Bill
	db.Model(&trade).Related(&bills, "Bills")
	t.Log(len(bills))
}

func TestQueryPreloadBillsFromTrade(t *testing.T) {
	var trades []Trade
	db.Where([]int{1490826}).Preload("Bills").Find(&trades)
	t.Log(len(trades[0].Bills))
}

func TestQueryPreloadBillsFromOrder(t *testing.T) {
	var order = Order{ID: 1486}
	var trades []Trade

	db.First(&order)
	db = db.Model(&order).Related(&trades, strings.ToLower(order.Type[5:])+"_id")

	var tradeIds []int
	for _, trade := range trades {
		tradeIds = append(tradeIds, trade.ID)
	}

	db.Where(tradeIds).Preload("Bills").Find(&trades)
	t.Log(len(trades))

	for _, trade := range trades {
		t.Logf("trade: %d, bills: %d", trade.ID, len(trade.Bills))
	}
}

// testing CRUD
func TestCreate(t *testing.T) {
	// db.Create
	var member = Member{ID: 111111}
	db.Create(&member)
}

func TestQueryFirst(t *testing.T) {
	var market Market

	// first - order by primary key
	db.First(&market)
	t.Log(market.ID, market.AskFee.String(), market.BidFee.String())
}

func TestQueryFirstByPrimaryKey(t *testing.T) {
	var member Member

	// first - order by primary key
	db.First(&member, 1)
	t.Log(member)
}

func TestQueryTake(t *testing.T) {
	var market Market

	// take - no specified order
	db.Take(&market)
	t.Log(market.ID, market.AskFee.String(), market.BidFee.String())
}

func TestQueryLast(t *testing.T) {
	var market Market

	// last - order by primary key
	db.Last(&market)
	t.Log(market.ID, market.AskFee.String(), market.BidFee.String())
}

func TestQueryFind(t *testing.T) {
	var orders []Order

	// find - get all records
	db.Find(&orders)
	t.Log(len(orders))
}

func TestQueryWherePlainSQL(t *testing.T) {
	var order Order
	db.Where("market_id = ?", "wonbtc").First(&order)
	t.Log(order.ID)

	var orders []Order
	db.Where("market_id = ?", "wonbtc").Find(&orders)
	t.Log(len(orders))

	db.Where("market_id <> ?", "wonbtc").Find(&orders)
	t.Log(len(orders))

	db.Where("market_id IN (?)", []string{"topwon", "wonbtc"}).Find(&orders)
	t.Log(len(orders))

	db.Where("market_id LIKE ?", "%won%").Find(&orders)
	t.Log(len(orders))

	db.Where("market_id = ? AND member_id = ?", "topwon", 1).Find(&orders)
	t.Log(len(orders))

	db.Where("created_at > ?", "2019-01-01").Find(&orders)
	t.Log(len(orders))

	db.Where("created_at BETWEEN ? AND ?", "2019-01-01", "2019-03-01").Find(&orders)
	t.Log(len(orders))
}

func TestQueryWhereStructMapSlice(t *testing.T) {
	// struct

	// map

	// slice
}

func TestQueryNotOr(t *testing.T) {
	// Not

	// Or
}

func TestQueryForUpdate(t *testing.T) {
	// db.Set("gorm:query_option", "FOR UPDATE").First(&user, 10)
}

func TestQueryFirstOrInit(t *testing.T) {
	var member Member
	db.FirstOrInit(&member, Member{ID: 1})
	t.Log(member.ID)

	var market Market
	db.FirstOrInit(&market, map[string]interface{}{"ask_unit": "won"})
	t.Log(market.ID)

	// Attrs
	// Initialize struct with argument if record not found

	// Assign
	// Assign argument to struct regardless it is found or not
}

func TestQueryFirstOrCreate(t *testing.T) {
	// Get first matched record, or create a new one with given conditions (only works with struct, map conditions)

	// Attrs
	// Initialize struct with argument if record not found

	// Assign
	// Assign argument to struct regardless it is found or not
}

func TestQueryAdvanced(t *testing.T) {
	// SubQuery

	// Select

	// Order

	// Limit

	// Offset

	// Count

	// Group & Having

	// Joins

	// Pluck
	var ids []string
	db.Model(&Market{}).Pluck("Id", &ids)
	t.Log(len(ids))

	// Scan
}

func TestUpdate(t *testing.T) {
	// Save
	var member Member
	db.First(&member)
	member.Sn = "WONBQ1BRQMOAPI"
	db.Save(&member)

	// Update & Updates
	db.Model(&Member{ID: 1}).UpdateColumn("sn", "WONAQ1BRQMOAPI")

	// UpdateColumn & UpdateColumns

	// gorm.Expr
}

func TestDelete(t *testing.T) {
	// Delete Record

	// Batch Delete

	// Soft Delete

	// Delete record permanently
}

func TestHooks(t *testing.T) {
	// Create
	// BeforeCreate() (err error)
	// BeforeSave() (err error)
	// AfterCreate(tx *gorm.DB)
	// AfterSave() (err error)

	// Update
	// BeforeUpdate() (err error)
	// BeforeSave() (err error)
	// AfterUpdate()
	// AfterSave() (err error)

	// Query
	// AfterFind()

	// Delete
	// BeforeDelete() (err error)
	// AfterDelete() (err error)
}

// deadlock
func TestTransaction(t *testing.T) {
    tx1 := db.Begin()

    tx1.Model(&Member{ID: 1}).UpdateColumn("sn", "WONAQ1BRQMOAPI")

    tx2 := db.Begin()
    tx2.Model(&Member{ID: 1}).UpdateColumn("sn", "WON4926QRQPAPI")
    tx2.Commit()

    tx3 := db.Begin()
    tx3.Model(&Member{ID: 1}).UpdateColumn("sn", "WONO5Y5ZVO9API")
    tx3.Commit()

    tx1.Commit()
}

func TestRaw(t *testing.T) {
    var count int
    db.Raw("select count(*) from members").Count(&count)
    t.Log(count)
}
