package callback

import (
	"demo-gorm-cb/aes"
	"fmt"
	"gorm.io/gorm"
	"reflect"
)

func BeforeUpdate(db *gorm.DB) {
	defer func() {
		if err := recover(); err != nil {
			db.AddError(fmt.Errorf("recover panic:%v", err))
		}
	}()
	if db.Error == nil && !db.Statement.SkipHooks {
		if db.Statement.Dest != nil {
			switch db.Statement.Dest.(type) {
			case map[string]interface{}: //map
				dest := db.Statement.Dest.(map[string]interface{})
				db.Statement.Dest = encryptDestMap(db.Statement.Table, dest)
				return
			case *map[string]interface{}: //*map
				dest := db.Statement.Dest.(*map[string]interface{})
				db.Statement.Dest = encryptDestMap(db.Statement.Table, *dest)
				return
			}
		}

		if db.Statement.Schema == nil {
			return
		}
		destReflectValue := getReflectValueElem(db.Statement.Dest)
		for _, field := range db.Statement.Schema.Fields {
			switch destReflectValue.Kind() {
			case reflect.Struct: //struct
				if fieldValue, isZero := field.ValueOf(destReflectValue); !isZero { // 从字段中获取数值
					if dbAesTableColumnMap[[2]string{db.Statement.Schema.Table, field.DBName}] {
						if s, ok := fieldValue.(string); ok {
							//todo
							if destReflectValue.CanSet() {
								_ = db.AddError(field.Set(destReflectValue, aes.EncryptMock(s)))
							} else {
								newStruct := reflect.New(destReflectValue.Type())
								for i := 0; i < destReflectValue.Type().NumField(); i++ {
									if destReflectValue.Type().Field(i).Name == field.Name {
										newStruct.Elem().Field(i).SetString(aes.EncryptMock(s))
									} else {
										newStruct.Elem().Field(i).Set(destReflectValue.Field(i))
									}
								}
								db.Statement.Dest = newStruct.Elem().Interface()
							}
						}
					}
				}
			}
		}
	}
}

func AfterUpdate(db *gorm.DB) {
	defer func() {
		if err := recover(); err != nil {
			db.AddError(fmt.Errorf("recover panic:%v", err))
		}
	}()
	if db.Error == nil && db.Statement.Schema != nil {
		destIsMap := false
		if db.Statement.Dest != nil {
			switch db.Statement.Dest.(type) {
			case []map[string]interface{}, map[string]interface{}, *[]map[string]interface{}, *map[string]interface{}:
				destIsMap = true
			}
		}

		fn := func(reflectValue reflect.Value) {
			for _, field := range db.Statement.Schema.Fields {
				switch reflectValue.Kind() {
				case reflect.Struct: //struct
					if fieldValue, isZero := field.ValueOf(reflectValue); !isZero { // 从字段中获取数值
						if dbAesTableColumnMap[[2]string{db.Statement.Schema.Table, field.DBName}] {
							if s, ok := fieldValue.(string); ok {
								if reflectValue.CanSet() {
									_ = db.AddError(field.Set(reflectValue, aes.DecryptMock(s)))
								}
							}
						}
					}
				}
			}
		}
		destReflectValue := getReflectValueElem(db.Statement.Dest)
		if !destIsMap {
			fn(destReflectValue)
		}
		if destReflectValue != db.Statement.ReflectValue {
			fn(db.Statement.ReflectValue)
		}
	}
}
