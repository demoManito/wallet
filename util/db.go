package util

import "gorm.io/gorm"

func ForUpdateLock(db *gorm.DB, record interface{}, id int64, fc func(tx *gorm.DB)) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(record, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 业务逻辑
	fc(tx)
	// ======
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
