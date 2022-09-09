package testkit

import "gorm.io/gorm"

func PrepareTableAndData(db *gorm.DB) error {
	err := db.Migrator().CreateTable(&TestModel{})
	if err != nil {
		return err
	}

	models := make([]TestModel, 0)
	for i := 1; i <= testSize; i++ {
		_pValue := int64(i)
		model := TestModel{
			ID:        int64(i),
			Value1:    int64(i),
			Value2:    int64(i),
			Value3:    int64(i),
			Value4:    int64(i),
			Value5:    int64(i),
			Value6:    int64(i),
			Value7:    int64(i),
			Value8:    int64(i),
			PtrValue1: &_pValue,
		}
		models = append(models, model)
	}

	db.CreateInBatches(models, 2000)
	return db.Error
}

func CleanTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&TestModel{})
}
