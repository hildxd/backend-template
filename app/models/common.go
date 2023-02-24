package models

// 自增ID主键
type ID struct {
	ID uint `json:"id" gorm:"primary_key"`
}

// 创建/更新时间
type Timestamps struct {
	CreatedAt int64 `json:"created_at" gorm:"comment:'创建时间'"`
	UpdatedAt int64 `json:"updated_at" gorm:"comment:'更新时间'"`
}

// 软删除
type SoftDelete struct {
	DeletedAt int64 `json:"deleted_at" gorm:"comment:'删除时间'"`
}
