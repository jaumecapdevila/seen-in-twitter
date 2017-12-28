package persistence

import (
	"github.com/jinzhu/gorm"
)

// Tweet model
type Tweet struct {
	gorm.Model
	Text string
}
