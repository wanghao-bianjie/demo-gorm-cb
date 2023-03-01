package callback

import (
	"demo-gorm-cb/aes"
	"fmt"
	"gorm.io/gorm"
	"reflect"
)

func BeforeCreate(db *gorm.DB) {
	defer func() {
		if err := recover(); err != nil {
			db.AddError(fmt.Errorf("recover panic:%v", err))
		}
	}()
	if db.Error == nil && db.Statement.Schema != nil && !db.Statement.SkipHooks {
		if db.Statement.Dest != nil {
			switch db.Statement.Dest.(type) {
			case []map[string]interface{}: //[]map
				dest := db.Statement.Dest.([]map[string]interface{})
				db.Statement.Dest = encryptDestMapSlice(db.Statement.Schema.Table, dest)
				return
			case *[]map[string]interface{}: //*[]map
				dest := db.Statement.Dest.(*[]map[string]interface{})
				db.Statement.Dest = encryptDestMapSlice(db.Statement.Schema.Table, *dest)
				return
			case map[string]interface{}: //map
				dest := db.Statement.Dest.(map[string]interface{})
				db.Statement.Dest = encryptDestMap(db.Statement.Schema.Table, dest)
				return
			case *map[string]interface{}: //*map
				dest := db.Statement.Dest.(*map[string]interface{})
				db.Statement.Dest = encryptDestMap(db.Statement.Schema.Table, *dest)
				return
			}
		}

		destReflectValue := getReflectValueElem(db.Statement.Dest)
		for _, field := range db.Statement.Schema.Fields {
			switch destReflectValue.Kind() {
			case reflect.Slice, reflect.Array: //[]struct
				for i := 0; i < destReflectValue.Len(); i++ {
					index := destReflectValue.Index(i)
					if destReflectValue.Index(i).Kind() != reflect.Struct {
						continue
					}
					if fieldValue, isZero := field.ValueOf(index); !isZero { // 从字段中获取数值
						if dbAesTableColumnMap[[2]string{db.Statement.Schema.Table, field.DBName}] {
							_ = db.AddError(field.Set(index, aes.EncryptMock(fieldValue.(string))))
						}
					}
				}
			case reflect.Struct: //struct
				// 从字段中获取数值
				if fieldValue, isZero := field.ValueOf(destReflectValue); !isZero {
					if dbAesTableColumnMap[[2]string{db.Statement.Schema.Table, field.DBName}] {
						_ = db.AddError(field.Set(destReflectValue, aes.EncryptMock(fieldValue.(string))))
					}
				}
			}
		}
	}
}

func AfterCreate(db *gorm.DB) {
	defer func() {
		if err := recover(); err != nil {
			db.AddError(fmt.Errorf("recover panic:%v", err))
		}
	}()
	if db.Error == nil && db.Statement.Schema != nil && !db.Statement.SkipHooks {
		if db.Statement.Dest != nil {
			switch db.Statement.Dest.(type) {
			case []map[string]interface{}, map[string]interface{}, *[]map[string]interface{}, *map[string]interface{}:
				return
			}
		}

		destReflectValue := getReflectValueElem(db.Statement.Dest)
		for _, field := range db.Statement.Schema.Fields {
			switch destReflectValue.Kind() {
			case reflect.Slice, reflect.Array: //[]struct
				for i := 0; i < destReflectValue.Len(); i++ {
					// 从字段中获取数值
					if destReflectValue.Index(i).Kind() != reflect.Struct {
						continue
					}
					if fieldValue, isZero := field.ValueOf(destReflectValue.Index(i)); !isZero {
						if dbAesTableColumnMap[[2]string{db.Statement.Schema.Table, field.DBName}] {
							if s, ok := fieldValue.(string); ok {
								_ = db.AddError(field.Set(destReflectValue.Index(i), aes.DecryptMock(s)))
							}
						}
					}
				}
			case reflect.Struct: //struct
				// 从字段中获取数值
				if fieldValue, isZero := field.ValueOf(destReflectValue); !isZero {
					if dbAesTableColumnMap[[2]string{db.Statement.Schema.Table, field.DBName}] {
						if s, ok := fieldValue.(string); ok {
							_ = db.AddError(field.Set(destReflectValue, aes.DecryptMock(s)))
						}
					}
				}
			}
		}
	}
}
