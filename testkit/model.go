package testkit

type TestModel struct {
	ID        int64  `gorm:"column:id;primary_key"`
	Value1    int64  `gorm:"column:value1"`
	Value2    int64  `gorm:"column:value2"`
	Value3    int64  `gorm:"column:value3"`
	Value4    int64  `gorm:"column:value4"`
	Value5    int64  `gorm:"column:value5"`
	Value6    int64  `gorm:"column:value6"`
	Value7    int64  `gorm:"column:value7"`
	Value8    int64  `gorm:"column:value8"`
	PtrValue1 *int64 `gorm:"column:ptr_value1"`
}

const (
	TestModelTableName = "gorm_cache_model"
)

func (m *TestModel) TableName() string {
	return TestModelTableName
}
