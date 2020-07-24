package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

//configure is a operator for default config
//which load data from config file.
//Use newConfigure func to create configure.
type Configure struct {
	Version     string
	Environment string
	Service     string
	Logger      *loggerConfig
	Mysql       *mysqlConfig
	MQ          *mqConfig
	HTTP        *httpConfig
	Redis       *redisConfig
	Others      map[string]json.RawMessage
	Websocket   *websocketConfig
}

//type lotteryConfig struct {
//	ID int
//}

type loggerConfig struct {
	Enable   bool // true=enable, false=disable
	FullPath bool // true=log with full path, false= log with short path
	FileName string
	Level    string //default log level:"debug","info","warning","error","panic"
	Address  string // host:port
}

//請看 mysql mysqlParameter
type mysqlConfig struct {
	Enable             bool
	MainDB             mysqlDBParam //Main DB setting，沒讀寫分離，就用這個。如果有讀寫分離，這個就是Read only
	ReadWriteSplitting bool         //是否讀寫分離，有讀寫分離，WriteOnly一定要設定，否則錯。
	ReadOnlyDB         mysqlDBParam //有毒蠍分離，這是Write only DB setting
}

type mysqlDBParam struct {
	DriverName     string
	DataSourceName string
	User           string
	Password       string
	Net            string
	Address        string
	DBName         string
	Timeout        string //string, 在EnableMySQL時改duration
	ReadTimeout    string //string, 在EnableMySQL時改duration
	WriteTimeout   string //string, 在EnableMySQL時改duration
	ParseTime      bool
}

type mqConfig struct {
	Enable                bool
	Url                   string
	ResponseTimeoutSecond int
	ResponseConsumerParam struct {
		Queue     string // queue name，使用 service name
		Consumer  string // consumer name
		AutoAck   bool   // consumer是否自動送回ack(刪除該訊息)
		Exclusive bool   //
		NoLocal   bool   // no-local
		NoWait    bool   // no-wait
		Args      map[string]interface{}
	}
	CommandConsumerParam struct {
		Queue     string // target queue name
		Consumer  string // consumer name
		AutoAck   bool   // consumer是否自動送回ack(刪除該訊息)
		Exclusive bool   //
		NoLocal   bool   // no-local
		NoWait    bool   // no-wait
		Args      map[string]interface{}
	}
	Qos struct {
		PrefetchCount int
		PrefetchSize  int
		Global        bool
	}
	//queue define
	//Type 0=consume command, 1= consume response,2 = other
	QueueDefinition []struct {
		Name       string
		Durable    bool
		AutoDelete bool
		Exclusive  bool
		NoWaite    bool
		Args       map[string]interface{}
	}
}

type httpConfig struct {
	Enable             bool
	Address            string
	Port               string
	ReadTimeoutSecond  int
	WriteTimeoutSecond int
	IdleTimeoutSecond  int
	MaxHeaderBytes     int
	IsTLS              bool
	CertificateFile    string
	KeyFile            string
}

type redisConfig struct {
	Enable             bool
	Network            string
	Address            string
	Password           string
	DB                 int
	DialTimeoutSecond  int
	ReadTimeoutSecond  int
	WriteTimeoutSecond int
	PoolSize           int
}
type websocketConfig struct {
	Enable          bool
	Address         string
	Port            string
	ConnPoolSize    int
	ChanPoolSize    int
	AcceptTimeout   string //3s
	AliveTimeout    string //eg. 2h3m
	ReadBufferSize  int
	WriteBufferSize int
}

//newConfigure returns configure.
func newConfigure(fileName string) (*Configure, error) {

	configFilePath := getFilePosition(fileName)

	cf := &Configure{
		//data: configData{},
	}

	b, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	//unmarshal to struct
	if err := json.Unmarshal(b, &cf); err != nil {
		return nil, err
	}

	return cf, nil
}

func (c *Configure) loadConfig(filePath string) error {

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	//unmarshal to struct
	if err := json.Unmarshal(b, &c); err != nil {
		return err
	}

	return nil
}

func (c *Configure) GetString(key string) (rlt string, err error) {
	if v, ok := c.Others[key]; ok {
		err = json.Unmarshal(v, &rlt)
		return
	}
	err = errors.New("unknown key")
	return
}

func (c *Configure) GetInt64(key string) (rlt int64, err error) {
	if v, ok := c.Others[key]; ok {
		err = json.Unmarshal(v, &rlt)
		return
	}
	err = errors.New("unknown key")
	return
}

func (c *Configure) GetFloat64(key string) (rlt float64, err error) {
	if v, ok := c.Others[key]; ok {
		err = json.Unmarshal(v, &rlt)
		return
	}
	err = errors.New("unknown key")
	return
}

func (c *Configure) GetArray(key string, array interface{}) (err error) {
	if v, ok := c.Others[key]; ok {
		err = json.Unmarshal(v, &array)
		return
	}
	err = errors.New("unknown key")
	return
}

func (c *Configure) GetMap(key string, s interface{}) (err error) {
	if v, ok := c.Others[key]; ok {
		err = json.Unmarshal(v, &s)
		return
	}
	err = errors.New("unknown key")
	return
}

func getFilePosition(fileName string) string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
	}

	var buf bytes.Buffer
	if fileName[0] == '/' {
		buf.WriteString(fileName)
	} else {
		buf.WriteString(dir)
		buf.WriteString("/")
		buf.WriteString(fileName)
	}

	return buf.String()
}
