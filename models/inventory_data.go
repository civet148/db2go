package models

import github_com_civet148_db2go_types "github.com/civet148/db2go/types"
import "github.com/civet148/sqlca/v3"

const TableNameInventoryData = "inventory_data" //产品库存数据表

const (
	INVENTORY_DATA_COLUMN_ID            = "id"
	INVENTORY_DATA_COLUMN_CREATE_ID     = "create_id"
	INVENTORY_DATA_COLUMN_CREATE_NAME   = "create_name"
	INVENTORY_DATA_COLUMN_CREATE_TIME   = "create_time"
	INVENTORY_DATA_COLUMN_UPDATE_ID     = "update_id"
	INVENTORY_DATA_COLUMN_UPDATE_NAME   = "update_name"
	INVENTORY_DATA_COLUMN_UPDATE_TIME   = "update_time"
	INVENTORY_DATA_COLUMN_IS_FROZEN     = "is_frozen"
	INVENTORY_DATA_COLUMN_NAME          = "name"
	INVENTORY_DATA_COLUMN_SERIAL_NO     = "serial_no"
	INVENTORY_DATA_COLUMN_QUANTITY      = "quantity"
	INVENTORY_DATA_COLUMN_PRICE         = "price"
	INVENTORY_DATA_COLUMN_PRODUCT_EXTRA = "product_extra"
	INVENTORY_DATA_COLUMN_LOCATION      = "location"
)

type InventoryData struct {
	github_com_civet148_db2go_types.BaseModel
	Id           uint64        `json:"id,omitempty" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`                                  //产品ID
	CreateId     uint64        `json:"create_id,omitempty" db:"create_id" gorm:"column:create_id;type:bigint unsigned;default:0;"`       //创建人ID
	CreateName   string        `json:"create_name,omitempty" db:"create_name" gorm:"column:create_name;type:varchar(64);"`               //创建人姓名
	UpdateId     uint64        `json:"update_id,omitempty" db:"update_id" gorm:"column:update_id;type:bigint unsigned;default:0;"`       //更新人ID
	UpdateName   string        `json:"update_name,omitempty" db:"update_name" gorm:"column:update_name;type:varchar(64);"`               //更新人姓名
	IsFrozen     int8          `json:"is_frozen,omitempty" db:"is_frozen" gorm:"column:is_frozen;type:tinyint(1);default:0;"`            //冻结状态(0: 未冻结 1: 已冻结)
	Name         string        `json:"name,omitempty" db:"name" gorm:"column:name;type:varchar(255);"`                                   //产品名称
	SerialNo     string        `json:"serial_no,omitempty" db:"serial_no" gorm:"column:serial_no;type:varchar(64);"`                     //产品编号
	Quantity     sqlca.Decimal `json:"quantity,omitempty" db:"quantity" gorm:"column:quantity;type:decimal(16,3);default:0.000;"`        //产品库存
	Price        sqlca.Decimal `json:"price,omitempty" db:"price" gorm:"column:price;type:decimal(16,2);default:0.00;"`                  //产品均价
	ProductExtra string        `json:"product_extra,omitempty" db:"product_extra" gorm:"column:product_extra;type:text;" sqlca:"isnull"` //产品附带数据(JSON文本)
	Location     sqlca.Point   `json:"location,omitempty" db:"location" gorm:"column:location;type:point;" sqlca:"isnull"`               //地理位置
}

func (do InventoryData) TableName() string            { return "inventory_data" }
func (do InventoryData) GetId() uint64                { return do.Id }
func (do *InventoryData) SetId(v uint64)              { do.Id = v }
func (do InventoryData) GetCreateId() uint64          { return do.CreateId }
func (do *InventoryData) SetCreateId(v uint64)        { do.CreateId = v }
func (do InventoryData) GetCreateName() string        { return do.CreateName }
func (do *InventoryData) SetCreateName(v string)      { do.CreateName = v }
func (do InventoryData) GetCreateTime() string        { return do.CreateTime }
func (do *InventoryData) SetCreateTime(v string)      { do.CreateTime = v }
func (do InventoryData) GetUpdateId() uint64          { return do.UpdateId }
func (do *InventoryData) SetUpdateId(v uint64)        { do.UpdateId = v }
func (do InventoryData) GetUpdateName() string        { return do.UpdateName }
func (do *InventoryData) SetUpdateName(v string)      { do.UpdateName = v }
func (do InventoryData) GetUpdateTime() string        { return do.UpdateTime }
func (do *InventoryData) SetUpdateTime(v string)      { do.UpdateTime = v }
func (do InventoryData) GetIsFrozen() int8            { return do.IsFrozen }
func (do *InventoryData) SetIsFrozen(v int8)          { do.IsFrozen = v }
func (do InventoryData) GetName() string              { return do.Name }
func (do *InventoryData) SetName(v string)            { do.Name = v }
func (do InventoryData) GetSerialNo() string          { return do.SerialNo }
func (do *InventoryData) SetSerialNo(v string)        { do.SerialNo = v }
func (do InventoryData) GetQuantity() sqlca.Decimal   { return do.Quantity }
func (do *InventoryData) SetQuantity(v sqlca.Decimal) { do.Quantity = v }
func (do InventoryData) GetPrice() sqlca.Decimal      { return do.Price }
func (do *InventoryData) SetPrice(v sqlca.Decimal)    { do.Price = v }
func (do InventoryData) GetProductExtra() string      { return do.ProductExtra }
func (do *InventoryData) SetProductExtra(v string)    { do.ProductExtra = v }
func (do InventoryData) GetLocation() sqlca.Point     { return do.Location }
func (do *InventoryData) SetLocation(v sqlca.Point)   { do.Location = v }

////////////////////// ----- 自定义代码请写在下面 ----- //////////////////////
