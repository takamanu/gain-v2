package database

import (
	"fmt"

	// DataFashion "gain-v2/features/fashions/data"

	DataUser "gain-v2/features/users/data"

	// DataVoucher "gain-v2/features/vouchers/data"

	"gorm.io/gorm"
)

func MigrateWithDrop(db *gorm.DB) {

	db.Exec("DROP SCHEMA public CASCADE;")
	db.Exec("CREATE SCHEMA public;")

	db.Exec("CREATE TYPE roles AS ENUM ('customer', 'super_admin', 'staff', 'agent');")
	db.Exec("CREATE TYPE status AS ENUM ('active', 'not_active');")
	fmt.Println("[MIGRATION] Success creating enum types for roles and status")

	// USER DATA MANAGEMENT \\
	db.AutoMigrate(DataUser.User{})
	fmt.Println("[MIGRATION] Success creating user")

}
