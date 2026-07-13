package models

import "time"
import "github.com/civet148/sqlca/v3"

const TableNameInventoryIn = "inventory_in" //

const (
	InventoryInColumn_Id         = "id"
	InventoryInColumn_CreatedAt  = "created_at"
	InventoryInColumn_UpdatedAt  = "updated_at"
	InventoryInColumn_IsDeleted  = "is_deleted"
	InventoryInColumn_DeleteTime = "delete_time"
	InventoryInColumn_ProductId  = "product_id"
	InventoryInColumn_OrderNo    = "order_no"
	InventoryInColumn_UserId     = "user_id"
	InventoryInColumn_UserName   = "user_name"
	InventoryInColumn_Quantity   = "quantity"
	InventoryInColumn_Weight     = "weight"
	InventoryInColumn_Remark     = "remark"
	InventoryInColumn_CreateId   = "create_id"
	InventoryInColumn_CreateName = "create_name"
	InventoryInColumn_UpdateId   = "update_id"
	InventoryInColumn_UpdateName = "update_name"
)

type InventoryIn struct {
	Id         uint64        `json:"id" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`                                                                          //
	IsDeleted  int8          `json:"is_deleted" db:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);default:0;" sqlca:"isnull"`                                  //
	DeleteTime *time.Time    `json:"delete_time" db:"delete_time" gorm:"column:delete_time;type:datetime;default:null;" sqlca:"isnull"`                              //
	ProductId  uint64        `json:"product_id" db:"product_id" gorm:"column:product_id;type:bigint unsigned;index:idx_prod_create_id;default:null;" sqlca:"isnull"` //
	OrderNo    string        `json:"order_no" db:"order_no" gorm:"column:order_no;type:varchar(64);uniqueIndex:UNIQ_ORDER_NO;default:null;" sqlca:"isnull"`          //
	UserId     uint64        `json:"user_id" db:"user_id" gorm:"column:user_id;type:bigint unsigned;default:0;" sqlca:"isnull"`                                      //
	UserName   string        `json:"user_name" db:"user_name" gorm:"column:user_name;type:varchar(64);default:null;" sqlca:"isnull"`                                 //
	Quantity   sqlca.Decimal `json:"quantity" db:"quantity" gorm:"column:quantity;type:decimal(16,6);default:0.000000;" sqlca:"isnull"`                              //
	Weight     sqlca.Decimal `json:"weight" db:"weight" gorm:"column:weight;type:decimal(16,6);default:0.000000;" sqlca:"isnull"`                                    //
	Remark     string        `json:"remark" db:"remark" gorm:"column:remark;type:varchar(512);default:null;" sqlca:"isnull"`                                         //
	CreateId   uint64        `json:"create_id" db:"create_id" gorm:"column:create_id;type:bigint unsigned;index:idx_prod_create_id;default:0;" sqlca:"isnull"`       //
	CreateName string        `json:"create_name" db:"create_name" gorm:"column:create_name;type:varchar(64);default:null;" sqlca:"isnull"`                           //
	UpdateId   uint64        `json:"update_id" db:"update_id" gorm:"column:update_id;type:bigint unsigned;default:0;" sqlca:"isnull"`                                //
	UpdateName string        `json:"update_name" db:"update_name" gorm:"column:update_name;type:varchar(64);default:null;" sqlca:"isnull"`                           //
}

func (do InventoryIn) TableName() string { return "inventory_in" }

func (do InventoryIn) GetId() uint64 { return do.Id }

func (do InventoryIn) GetIsDeleted() int8 { return do.IsDeleted }

func (do InventoryIn) GetDeleteTime() *time.Time { return do.DeleteTime }

func (do InventoryIn) GetProductId() uint64 { return do.ProductId }

func (do InventoryIn) GetOrderNo() string { return do.OrderNo }

func (do InventoryIn) GetUserId() uint64 { return do.UserId }

func (do InventoryIn) GetUserName() string { return do.UserName }

func (do InventoryIn) GetQuantity() sqlca.Decimal { return do.Quantity }

func (do InventoryIn) GetWeight() sqlca.Decimal { return do.Weight }

func (do InventoryIn) GetRemark() string { return do.Remark }

func (do InventoryIn) GetCreateId() uint64 { return do.CreateId }

func (do InventoryIn) GetCreateName() string { return do.CreateName }

func (do InventoryIn) GetUpdateId() uint64 { return do.UpdateId }

func (do InventoryIn) GetUpdateName() string { return do.UpdateName }

func (do *InventoryIn) SetId(v uint64) { do.Id = v }

func (do *InventoryIn) SetIsDeleted(v int8) { do.IsDeleted = v }

func (do *InventoryIn) SetDeleteTime(v *time.Time) { do.DeleteTime = v }

func (do *InventoryIn) SetProductId(v uint64) { do.ProductId = v }

func (do *InventoryIn) SetOrderNo(v string) { do.OrderNo = v }

func (do *InventoryIn) SetUserId(v uint64) { do.UserId = v }

func (do *InventoryIn) SetUserName(v string) { do.UserName = v }

func (do *InventoryIn) SetQuantity(v sqlca.Decimal) { do.Quantity = v }

func (do *InventoryIn) SetWeight(v sqlca.Decimal) { do.Weight = v }

func (do *InventoryIn) SetRemark(v string) { do.Remark = v }

func (do *InventoryIn) SetCreateId(v uint64) { do.CreateId = v }

func (do *InventoryIn) SetCreateName(v string) { do.CreateName = v }

func (do *InventoryIn) SetUpdateId(v uint64) { do.UpdateId = v }

func (do *InventoryIn) SetUpdateName(v string) { do.UpdateName = v }
