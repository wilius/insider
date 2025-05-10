package message

import (
	"gorm.io/gorm"
	"insider/constants"
	"time"
)

type entity struct {
	ID                int64                   `gorm:"primaryKey;column:id"`
	PhoneNumber       string                  `gorm:"column:phone_number;not null"`
	Message           string                  `gorm:"column:message;not null"`
	Status            constants.MessageStatus `gorm:"column:status;not null"`
	ProviderMessageID *string                 `gorm:"column:provider_message_id"`
	Provider          *string                 `gorm:"column:provider"`
	CreateDate        time.Time               `gorm:"column:create_date;not null"`
	UpdateDate        *time.Time              `gorm:"column:update_date;not null"`
}

func (*entity) TableName() string {
	return "notifications.message"
}

func (c *entity) BeforeCreate(tx *gorm.DB) (err error) {
	var nextID int64
	err = tx.Raw("SELECT nextval('notifications.seq__message_id')").Scan(&nextID).Error
	if err != nil {
		return err
	}
	c.ID = nextID
	now := time.Now().UTC()
	c.CreateDate = now
	return nil
}
