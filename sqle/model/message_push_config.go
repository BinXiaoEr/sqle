package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Strings []string

func (t *Strings) Scan(value interface{}) error {
	bytesValue, _ := value.([]byte)
	return json.Unmarshal(bytesValue, t)
}

func (t Strings) Value() (driver.Value, error) {
	return json.Marshal(t)
}

type MessagePushConfig struct {
	Model
	ProjectId         string    `json:"project_id" gorm:"type:varchar(255)"`
	Type              string    `json:"type" gorm:"type:varchar(255)"`
	PushFrequencyType string    `json:"push_frequency_type"  gorm:"type:varchar(255)"`
	PushFrequencyCron string    `json:"cron" gorm:"type:varchar(255)"`
	PushUserType      string    `json:"push_user_Type" gorm:"type:varchar(255)"`
	PushUserList      Strings   `json:"push_user_list"`
	LastPushTime      time.Time `json:"last_push_time" gorm:"type:datetime(3)"`
	Status            string    `json:"status" gorm:"type:varchar(255)"`
}
