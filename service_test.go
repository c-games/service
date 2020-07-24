package service

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestService(t *testing.T) {
	wd, _ := os.Getwd()
	for !strings.HasSuffix(wd, "service") {
		wd = filepath.Dir(wd)
	}
	//t.Logf("wd = %s", wd)

	app, err := NewService(wd + "/serviceConfig.json")
	if err != nil {
		t.Fatalf("err=%s\n", err.Error())
	}

	t.Logf("logger %+v", app.configure.Logger)
	t.Logf("mq %+v", app.configure.MQ)
	t.Logf("http %+v", app.configure.HTTP)
	t.Logf("redis %+v", app.configure.Redis)
	t.Logf("mysql %+v", app.configure.Mysql)

	//show mysql stat

	mysql, _ := app.GetMySQL()
	t.Logf("DB stat %+v", mysql.mainDB.Stats())

	//test publish wallet
	//num := 5
	//testChan := make(chan []byte)
	//
	//for i := num; i > 0; i-- {
	//	err = testPublishProcdureWalletQuery(app, testChan)
	//	if err != nil {
	//		t.Errorf("PublishProcedure error %s", err.Error())
	//		return
	//	}
	//}
	//
	//for {
	//	b := <-testChan
	//	t.Logf("testChan %s", string(b))
	//	num -= 1
	//	if num == 0 {
	//		return
	//	}
	//}

	app.Hold()

}

func testPublishQueryUserLimit(app *Service, testChan chan []byte) error {
	pubData := []publishUserQueryUserLimitData{{1000001, 201}}

	bd, _ := json.Marshal(&pubData)
	sn := randomSerial()
	body, _ := getCGMessage(sn, "query_user_limit", "cruise-test", bd, true)
	pp := &PublishParam{
		"cg-user",
		body,
		testChan,
	}

	return app.mq.PublishProcedure(pp)
}

//todo 這個是publish 不是 publishProcedure
func testPublishWalletQuery(app *Service, testChan chan []byte) error {
	pubData := []publishWalletQueryData{{101, 1002, 1000001}}

	bd, _ := json.Marshal(&pubData)
	sn := randomSerial()
	body, _ := getCGMessage(sn, "query", "cruise-test", bd, true)

	return app.mq.Publish("cg-wallet", body)
}

func testPublishProcdureWalletQuery(app *Service, testChan chan []byte) error {
	pubData := []publishWalletQueryData{{101, 1002, 1000001}}

	bd, _ := json.Marshal(&pubData)
	sn := randomSerial()
	body, _ := getCGMessage(sn, "query", "cruise-test", bd, true)
	pp := &PublishParam{
		"cg-wallet",
		body,
		testChan,
	}

	return app.mq.PublishProcedure(pp)
}

type publishWalletQueryData struct {
	MasterAgentID int   `json:"master_agent_id"`
	AgentID       int   `json:"agent_id"`
	UserID        int64 `json:"user_id"`
}

type publishUserQueryUserLimitData struct {
	UserID int64 `json:"user_id"`
	GameID int   `json:"game_id"`
}

func getCGMessage(serial int, cmd string, responseQueue string, data []byte, waitResponse bool) ([]byte, error) {

	res := &CGMessage{
		Serial:        int64(serial),
		ResponseQueue: responseQueue,
		Command:       cmd,
		Data:          json.RawMessage(data),
		WaitResponse:  waitResponse,
	}

	return json.Marshal(res)
}

func Testgg(t *testing.T) {

	type st struct {
		Code string `json:code`
		Info string `json:info`
	}

	a := `{"code":"1003","info":""}`
	aa := []byte(a)
	ss := &st{}

	if err := json.Unmarshal(aa, ss); err != nil {
		panic(err)
	}
}
