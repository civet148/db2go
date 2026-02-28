package models

import "time"
import "github.com/civet148/sqlca/v3"

const TableNameInventoryData = "inventory_data" //

const (
	INVENTORY_DATA_COLUMN_ID            = "id"
	INVENTORY_DATA_COLUMN_CREATED_AT    = "created_at"
	INVENTORY_DATA_COLUMN_UPDATED_AT    = "updated_at"
	INVENTORY_DATA_COLUMN_IS_FROZEN     = "is_frozen"
	INVENTORY_DATA_COLUMN_NAME          = "name"
	INVENTORY_DATA_COLUMN_SERIAL_NO     = "serial_no"
	INVENTORY_DATA_COLUMN_QUANTITY      = "quantity"
	INVENTORY_DATA_COLUMN_PRICE         = "price"
	INVENTORY_DATA_COLUMN_LOCATION      = "location"
	INVENTORY_DATA_COLUMN_PRODUCT_EXTRA = "product_extra"
	INVENTORY_DATA_COLUMN_CREATE_ID     = "create_id"
	INVENTORY_DATA_COLUMN_CREATE_NAME   = "create_name"
	INVENTORY_DATA_COLUMN_UPDATE_ID     = "update_id"
	INVENTORY_DATA_COLUMN_UPDATE_NAME   = "update_name"
)

type InventoryData struct {
	BaseModel
	Id           uint64        `json:"id,omitempty" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`                                                             //
	IsFrozen     int8          `json:"is_frozen,omitempty" db:"is_frozen" gorm:"column:is_frozen;type:tinyint(1);default:0;" sqlca:"isnull"`                        //
	Name         string        `json:"name,omitempty" db:"name" gorm:"column:name;type:varchar(255);comment:产品：名称；不能为空;" sqlca:"isnull"`                            //产品：名称；不能为空
	SerialNo     string        `json:"serial_no,omitempty" db:"serial_no" gorm:"column:serial_no;type:varchar(64);index:i_serial_no;comment:产品序列号;" sqlca:"isnull"` //产品序列号
	Quantity     sqlca.Decimal `json:"quantity,omitempty" db:"quantity" gorm:"column:quantity;type:decimal(16,3);default:0.000;" sqlca:"isnull"`                    //
	Price        sqlca.Decimal `json:"price,omitempty" db:"price" gorm:"column:price;type:decimal(16,2);default:0.00;" sqlca:"isnull"`                              //
	Location     sqlca.Point   `json:"location,omitempty" db:"location" gorm:"column:location;type:point;" sqlca:"isnull"`                                          //
	ProductExtra struct{}      `json:"product_extra,omitempty" db:"product_extra" gorm:"column:product_extra;type:json;" sqlca:"isnull"`                            //
	CreateId     uint64        `json:"create_id,omitempty" db:"create_id" gorm:"column:create_id;type:bigint unsigned;default:0;" sqlca:"isnull"`                   //
	CreateName   string        `json:"create_name,omitempty" db:"create_name" gorm:"column:create_name;type:varchar(64);" sqlca:"isnull"`                           //
	UpdateId     uint64        `json:"update_id,omitempty" db:"update_id" gorm:"column:update_id;type:bigint unsigned;default:0;" sqlca:"isnull"`                   //
	UpdateName   string        `json:"update_name,omitempty" db:"update_name" gorm:"column:update_name;type:varchar(64);" sqlca:"isnull"`                           //
}

func (do InventoryData) TableName() string { return "inventory_data" }

func (do InventoryData) GetId() uint64              { return do.Id }
func (do InventoryData) GetCreatedAt() time.Time    { return do.CreatedAt }
func (do InventoryData) GetUpdatedAt() time.Time    { return do.UpdatedAt }
func (do InventoryData) GetIsFrozen() int8          { return do.IsFrozen }
func (do InventoryData) GetName() string            { return do.Name }
func (do InventoryData) GetSerialNo() string        { return do.SerialNo }
func (do InventoryData) GetQuantity() sqlca.Decimal { return do.Quantity }
func (do InventoryData) GetPrice() sqlca.Decimal    { return do.Price }
func (do InventoryData) GetLocation() sqlca.Point   { return do.Location }
func (do InventoryData) GetProductExtra() struct{}  { return do.ProductExtra }
func (do InventoryData) GetCreateId() uint64        { return do.CreateId }
func (do InventoryData) GetCreateName() string      { return do.CreateName }
func (do InventoryData) GetUpdateId() uint64        { return do.UpdateId }
func (do InventoryData) GetUpdateName() string      { return do.UpdateName }

func (do *InventoryData) SetId(v uint64)              { do.Id = v }
func (do *InventoryData) SetCreatedAt(v time.Time)    { do.CreatedAt = v }
func (do *InventoryData) SetUpdatedAt(v time.Time)    { do.UpdatedAt = v }
func (do *InventoryData) SetIsFrozen(v int8)          { do.IsFrozen = v }
func (do *InventoryData) SetName(v string)            { do.Name = v }
func (do *InventoryData) SetSerialNo(v string)        { do.SerialNo = v }
func (do *InventoryData) SetQuantity(v sqlca.Decimal) { do.Quantity = v }
func (do *InventoryData) SetPrice(v sqlca.Decimal)    { do.Price = v }
func (do *InventoryData) SetLocation(v sqlca.Point)   { do.Location = v }
func (do *InventoryData) SetProductExtra(v struct{})  { do.ProductExtra = v }
func (do *InventoryData) SetCreateId(v uint64)        { do.CreateId = v }
func (do *InventoryData) SetCreateName(v string)      { do.CreateName = v }
func (do *InventoryData) SetUpdateId(v uint64)        { do.UpdateId = v }
func (do *InventoryData) SetUpdateName(v string)      { do.UpdateName = v }

/*
CREATE TABLE `inventory_data` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `is_frozen` tinyint(1) DEFAULT '0',
  `name` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '产品：名称；不能为空',
  `serial_no` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '产品序列号',
  `quantity` decimal(16,3) DEFAULT '0.000',
  `price` decimal(16,2) DEFAULT '0.00',
  `location` point DEFAULT NULL,
  `product_extra` json DEFAULT NULL,
  `create_id` bigint unsigned DEFAULT '0',
  `create_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `update_id` bigint unsigned DEFAULT '0',
  `update_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_inventory_data_created_at` (`created_at`),
  KEY `idx_inventory_data_updated_at` (`updated_at`),
  KEY `i_serial_no` (`serial_no`)
) ENGINE=InnoDB AUTO_INCREMENT=2027587914967289858 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
*/

////////////////////// ----- 自定义代码请写在下面 ----- //////////////////////
