package service

import (
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/proullon/ramsql/driver"
	"testing"
	"time"
)

func TestUpdateUserAddress(t *testing.T) {
	batch := []string{
		`CREATE TABLE address (id BIGSERIAL PRIMARY KEY, street TEXT, street_number INT);`,
		`CREATE TABLE user_addresses (address_id INT, user_id INT);`,
		`INSERT INTO address (street, street_number) VALUES ('rue Victor Hugo', 32);`,
		`INSERT INTO user_addresses (address_id, user_id) VALUES (2, 1);`,
	}

	d, _ := time.ParseDuration("10s")
	ms := &MysqlParameter{
		"ramsql",
		"whatever",
		"whatever",
		"whatever",
		"whatever",
		"address",
		d,
		d,
		d,
		false,
	}
	mysql, err := NewMySQL(ms, false, nil)

	if err != nil {
		t.Fatalf("sql.Open : Error : %s\n", err)
	}

	//create ramsql
	for _, b := range batch {
		_,err=mysql.SQLExec(b)
		if err != nil {
			t.Fatalf("sql.Exec: Error: %s\n", err)
		}
	}

	err=mysql.SQLUpdate("update address set street =? where street_number = ?","street",32,)

	if err != nil {
		t.Fatalf("Too bad! unexpected error: %s", err)
	}
}

func TestORM(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Create mock failed")
	}

	mysql := &MySQL{mainDB: db, readWriteSplitting: false}

	type OOO struct {
		*ITable `name:"ooo_table"`
		Id      int64 `sql:"id"`
	}

	// basic query
	mock.ExpectQuery("SELECT (.+) FROM `ooo_table` WHERE 1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "test2"}).AddRow(10000001, "Don't assign this"))

	ret, err := mysql.ORM(OOO{}).Select(nil).Rows()

	if err != nil {
		t.Error(err.Error())
	}
	if ret[0].(OOO).Id != 10000001 {
		t.Errorf("NOT equal")
	}

	// with where
	mock.ExpectQuery("SELECT (.+) FROM `ooo_table` WHERE id = \\? AND test > \\?").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))

	ret, err = mysql.ORM(OOO{}).Select(nil).
		Where(func() (string, []interface{}) {
			return "id = ? AND test > ?", []interface{}{10, "test"}
		}).Rows()

	if err != nil {
		t.Error(err.Error())
	}
	if ret[0].(OOO).Id != 10 {
		t.Errorf("NOT equal")
	}

	//
	mock.ExpectQuery("SELECT id,test FROM `ooo_table` WHERE id = \\? AND test > \\? ORDER BY test DESC LIMIT 10").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))

	ret, err = mysql.ORM(OOO{}).Select([]string{"id", "test"}).
		WhereAnd(WhereCond{Column: "id", Cond: "=", Value: 10}).
		WhereAnd(WhereCond{Column: "test", Cond: ">", Value: 10}).
		OrderBy("test", DESC).
		Limit(10).
		Rows()
	if err != nil {
		t.Error(err.Error())
	}
	if ret[0].(OOO).Id != 10 {
		t.Errorf("NOT equal")
	}

}

func TestORM_insert(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatal("Create mock failed")
	}

	mysql := &MySQL{mainDB: db, readWriteSplitting: false}

	type OOO struct {
		*ITable `name:"ooo_table"`
		Id      int64 `sql:"id"`
	}

	// basic query
	mock.ExpectExec("INSERT INTO ooo_table (id,name,email) VALUES (?,?,?),(?,?,?)").
		WithArgs(10, "Bob", "test@123.tw", 11, "Any", "any@123.tw").
		WillReturnResult(sqlmock.NewResult(2, 2))

	_, err = mysql.ORM(OOO{}).Insert("id", "name", "email").
		Values(10, "Bob", "test@123.tw").
		Values(11, "Any", "any@123.tw").
		Exec()

	if err != nil {
		t.Error(err.Error())
	}
}
