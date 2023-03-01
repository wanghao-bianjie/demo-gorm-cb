package callback

import (
	"demo-gorm-cb/aes"
	"fmt"
	"gorm.io/gorm"
	"reflect"
)

//AfterQuery 只支持查询返回的结果是单表model的结构体，map、自定义的其他结构体（联表、部分字段）暂不支持
func AfterQuery(db *gorm.DB) {
	defer func() {
		if err := recover(); err != nil {
			db.AddError(fmt.Errorf("recover panic:%v", err))
		}
	}()
	if db.Error == nil && db.Statement.Schema != nil && db.RowsAffected > 0 && !db.Statement.SkipHooks {
		/*if db.Statement.Dest != nil {
			valueOfDest := reflect.ValueOf(db.Statement.Dest)
			typeOfDest := reflect.TypeOf(db.Statement.Dest)
			println(typeOfDest.Kind().String())
			switch typeOfDest.Kind() {
			case reflect.Ptr:
				indirect := reflect.Indirect(valueOfDest)
				println(indirect.Kind().String())
				switch indirect.Kind() {
				case reflect.Slice, reflect.Array:
					//for i := 0; i < valueOfDest.Len(); i++ {
					//	typeOfDest := reflect.TypeOf(db.Statement.Dest)
					//	for i := 0; i < indirect.NumField(); i++ {
					//		field := typeOfDest.Field(i)
					//		println(field.Name)
					//	}
					//}

				case reflect.Struct:
					//typeOfDest := reflect.TypeOf(db.Statement.Dest)

					fn := func(v reflect.Value) {
						println(v.NumField())
						for i := 0; i < indirect.NumField(); i++ {
							field := indirect.Field(i)
							switch field.Kind() {
							case reflect.Struct:
								//fn()

							case reflect.String:
								field.SetString(AesDecryptMock(field.String()))
							}
							println(field.Kind().String())
						}
					}
					fn(indirect)

				}
			default:
				return
			}
		}*/
		destReflectValue := getReflectValueElem(db.Statement.Dest)
		for fieldIndex, field := range db.Statement.Schema.Fields {
			switch destReflectValue.Kind() {
			case reflect.Slice, reflect.Array: //[]struct
				for i := 0; i < destReflectValue.Len(); i++ {
					index := destReflectValue.Index(i)
					if index.Kind() != reflect.Struct {
						continue
					}
					if index.NumField() != len(db.Statement.Schema.Fields) {
						return
					}
					if fieldValue, isZero := field.ValueOf(index); !isZero { // 从字段中获取数值
						if index.Type().Field(fieldIndex).Name != field.Name || index.Type().Field(fieldIndex).Type.Kind() != field.FieldType.Kind() {
							return
						}
						if dbAesTableColumnMap[[2]string{db.Statement.Schema.Table, field.DBName}] {
							if s, ok := fieldValue.(string); ok {
								_ = db.AddError(field.Set(destReflectValue.Index(i), aes.DecryptMock(s)))
							}
						}
					}
				}
			case reflect.Struct: //struct
				if destReflectValue.NumField() != len(db.Statement.Schema.Fields) {
					return
				}
				if fieldValue, isZero := field.ValueOf(destReflectValue); !isZero { // 从字段中获取数值
					if destReflectValue.Type().Field(fieldIndex).Name != field.Name {
						return
					}
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
