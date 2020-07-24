package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/c-games/common/coll"
	"github.com/c-games/common/fail"
	strutil "github.com/c-games/common/str"
	goMysql "github.com/go-sql-driver/mysql"
	"reflect"
	"strings"
	"time"
)

type sqlQueryExec func(rows *sql.Rows) (json.RawMessage, error)
type sqlQueryExec_retArrInterface func(rows *sql.Rows) ([]interface{}, error)

//mysql is a operator for mysql.
type MySQL struct {
	mainDB             *sql.DB
	readDB             *sql.DB
	mainDBConfig       *MysqlParameter
	readDBConfig       *MysqlParameter
	readWriteSplitting bool
}

//請看 mysql mysqlParameter
type MysqlParameter struct {
	DriverName   string
	User         string        // Username
	Password     string        // Password (requires User)
	Net          string        // Network type
	Address      string        // Network address (requires Net)
	DBName       string        // Database name
	Timeout      time.Duration // Dial timeout
	ReadTimeout  time.Duration // I/O read timeout
	WriteTimeout time.Duration // I/O write timeout
	ParseTime    bool          // parse db time
}

/*
type DBStats struct {
	MaxOpenConnections int // Maximum number of open connections to the database; added in Go 1.11

	// Pool Status
	OpenConnections int // The number of established connections both in use and idle.
	InUse           int // The number of connections currently in use; added in Go 1.11
	Idle            int // The number of idle connections; added in Go 1.11

	// Counters
	WaitCount         int64         // The total number of connections waited for; added in Go 1.11
	WaitDuration      time.Duration // The total time blocked waiting for a new connection; added in Go 1.11
	MaxIdleClosed     int64         // The total number of connections closed due to SetMaxIdleConns; added in Go 1.11
	MaxLifetimeClosed int64
}
*/

//newMysql returns mysql.
//readTimeout I/O read timeout 30s, 0.5m, 1m30s
//writeTimeout I/O write timeout
//timeout timeout for establishing connections
func NewMySQL(mainDBConfig *MysqlParameter, readWriteSplitting bool, readOnlyDBConfig *MysqlParameter) (*MySQL, error) {

	mainConf, err := getMySQLConfig(mainDBConfig)
	if err != nil {
		return nil, err
	}
	formatStr := mainConf.FormatDSN()

	mainDB, err := sql.Open(mainDBConfig.DriverName, formatStr)
	if err != nil {
		return nil, err
	}

	err = mainDB.Ping()
	if err != nil {
		return nil, err
	}

	var readOnlyConf *goMysql.Config
	formatStr = ""
	var readDB *sql.DB

	//有讀寫分離
	if readWriteSplitting {

		readOnlyConf, err = getMySQLConfig(readOnlyDBConfig)
		if err != nil {
			return nil, err
		}

		formatStr = readOnlyConf.FormatDSN()
		readDB, err := sql.Open(readOnlyDBConfig.DriverName, formatStr)
		if err != nil {
			return nil, err
		}
		err = readDB.Ping()
		if err != nil {
			return nil, err
		}
	}

	return &MySQL{
		mainDB:             mainDB,
		readDB:             readDB,
		mainDBConfig:       mainDBConfig,
		readDBConfig:       readOnlyDBConfig,
		readWriteSplitting: readWriteSplitting,
	}, nil
}

func getMySQLConfig(dbConfig *MysqlParameter) (*goMysql.Config, error) {

	if len(dbConfig.User) < 1 {
		return nil, fmt.Errorf("mysqlParameter length of %s =%d", "user", len(dbConfig.User))
	}
	if len(dbConfig.Password) < 1 {
		return nil, fmt.Errorf("mysqlParameter length of %s =%d", "password", len(dbConfig.Password))
	}
	if len(dbConfig.Net) < 1 {
		return nil, fmt.Errorf("mysqlParameter length of %s =%d", "Net", len(dbConfig.Net))
	}
	if len(dbConfig.Address) < 1 {
		return nil, fmt.Errorf("mysqlParameter length of %s =%d", "Addr", len(dbConfig.Address))
	}
	if len(dbConfig.DBName) < 1 {
		return nil, fmt.Errorf("mysqlParameter length of %s =%d", "DBName", len(dbConfig.DBName))
	}
	if dbConfig.Timeout.Seconds() < 1 {
		return nil, fmt.Errorf("mysqlParameter Timeout Second =%f", dbConfig.Timeout.Seconds())
	}
	if dbConfig.ReadTimeout.Seconds() < 1 {
		return nil, fmt.Errorf("mysqlParameter ReadTimeout Second =%f", dbConfig.ReadTimeout.Seconds())
	}
	if dbConfig.WriteTimeout.Seconds() < 1 {
		return nil, fmt.Errorf("mysqlParameter WriteTimeout Second =%f", dbConfig.WriteTimeout.Seconds())
	}

	conf := goMysql.NewConfig()
	conf.User = dbConfig.User
	conf.Passwd = dbConfig.Password
	conf.Net = dbConfig.Net
	conf.Addr = dbConfig.Address
	conf.DBName = dbConfig.DBName
	conf.WriteTimeout = dbConfig.WriteTimeout
	conf.ReadTimeout = dbConfig.ReadTimeout
	conf.Timeout = dbConfig.Timeout
	conf.ParseTime = dbConfig.ParseTime

	return conf, nil
}

//GetSQLDB returns real sql db operator.
//func (ms *mysql) GetDB() *sql.DB {
//	return ms.sqlDB
//}

//Close close sql db connection.
func (ms *MySQL) Close() {
	if ms.mainDB != nil {
		ms.mainDB.Close()
	}
	if ms.readDB != nil {
		ms.readDB.Close()
	}
}

//SetMainDB for sql mock test
func (ms *MySQL)SetMainDB(db *sql.DB){
	ms.mainDB=db
}

//SelectTableNames returns all table name of target db.
func (ms *MySQL) SelectTableNames(dbName string) ([]string, error) {
	res := make([]string, 0)

	var tableName string
	rows, err := ms.mainDB.Query("SELECT table_name FROM information_schema.tables where table_schema  = ?", dbName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&tableName)
		if err != nil {
			return nil, err
		}
		res = append(res, tableName)
	}
	return res, nil
}

//QueryRow executes a query that is expected to return at most one row.
// QueryRow always returns a non-nil value. Errors are deferred until Row's Scan method is called.
// If the query selects no rows, the *Row's Scan will return ErrNoRows. Otherwise,
// the *Row's Scan scans the first selected row and discards the rest.
func (ms *MySQL) QueryRow(query string, args ...interface{}) *sql.Row {

	var doDB = ms.mainDB
	if ms.readWriteSplitting && ms.readDB != nil {
		doDB = ms.readDB
	}

	if len(args) == 0 {
		return doDB.QueryRow(query)
	} else {
		return doDB.QueryRow(query, args...)
	}
}

//**************************************
func (ms *MySQL) SQLQuery(exec sqlQueryExec, query string, args ...interface{}) (retJson json.RawMessage, retErr error) {
	logPrefix := "SQLQuery() "

	var doDB = ms.mainDB
	if ms.readWriteSplitting && ms.readDB != nil {
		doDB = ms.readDB
	}

	stmt, err := doDB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			retJson = nil
			retErr = fmt.Errorf("%s stmt.Close() err %s", logPrefix, err.Error())
		}
	}()

	rows, err := stmt.Query(args...)
	if err != nil {
		retJson = nil
		retErr = fmt.Errorf("%s stmt.Query() err %s", logPrefix, err.Error())
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			retJson = nil
			retErr = fmt.Errorf("%s rows.Close() err %s", logPrefix, err.Error())
		}
	}()

	retJson, retErr = exec(rows)
	if retErr != nil {
		retJson = nil
		retErr = fmt.Errorf("%s exec err %s", logPrefix, retErr.Error())
		return
	}

	return

}

func (ms *MySQL) SQLQuery_retArrInterface(exec sqlQueryExec_retArrInterface, query string, args ...interface{}) (retJson []interface{}, retErr error) {
	logPrefix := "SQLQuery_retArrInterface() "

	var doDB = ms.mainDB
	if ms.readWriteSplitting && ms.readDB != nil {
		doDB = ms.readDB
	}

	stmt, err := doDB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			retJson = nil
			retErr = fmt.Errorf("%s stmt.Close() err %s", logPrefix, err.Error())
		}
	}()

	rows, err := stmt.Query(args...)
	if err != nil {
		retJson = nil
		retErr = fmt.Errorf("%s stmt.Query() err %s", logPrefix, err.Error())
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			retJson = nil
			retErr = fmt.Errorf("%s rows.Close() err %s", logPrefix, err.Error())
		}
	}()

	retJson, retErr = exec(rows)
	if retErr != nil {
		retJson = nil
		retErr = fmt.Errorf("%s exec err %s", logPrefix, retErr.Error())
		return
	}

	return
}

//db
//sql query row 使用在只取一個row的時候，只有一個，且一定有一個
func (ms *MySQL) SQLQueryRow(query string, args ...interface{}) (retRow *sql.Row, retErr error) {
	logPrefix := "SQLQueryRow() "

	var doDB = ms.mainDB
	if ms.readWriteSplitting && ms.readDB != nil {
		doDB = ms.readDB
	}

	stmt, err := doDB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			retRow = nil
			retErr = fmt.Errorf("%s stmt.Close() err %s", logPrefix, err.Error())
		}
	}()

	retRow = stmt.QueryRow(args...)

	return
}

func (ms *MySQL) SQLUpdate(query string, args ...interface{}) (retErr error) {
	logPrefix := "SQLUpdate() "

	stmt, err := ms.mainDB.Prepare(query)

	if err != nil {
		retErr = fmt.Errorf("%s Prepare() err %s", logPrefix, err.Error())
		return
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			retErr = fmt.Errorf("%s stmt.Close() err %s", logPrefix, err.Error())
		}
	}()

	_, err = stmt.Exec(args...)
	if err != nil {
		retErr = fmt.Errorf("%s stmt.Exec() err %s", logPrefix, err.Error())
		return
	}

	//rowCnt=0，有可能是update內容跟原本資料相同，會回傳 0
	//rowCnt, err := res.RowsAffected()
	//if err != nil {
	//	return CodeSQLExecResRowAffectedError, err
	//}
	//
	//if rowCount != -1 {
	//	if rowCnt != rowCount {
	//		return CodeSQLExecResRowAffectedCountNotMatch, err
	//	}
	//}
	retErr = nil
	return
}

// 回傳 sql.Result 的版本
func (ms *MySQL) SQLUpdate_Result(query string, args ...interface{}) (sql.Result, error) {
	logPrefix := "SQLUpdate() "

	stmt, err := ms.mainDB.Prepare(query)
	var retErr error

	if err != nil {
		retErr = fmt.Errorf("%s Prepare() errcode %s", logPrefix, err.Error())
		return nil, retErr
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			retErr = fmt.Errorf("%s stmt.Close() errcode %s", logPrefix, err.Error())
		}
	}()

	return stmt.Exec(args...)
}

func (ms *MySQL) SQLInsert(query string, args ...interface{}) (lastID int64, retErr error) {
	logPrefix := "SQLInsert() "

	stmt, err := ms.mainDB.Prepare(query)
	if err != nil {
		lastID = -1
		retErr = fmt.Errorf("%s err %s", logPrefix, err.Error())
		return
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			retErr = fmt.Errorf("%s stmt.Close() err %s", logPrefix, err.Error())
		}
	}()

	res, err := stmt.Exec(args...)
	if err != nil {
		lastID = -1
		retErr = fmt.Errorf("%s err %s", logPrefix, err.Error())
		return
	}
	lastID, err = res.LastInsertId()
	if err != nil {
		lastID = -1
		retErr = fmt.Errorf("%s err %s", logPrefix, err.Error())
		return
	}

	return lastID, nil
}

func (ms *MySQL) SQLDelete(query string, args ...interface{}) (sql.Result, error) {
	return ms.mainDB.Exec(query, args...)
}

func (ms *MySQL) SQLExec(query string, args ...interface{}) (sql.Result, error) {
	return ms.mainDB.Exec(query, args...)
}

type ITable struct {
}

type DbWhere map[string]interface{}
type DbOrderBy map[string]OrderByDir
type DbContext struct {
	db  *MySQL
	out interface{}
}
type QueryContext struct {
	//db        *MySQL
	ctx       *DbContext
	tableName string
	selects   []string
	//where     DbWhere
	where  string
	dirs   DbOrderBy
	limit  int
	offset int
	args   []interface{}
}

func (ms *MySQL) GetGetTableDefine(out interface{}) string {
	rfs := reflect.TypeOf(out)
	if rfs.Kind() != reflect.Struct {
		panic(fmt.Sprintf("out is a %s, not a struct", rfs.Kind().String()))
	}
	return strutil.Pascal2Snake(rfs.Name())
}

func (ms *MySQL) ORM(out interface{}) *DbContext {
	return &DbContext{
		db:  ms,
		out: out,
	}
}

func (ctx *DbContext) Select(s []string) *QueryContext {
	rfs := reflect.TypeOf(ctx.out)
	var tblName string
	for nf := 0; nf < rfs.NumField(); nf++ {
		f := rfs.Field(nf)

		// NOTE out struct must have a *service.ITable field
		if f.Name == reflect.TypeOf(ITable{}).Name() && f.Type.String() == "*service.ITable" {
			tblName = f.Tag.Get("name")
			if tblName == "" {
				tblName = strutil.Pascal2Snake(rfs.Name())
			}
			return &QueryContext{
				ctx:       ctx,
				tableName: tblName,
				selects:   s,
			}
		}
	}

	panic("db table definition need a *service.ITable")
}

func (ctx *DbContext) SelectWithSpecificTableName(s []string, tblName string) *QueryContext {

	return &QueryContext{
		ctx:       ctx,
		tableName: tblName,
		selects:   s,
	}
}

func (ctx *QueryContext) Join(targetTable, targetColumn, foreignColumn string) *QueryContext {
	panic("implement me")
	return ctx
}

type WhereCond struct {
	Column   string
	Cond     string
	Value    interface{}
	SubQuery bool
}

func (ctx *QueryContext) Where(fn func() (string, []interface{})) *QueryContext {
	where, args := fn()
	ctx.where = where
	ctx.args = args
	return ctx
}

func (ctx *QueryContext) _where(operator string, cond WhereCond) {
	if cond.SubQuery {
		if ctx.where == "" || operator == "" {
			ctx.where += fmt.Sprintf("%s %s %s ", cond.Column, cond.Cond, cond.Value)
		} else {
			ctx.where += fmt.Sprintf("%s %s %s %s ", operator, cond.Column, cond.Cond, cond.Value)
		}
		//ctx.where += fmt.Sprintf("%s %s %s %s ", operator, cond.Column, cond.Cond, cond.Value)
	} else {
		ctx.args = append(ctx.args, cond.Value)
		if ctx.where == "" || operator == "" {
			ctx.where += fmt.Sprintf("%s %s ? ", cond.Column, cond.Cond)
		} else {
			ctx.where += fmt.Sprintf("%s %s %s ? ", operator, cond.Column, cond.Cond)
		}
	}
}

func (ctx *QueryContext) _whereAppend(sqlStr string) {
	ctx.where += sqlStr
}

func (ctx *QueryContext) WhereAnd(cond WhereCond) *QueryContext {
	ctx._where("AND", cond)
	return ctx
}

func (ctx *QueryContext) WhereAndSubOrs(conds ...WhereCond) *QueryContext {
	ctx._whereAppend("AND ( ")
	ctx._where("", conds[0])
	for _, cond := range conds[1:] {
		ctx._where("OR", cond)
	}
	ctx._whereAppend(")")
	return ctx
}

func (ctx *QueryContext) WhereOr(cond WhereCond) *QueryContext {
	ctx._where("OR", cond)
	return ctx
}

type OrderByDir int

const (
	ASC  OrderByDir = 0
	DESC OrderByDir = 1
)

func ParseOrderByDirByInt(dir int) OrderByDir {
	if dir == 1 {
		return DESC
	} else {
		return ASC
	}
}

func (ctx *QueryContext) OrderBy(column string, dir OrderByDir) *QueryContext {
	if ctx.dirs == nil {
		ctx.dirs = map[string]OrderByDir{
			column: dir,
		}
	} else {
		ctx.dirs[column] = dir
	}
	return ctx
}

func (ctx *QueryContext) GroupBy() *QueryContext {
	panic("implement me")
	return ctx
}

func (ctx *QueryContext) Limit(limit int) *QueryContext {
	ctx.limit = limit
	return ctx
}

func (ctx *QueryContext) Offset(offset int) *QueryContext {
	ctx.offset = offset
	return ctx
}

// Final

func (ctx *QueryContext) Args() []interface{} {
	return ctx.args
}

func (ctx *QueryContext) SQL() string {
	where := ctx.where

	if where == "" {
		where = "1"
	}

	// parse order by
	orderBy := parseOrderBy(ctx.dirs)
	limit := parseLimit(ctx.limit)
	offset := parseOffset(ctx.offset)
	var selects string
	groupBy := "" // TODO

	if ctx.selects == nil || len(ctx.selects) == 0 || ctx.selects[0] == "*" {
		selects = "*"
	} else {
		selects = coll.JoinString(ctx.selects, ",")
	}

	sqlStr := fmt.Sprintf("SELECT %s FROM `%s` WHERE %s %s %s %s %s",
		selects,
		ctx.tableName,
		where,
		groupBy,
		orderBy,
		limit,
		offset)

	return strings.TrimRight(sqlStr, " ")
}

func (ctx *QueryContext) First() (interface{}, error) {
	rlt, err := ctx.Rows()
	if err != nil {
		return nil, err
	} else {
		return rlt[0], err
	}
}

func (ctx *QueryContext) Rows() ([]interface{}, error) {
	sqlStr := ctx.SQL()

	rows, err := ctx.ctx.db.mainDB.Query(sqlStr, ctx.args...)
	// fail.FailedOnError(err, "")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	fail.FailedOnError(err, "")
	columnTypes, err := rows.ColumnTypes()
	fail.FailedOnError(err, "get type failed")

	var rlt []interface{}
	for rows.Next() {
		o := reflect.New(reflect.TypeOf(ctx.ctx.out))
		scan(rows, columns, columnTypes, o.Elem())
		rlt = append(rlt, o.Elem().Interface())
	}
	return rlt, nil
}

func (ctx *QueryContext) ScanFirst(args ...interface{}) error {
	sqlStr := ctx.SQL()

	rows, err := ctx.ctx.db.mainDB.Query(sqlStr, ctx.args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()
	return rows.Scan(args...)
}

// NOTE
// 特列的 api，會用 ctx 設好的條件去算數量
func (ctx *QueryContext) Count() (int64, error) {
	where := ctx.where

	if where == "" {
		where = "1"
	}
	selects := "count(*)"

	sqlStr := fmt.Sprintf("SELECT %s FROM `%s` WHERE %s",
		selects,
		ctx.tableName,
		where)

	rows, err := ctx.ctx.db.mainDB.Query(sqlStr, ctx.args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	rows.Next()
	var count int64
	err = rows.Scan(&count)

	return count, err
}

// ---------------------------------------------------
// Prepare
type PrepareContext struct {
	db        *MySQL
	tableName string
	stmt      *sql.Stmt
	where     []string
	dirs      DbOrderBy
	limit     int
	offset    int
	out       interface{}
	args      []interface{}
}

func (ms *MySQL) Prepare(out interface{}) *PrepareContext {
	rfs := reflect.TypeOf(out)
	for nf := 0; nf < rfs.NumField(); nf++ {
		f := rfs.Field(nf)

		// NOTE out struct must have a *service.ITable field
		if f.Name == reflect.TypeOf(ITable{}).Name() && f.Type.String() == "*service.ITable" {
			return &PrepareContext{
				db:        ms,
				tableName: strutil.Pascal2Snake(rfs.Name()),
				out:       out,
			}
		}
	}
	panic("db table definition need a *service.ITable")
}

func (ctx *PrepareContext) Where(wheres []string) *PrepareContext {
	ctx.where = wheres
	return ctx
}

func (ctx *PrepareContext) OrderBy(dirs map[string]OrderByDir) *PrepareContext {
	ctx.dirs = dirs
	return ctx
}

func (ctx *PrepareContext) GroupBy() *PrepareContext {
	panic("implement me")
	return ctx
}

func (ctx *PrepareContext) Limit(limit int) *PrepareContext {
	ctx.limit = limit
	return ctx
}

func (ctx *PrepareContext) Offset(offset int) *PrepareContext {
	ctx.offset = offset
	return ctx
}

func (ctx *PrepareContext) Exec(values ...interface{}) (interface{}, error) {
	if ctx.stmt == nil {
		var where string
		for _, val := range ctx.where {
			where += val + " = ?, "
		}
		where = where[:len(where)-2]

		orderBy := parseOrderBy(ctx.dirs)
		limit := parseLimit(ctx.limit)

		sql := fmt.Sprintf("SELECT * FROM `%s` WHERE %s %s %s ",
			ctx.tableName,
			where,
			orderBy,
			limit)

		newStmt, err := ctx.db.mainDB.Prepare(sql)
		fail.FailedOnError(err, "prepare failed")
		ctx.stmt = newStmt
	}

	rows, err := ctx.stmt.Query(values...)
	fail.FailedOnError(err, "query failed")

	columns, err := rows.Columns()
	fail.FailedOnError(err, "can't get row columns")
	columnTypes, err := rows.ColumnTypes()
	fail.FailedOnError(err, "get type failed")

	var rlt []interface{}
	for rows.Next() {
		o := reflect.New(reflect.TypeOf(ctx.out))
		scan(rows, columns, columnTypes, o.Elem())
		rlt = append(rlt, o.Elem().Interface())
	}
	return rlt, nil
}

func (ctx *PrepareContext) Query(values ...interface{}) (interface{}, error) {
	if ctx.stmt == nil {
		var where string
		for _, val := range ctx.where {
			where += val + " = ?, "
		}
		where = where[:len(where)-2]

		orderBy := parseOrderBy(ctx.dirs)
		limit := parseLimit(ctx.limit)

		sql := fmt.Sprintf(`SELECT * FROM %s WHERE %s %s %s `,
			ctx.tableName,
			where,
			orderBy,
			limit)

		newStmt, err := ctx.db.mainDB.Prepare(sql)
		ctx.stmt = newStmt
		fail.FailedOnError(err, "prepare failed")
	}

	rows, err := ctx.stmt.Query(values)
	fail.FailedOnError(err, "query failed")

	columns, err := rows.Columns()
	fail.FailedOnError(err, "")
	columnTypes, err := rows.ColumnTypes()
	fail.FailedOnError(err, "get type failed")

	var rlt []interface{}
	for rows.Next() {
		o := reflect.New(reflect.TypeOf(ctx.out))
		scan(rows, columns, columnTypes, o.Elem())
		rlt = append(rlt, o.Elem().Interface())
	}
	return rlt, nil
}

func (ctx *PrepareContext) First(values ...interface{}) (interface{}, error) {
	// TODO check SELECT
	ret, err := ctx.Query(values...)
	if err == nil && ret != nil {
		return ret.([]interface{})[0], err
	} else {
		return nil, err
	}
}

func (ctx *PrepareContext) Close() {
	if ctx.stmt == nil {
		panic("no stmt")
	} else {
		err := ctx.stmt.Close()
		fail.FailedOnError(err, "close stmt failed")
	}
}

// ---------------------------------------------------
// Insert
type InsertContext struct {
	ctx       *DbContext
	tableName string
	columns   []string
	values    [][]interface{}
}

func (ctx *DbContext) Insert(fields ...string) *InsertContext {

	rfs := reflect.TypeOf(ctx.out)
	var tblName string
	for nf := 0; nf < rfs.NumField(); nf++ {
		f := rfs.Field(nf)

		// NOTE out struct must have a *service.ITable field
		if f.Name == reflect.TypeOf(ITable{}).Name() && f.Type.String() == "*service.ITable" {
			tblName = f.Tag.Get("name")
			if tblName == "" {
				tblName = strutil.Pascal2Snake(rfs.Name())
			}
			return &InsertContext{
				ctx:       ctx,
				tableName: tblName,
				columns:   fields,
			}
		}
	}

	panic("db table definition need a *service.ITable")
}

func (ctx *InsertContext) Values(values ...interface{}) *InsertContext {
	fieldsLen := len(ctx.columns)
	if len(values) != fieldsLen {
		fail.FailedOnError(errors.New("columns 和 values 數量需要相同"), "")
	}
	ctx.values = append(ctx.values, values)
	return ctx
}

// TODO move to common
func repeatString(str string, repeatCount int) []string {
	var result []string
	for i := 0; i < repeatCount; i++ {
		result = append(result, str)
	}
	return result
}

func (ctx *InsertContext) Exec() (sql.Result, error) {
	sql := fmt.Sprintf("INSERT INTO %s ", ctx.tableName)
	sql = fmt.Sprintf(sql+"(%s) ", coll.JoinString(ctx.columns, ","))

	aValueSet := "(" + coll.JoinString(repeatString("?", len(ctx.columns)), ",") + ")"

	var valuesSet []string
	var args []interface{}
	for _, data := range ctx.values {
		valuesSet = append(valuesSet, aValueSet)
		args = append(args, data...)
	}
	values := coll.JoinString(valuesSet, ",")

	sql += "VALUES " + values

	rlt, err := ctx.ctx.db.mainDB.Exec(sql, args...)
	return rlt, err
}

// ---------------------------------------------------
// Update
type UpdateContext struct {
	ctx       *DbContext
	tableName string
	columns   []string
	args      []interface{}
	where     string
	whereArgs []interface{}
}

func (ctx *DbContext) UpdateFields(data map[string]interface{}) *UpdateContext {
	newCtx := UpdateContext{}
	for col, val := range data {
		newCtx.args = append(newCtx.args, val)
		newCtx.columns = append(newCtx.columns, col)
	}
	return &newCtx
}

func (ctx *DbContext) Update() *UpdateContext {
	rfs := reflect.TypeOf(ctx.out)
	var tblName string
	for nf := 0; nf < rfs.NumField(); nf++ {
		f := rfs.Field(nf)

		// NOTE out struct must have a *service.ITable field
		if f.Name == reflect.TypeOf(ITable{}).Name() && f.Type.String() == "*service.ITable" {
			tblName = f.Tag.Get("name")
			if tblName == "" {
				tblName = strutil.Pascal2Snake(rfs.Name())
			}
			return &UpdateContext{
				ctx:       ctx,
				tableName: tblName,
			}
		}
	}
	panic("db table definition need a *service.ITable")
}

func (ctx *UpdateContext) Set(column string, value interface{}) *UpdateContext {
	rfs := reflect.TypeOf(ctx.ctx.out)
	for nf := 0; nf < rfs.NumField(); nf++ {
		f := rfs.Field(nf)

		// NOTE out struct must have a *service.ITable field
		if strutil.Pascal2Snake(f.Name) == column {
			ctx.columns = append(ctx.columns, column)
			ctx.args = append(ctx.args, value)
			return ctx
		}
	}
	panic("Unexpected column")
}

func (ctx *UpdateContext) Where(fn func() (string, []interface{})) *UpdateContext {
	where, args := fn()
	ctx.where = where
	ctx.whereArgs = args
	return ctx
}

func (ctx *UpdateContext) _where(operator string, cond WhereCond) {
	ctx.whereArgs = append(ctx.whereArgs, cond.Value)
	if ctx.where == "" || operator == "" {
		ctx.where += fmt.Sprintf("%s %s ? ", cond.Column, cond.Cond)
	} else {
		ctx.where += fmt.Sprintf("%s %s %s ? ", operator, cond.Column, cond.Cond)
	}
}

func (ctx *UpdateContext) _whereAppend(sqlStr string) {
	ctx.where += sqlStr
}

func (ctx *UpdateContext) WhereAnd(cond WhereCond) *UpdateContext {
	ctx._where("AND", cond)
	return ctx
}

func (ctx *UpdateContext) WhereAndSubOrs(conds ...WhereCond) *UpdateContext {
	ctx._whereAppend("AND ( ")
	ctx._where("", conds[0])
	for _, cond := range conds[1:] {
		ctx._where("OR", cond)
	}
	ctx._whereAppend(")")
	return ctx
}

func (ctx *UpdateContext) WhereOr(cond WhereCond) *UpdateContext {
	ctx._where("OR", cond)
	return ctx
}

func (ctx *UpdateContext) Exec() (sql.Result, error) {
	var updateFields []string
	for _, column := range ctx.columns {
		updateFields = append(updateFields, fmt.Sprintf("`%s` = ?", column))
	}

	if ctx.tableName == "" {
		panic("table name is empty!")
	}

	sqlStr := fmt.Sprintf("UPDATE `%s` SET %s WHERE %s",
		ctx.tableName,
		coll.JoinString(updateFields, ","),
		ctx.where)
	args := append(ctx.args, ctx.whereArgs...)
	return ctx.ctx.db.mainDB.Exec(sqlStr, args...)
}

func (ctx *UpdateContext) SQL() string {
	var updateFields []string
	for _, column := range ctx.columns {
		updateFields = append(updateFields, fmt.Sprintf("%s = ?", column))
	}

	return fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		ctx.tableName,
		coll.JoinString(updateFields, ","),
		ctx.where)
}

func (ctx *UpdateContext) Args() []interface{} {
	args := append(ctx.args, ctx.whereArgs...)
	return args
}

// -------------------------------------------------
// Util functions
func parseLimit(limit int) string {
	var rlt string
	if limit == 0 {
		rlt = ""
	} else {
		rlt = fmt.Sprintf("LIMIT %v", limit)
	}
	return rlt
}

func parseOffset(offset int) string {
	var rlt string
	if offset == 0 {
		rlt = ""
	} else {
		rlt = fmt.Sprintf("OFFSET %v", offset)
	}
	return rlt
}

func parseOrderBy(dirs DbOrderBy) string {
	var orderBy string
	if dirs == nil || len(dirs) == 0 {
		orderBy = ""
	} else {
		orderBy = "ORDER BY "
		for k, v := range dirs {
			var dir string
			if v == ASC {
				dir = "ASC"
			} else if v == DESC {
				dir = "DESC"
			} else {
				panic("wrong dir in field: " + k)
			}
			orderBy += fmt.Sprintf("%s %s, ", k, dir)
		}
		orderBy = orderBy[:len(orderBy)-2]
	}
	return orderBy
}

func scan(rows *sql.Rows, columns []string, columnTypes []*sql.ColumnType, out reflect.Value) {
	outType := out.Type()

	var scanArr []interface{}
	for cIdx, column := range columns {
		it := columnTypes[cIdx]
		tt := it.ScanType()

		if _, ok := outType.FieldByName(strutil.Snake2Pascal(column)); ok {
			field := out.FieldByName(strutil.Snake2Pascal(column))
			reflectValue := reflect.New(field.Type())
			scanArr = append(scanArr, reflectValue.Interface())
		} else {
			reflectValue := reflect.New(tt)
			scanArr = append(scanArr, reflectValue.Interface())
		}

	}
	err := rows.Scan(scanArr...)
	fail.FailedOnError(err, "scan failed")
	for cIdx, column := range columns {
		if _, ok := outType.FieldByName(strutil.Snake2Pascal(column)); ok {
			field := out.FieldByName(strutil.Snake2Pascal(column))
			columnValue := scanArr[cIdx]
			if field.CanSet() {
				pv := reflect.ValueOf(columnValue)
				//tv := reflect.TypeOf(columnValue)
				// NOTE 可能還有其它 type 要轉
				if pv.Elem().Type().Name() == reflect.TypeOf(sql.RawBytes{}).Name() {
					//indir := string(pv.Elem().Interface().(sql.RawBytes))
					field.SetString(string(pv.Elem().Interface().(sql.RawBytes)))
				} else {
					//NOTE 用 sqlmock 後，addRows 產出出的 row value 會多一層 *interface{} 造成 test 失敗
					if pv.Kind().String() == "ptr" {
						iv := reflect.Indirect(pv)
						if iv.Kind() == reflect.Interface {
							// sqlmock 回傳是 *interface
							field.Set(iv.Elem())
						} else {
							field.Set(pv.Elem())
						}
					} else {
						// 非 pointer 有問題, db row 出來應該要是 ptr
						// 所以這情況可以不用寫(不應該發生)
					}
				}
			}

		}
	}
}

//Transaction
