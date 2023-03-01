package callback

import (
	"demo-gorm-cb/aes"
	"demo-gorm-cb/model"
	"reflect"
)

var user model.User

var dbAesTableColumnMap = map[[2]string]bool{
	[2]string{user.TableName(), model.GetUserColumn().Name}:        false,
	[2]string{user.TableName(), model.GetUserColumn().PhoneNumber}: true,
	[2]string{user.TableName(), model.GetUserColumn().Address}:     true,
	[2]string{user.TableName(), model.GetUserColumn().IdNo}:        true,
}

func getReflectValueElem(i interface{}) reflect.Value {
	value := reflect.ValueOf(i)
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return value
}

func encryptDestMap(table string, dest map[string]interface{}) map[string]interface{} {
	var newDest = make(map[string]interface{}, len(dest))
	for field, value := range dest {
		if dbAesTableColumnMap[[2]string{table, field}] {
			if s, ok := value.(string); ok {
				newDest[field] = aes.EncryptMock(s)
			} else {
				newDest[field] = value
			}
		} else {
			newDest[field] = value
		}
	}
	return newDest
}

func encryptDestMapSlice(table string, dest []map[string]interface{}) []map[string]interface{} {
	var newDest = make([]map[string]interface{}, 0, len(dest))
	for _, m := range dest {
		newDest = append(newDest, encryptDestMap(table, m))
	}
	return newDest
}
