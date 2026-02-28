package models

import "time"
import "github.com/civet148/sqlca/v3"

const TableNameInventoryIn = "inventory_in" //

const (
	INVENTORY_IN_COLUMN_ID          = "id"
	INVENTORY_IN_COLUMN_CREATED_AT  = "created_at"
	INVENTORY_IN_COLUMN_UPDATED_AT  = "updated_at"
	INVENTORY_IN_COLUMN_IS_DELETED  = "is_deleted"
	INVENTORY_IN_COLUMN_DELETE_TIME = "delete_time"
	INVENTORY_IN_COLUMN_PRODUCT_ID  = "product_id"
	INVENTORY_IN_COLUMN_ORDER_NO    = "order_no"
	INVENTORY_IN_COLUMN_USER_ID     = "user_id"
	INVENTORY_IN_COLUMN_USER_NAME   = "user_name"
	INVENTORY_IN_COLUMN_QUANTITY    = "quantity"
	INVENTORY_IN_COLUMN_WEIGHT      = "weight"
	INVENTORY_IN_COLUMN_REMARK      = "remark"
	INVENTORY_IN_COLUMN_CREATE_ID   = "create_id"
	INVENTORY_IN_COLUMN_CREATE_NAME = "create_name"
	INVENTORY_IN_COLUMN_UPDATE_ID   = "update_id"
	INVENTORY_IN_COLUMN_UPDATE_NAME = "update_name"
)

type InventoryIn struct {
	BaseModel
	Id         uint64        `json:"id,omitempty" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`                                                    //
	IsDeleted  int8          `json:"is_deleted,omitempty" db:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);default:0;" sqlca:"isnull"`            //
	DeleteTime *time.Time    `json:"delete_time,omitempty" db:"delete_time" gorm:"column:delete_time;type:datetime;" sqlca:"isnull"`                     //
	ProductId  uint64        `json:"product_id,omitempty" db:"product_id" gorm:"column:product_id;type:bigint unsigned;" sqlca:"isnull"`                 //
	OrderNo    string        `json:"order_no,omitempty" db:"order_no" gorm:"column:order_no;type:varchar(64);uniqueIndex:UNIQ_ORDER_NO;" sqlca:"isnull"` //
	UserId     uint64        `json:"user_id,omitempty" db:"user_id" gorm:"column:user_id;type:bigint unsigned;default:0;" sqlca:"isnull"`                //
	UserName   string        `json:"user_name,omitempty" db:"user_name" gorm:"column:user_name;type:varchar(64);" sqlca:"isnull"`                        //
	Quantity   sqlca.Decimal `json:"quantity,omitempty" db:"quantity" gorm:"column:quantity;type:decimal(16,6);default:0.000000;" sqlca:"isnull"`        //
	Weight     sqlca.Decimal `json:"weight,omitempty" db:"weight" gorm:"column:weight;type:decimal(16,6);default:0.000000;" sqlca:"isnull"`              //
	Remark     string        `json:"remark,omitempty" db:"remark" gorm:"column:remark;type:varchar(512);" sqlca:"isnull"`                                //
	CreateId   uint64        `json:"create_id,omitempty" db:"create_id" gorm:"column:create_id;type:bigint unsigned;default:0;" sqlca:"isnull"`          //
	CreateName string        `json:"create_name,omitempty" db:"create_name" gorm:"column:create_name;type:varchar(64);" sqlca:"isnull"`                  //
	UpdateId   uint64        `json:"update_id,omitempty" db:"update_id" gorm:"column:update_id;type:bigint unsigned;default:0;" sqlca:"isnull"`          //
	UpdateName string        `json:"update_name,omitempty" db:"update_name" gorm:"column:update_name;type:varchar(64);" sqlca:"isnull"`                  //
}

func (do InventoryIn) TableName() string { return "inventory_in" }

func (do InventoryIn) GetId() uint64              { return do.Id }
func (do InventoryIn) GetCreatedAt() time.Time    { return do.CreatedAt }
func (do InventoryIn) GetUpdatedAt() time.Time    { return do.UpdatedAt }
func (do InventoryIn) GetIsDeleted() int8         { return do.IsDeleted }
func (do InventoryIn) GetDeleteTime() *time.Time  { return do.DeleteTime }
func (do InventoryIn) GetProductId() uint64       { return do.ProductId }
func (do InventoryIn) GetOrderNo() string         { return do.OrderNo }
func (do InventoryIn) GetUserId() uint64          { return do.UserId }
func (do InventoryIn) GetUserName() string        { return do.UserName }
func (do InventoryIn) GetQuantity() sqlca.Decimal { return do.Quantity }
func (do InventoryIn) GetWeight() sqlca.Decimal   { return do.Weight }
func (do InventoryIn) GetRemark() string          { return do.Remark }
func (do InventoryIn) GetCreateId() uint64        { return do.CreateId }
func (do InventoryIn) GetCreateName() string      { return do.CreateName }
func (do InventoryIn) GetUpdateId() uint64        { return do.UpdateId }
func (do InventoryIn) GetUpdateName() string      { return do.UpdateName }

func (do *InventoryIn) SetId(v uint64)              { do.Id = v }
func (do *InventoryIn) SetCreatedAt(v time.Time)    { do.CreatedAt = v }
func (do *InventoryIn) SetUpdatedAt(v time.Time)    { do.UpdatedAt = v }
func (do *InventoryIn) SetIsDeleted(v int8)         { do.IsDeleted = v }
func (do *InventoryIn) SetDeleteTime(v *time.Time)  { do.DeleteTime = v }
func (do *InventoryIn) SetProductId(v uint64)       { do.ProductId = v }
func (do *InventoryIn) SetOrderNo(v string)         { do.OrderNo = v }
func (do *InventoryIn) SetUserId(v uint64)          { do.UserId = v }
func (do *InventoryIn) SetUserName(v string)        { do.UserName = v }
func (do *InventoryIn) SetQuantity(v sqlca.Decimal) { do.Quantity = v }
func (do *InventoryIn) SetWeight(v sqlca.Decimal)   { do.Weight = v }
func (do *InventoryIn) SetRemark(v string)          { do.Remark = v }
func (do *InventoryIn) SetCreateId(v uint64)        { do.CreateId = v }
func (do *InventoryIn) SetCreateName(v string)      { do.CreateName = v }
func (do *InventoryIn) SetUpdateId(v uint64)        { do.UpdateId = v }
func (do *InventoryIn) SetUpdateName(v string)      { do.UpdateName = v }

/*
CREATE TABLE `inventory_in` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `is_deleted` tinyint(1) DEFAULT '0',
  `delete_time` datetime DEFAULT NULL,
  `product_id` bigint unsigned DEFAULT NULL,
  `order_no` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `user_id` bigint unsigned DEFAULT '0',
  `user_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `quantity` decimal(16,6) DEFAULT '0.000000',
  `weight` decimal(16,6) DEFAULT '0.000000',
  `remark` varchar(512) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `create_id` bigint unsigned DEFAULT '0',
  `create_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `update_id` bigint unsigned DEFAULT '0',
  `update_name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_ORDER_NO` (`order_no`),
  KEY `idx_inventory_in_created_at` (`created_at`),
  KEY `idx_inventory_in_updated_at` (`updated_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2027587914958901249 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
*/

////////////////////// ----- 自定义代码请写在下面 ----- //////////////////////
