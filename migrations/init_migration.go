package migration

import (
	"friend_connection_rest_api/services/friendship"
	"friend_connection_rest_api/services/user"

	"gorm.io/gorm"
)

func InitMigration(dbconn *gorm.DB) {

	if oke := dbconn.Migrator().HasTable(&user.Users{}); !oke {
		dbconn.AutoMigrate(&user.Users{})
	}

	if oke := dbconn.Migrator().HasTable(&friendship.Friendship{}); !oke {
		dbconn.AutoMigrate(&friendship.Friendship{})
	}
}
