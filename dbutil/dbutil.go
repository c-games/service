package dbutil

import (
	"database/sql"
	"fmt"
	"github.com/c-games/common/coll"
	"github.com/c-games/common/str"
	"reflect"
)


func GenDropTable(s interface{}) string {
	rfs := reflect.TypeOf(s)

	tableName := str.Pascal2Snake(rfs.Name()) // use struct name as default name
	fields := ""
	var pk []string
	for idx := 0 ; idx < rfs.NumField() ; idx++ {
		f := rfs.Field(idx)

		if f.Name == "ITable" {
			tableName = f.Tag.Get("name")
			continue
		}

		name := "`" + str.Pascal2Snake(f.Name) + "`"

		fields = fields + name + " " + f.Tag.Get("sql") + ",\n"

		_, ok := f.Tag.Lookup("pk")
		if ok {
			pk = append(pk, name)
		}

	}

	if tableName == "" {
		panic("you need a table name")
	} else {
		return "DROP TABLE `" + tableName + "`;"
	}
}

func GenCreateTable(s interface{}) string {
	rfs := reflect.TypeOf(s)

	tableName := str.Pascal2Snake(rfs.Name()) // use struct name as default name
	fields := ""
	var pk []string
	for idx := 0 ; idx < rfs.NumField() ; idx++ {
		f := rfs.Field(idx)

		if f.Name == "ITable" {
			tableName = f.Tag.Get("name")
			continue
		}

		name := "`" + str.Pascal2Snake(f.Name) + "`"

		fields = fields + name + " " + f.Tag.Get("sql") + ",\n"

		_, ok := f.Tag.Lookup("pk")
		if ok {
			pk = append(pk, name)
		}

	}

	var sqlString string
	if tableName ==  "" {
		panic("you need to set a table name")
	} else {
		sqlString = "CREATE TABLE `" + tableName + "` "
	}


	if len(fields) > 0 && len(pk) == 0{
		fields = fields[:len(fields) - 2]
	}

	pkStr := ""
	if len(pk) != 0 {
		pkStr = "PRIMARY KEY " + "(" + coll.JoinString(pk, ",") + ")"
	}


	sqlString = sqlString + "(\n" + fields + pkStr + ") ENGINE=INNODB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

	//fmt.Println(sqlString)
	return sqlString
}

func GenCreateIndex(s interface{}) string {

	rfs := reflect.TypeOf(s)
	sqlString := ""
	if rfs.Name() == "" {
		return ""
	}

	tableName := str.Pascal2Snake(rfs.Name())

	index := make(map[string][]string)
	var keys []string
	//var idx map[string][]string
	for idx := 0 ; idx < rfs.NumField() ; idx++ {
		f := rfs.Field(idx)
		indexName, ok := f.Tag.Lookup("index")
		if ok {
			_, ok := index[indexName]
			fieldName := "`" + str.Pascal2Snake(f.Name) + "`"
			if !ok {
				keys = append(keys, indexName)
				index[indexName] = []string{fieldName}
			} else {
				index[indexName] = append(index[indexName], fieldName)
			}

		}

	}

	// NOTE golang 的 map 不保證順序，所以要自己處理 keys
	for _, indexName := range keys {
		indexSet := index[indexName]
		idxs := coll.JoinString(indexSet, ",")
		sqlString = sqlString + "CREATE INDEX " + indexName + " ON `" + tableName + "` (" + idxs + ");"
	}

	return sqlString
}


func CompareParams(params []interface{}, expectParamCount int, expectParamTypes []reflect.Kind) error {
	totalParams := len(params)
	if totalParams != expectParamCount {
		panic("stmt 所需參數不同")
	}
	// check type
	for idx, param := range params {
		expectType := expectParamTypes[idx]
		paramType := reflect.TypeOf(param)

		if expectType != paramType.Kind() {
			panic(fmt.Sprintf("query stmt 資料型態有錯, expectType = %s, paramType = %s", expectType, paramType.Kind()))
		}
	}
	return nil
}

// NOTE 為了要共用 sql.rows 和 sql.row 的 Scan
// ---------------------------------------------------

type SqlRowLike interface {
	Scan(...interface{}) error
}

type Scannable interface {
	Scan(SqlRowLike) error
}

func QueryCondition(queryResult Scannable, rowLike SqlRowLike) error {
	err := queryResult.Scan(rowLike)
	if err == nil {
		return nil
	} else if err == sql.ErrNoRows  {
		return nil
	} else {
		// rows 的 Scan 有可能會因為沒有先 call Next() 而出錯，那需要自己處理，所以直接 return err
		return err
	}
}

func QueryCondition2(queryResultType reflect.Type, rowLike SqlRowLike) (interface{}, error) {
	newInstent := reflect.New(queryResultType)
	elementInstent := newInstent.Elem().Interface()
	// TODO 要 check elementInstent 有 Scannable
	err := QueryCondition(elementInstent.(Scannable), rowLike)
	return elementInstent, err
}


// key value map 轉成 sql update 的 string
func GenUpdateKeysAndValueList(updateData map[string]interface{}) (string, []interface{}) {
	kString := ""
	var vList []interface{}
	for k, v := range updateData {
		kString = kString + "," + k + "=?"
		vList = append(vList, v)
	}
	kString = kString[1:]
	return kString, vList
}

func GenInsertKeysAndValueList(updateData map[string]interface{}) (string, []interface{}) {
	kString := ""
	var vList []interface{}
	for k, v := range updateData {
		kString = kString + "," + k
		vList = append(vList, v)
	}
	kString = kString[1:]
	return kString, vList
}
