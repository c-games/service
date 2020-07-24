package service

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//Lottery is lottery server framework.
type Service struct {
	configure    *Configure
	mysql        *MySQL
	logger       *Logger
	mq           *MQ
	http         *HTTPServer
	redis        *Redis
	idGenerators map[int64]*idGenerator
	done         chan int
}

//NewService returns a app structure with config file path passed in.
//Config file includes mysql config.
func NewService(configFilePath string) (*Service, error) {

	//new configure
	conf, err := newConfigure(configFilePath)
	if err != nil {
		return nil, err
	}

	//new logger
	lg, err := newLogger(conf.Logger.FileName, conf.Logger.Level, conf.Logger.Address, conf.Environment,conf.Service)
	if err != nil {
		return nil, err
	}

	lt := &Service{
		logger:    lg,
		configure: conf,
		mysql:     nil,
		mq:        nil,
		http:      nil,
		done:      make(chan int),
	}

	if err = lt.enable(); err != nil {
		return nil, err
	}

	return lt, nil
}

func (s *Service) handleCtrlC() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		s.done <- 0
	}()
}

//Hold 如果主程式沒有hold住，就呼叫這個，主程式有就不需要。
func (s *Service) Hold() {
	s.handleCtrlC()

	defer func() {
		if r := recover(); r != nil {

		}
	}()

	for {
		select {
		//ctrl-c to shtdown
		case <-s.done:
			os.Exit(0)
		}
	}
}
func (s *Service) enable() error {

	idGen, err := newIDGenerator(0)
	if err != nil {
		return fmt.Errorf("idGenerator error %s", err.Error())
	}

	s.idGenerators = make(map[int64]*idGenerator)
	s.idGenerators[0] = idGen

	//if mySQL enable
	if s.configure.Mysql.Enable {
		if mysql, err := s.enableMySQL(); err != nil {
			return fmt.Errorf("MySQL %s", err.Error())
		} else {
			s.mysql = mysql
		}
	}

	//if mq enable
	if s.configure.MQ.Enable && len(s.configure.MQ.QueueDefinition) > 0 {
		if mq, err := s.enableMQ(); err != nil {
			return fmt.Errorf("MQ %s", err.Error())
		} else {
			s.mq = mq
		}
	}

	//if http enable
	if s.configure.HTTP.Enable {
		if hp, err := s.enableHTTP(); err != nil {
			return fmt.Errorf("HTTP %s", err.Error())
		} else {
			s.http = hp
			if err = s.http.Start(); err != nil {
				return fmt.Errorf("HTTP %s", err.Error())
			}
		}
	}

	//if redis enable
	if s.configure.Redis.Enable {
		if re, err := s.enableRedis(); err != nil {
			return fmt.Errorf("Redis %s", err.Error())
		} else {
			s.redis = re
		}
	}

	return nil
}
func (s *Service) enableMySQL() (*MySQL, error) {

	if s.configure.Mysql == nil {
		return nil, errors.New("mysql configure is not initialed")
	}
	mainDBConf := s.configure.Mysql.MainDB
	tm, _ := time.ParseDuration(mainDBConf.Timeout)
	rm, _ := time.ParseDuration(mainDBConf.ReadTimeout)
	wm, _ := time.ParseDuration(mainDBConf.WriteTimeout)

	mainDBParam := &MysqlParameter{
		DriverName:   mainDBConf.DriverName,
		User:         mainDBConf.User,
		Password:     mainDBConf.Password,
		Net:          mainDBConf.Net,
		Address:      mainDBConf.Address,
		DBName:       mainDBConf.DBName,
		Timeout:      tm,
		ReadTimeout:  rm,
		WriteTimeout: wm,
		ParseTime:    mainDBConf.ParseTime,
	}

	var readOnlyDBParam *MysqlParameter

	//有讀寫分離
	if s.configure.Mysql.ReadWriteSplitting {
		wDBConf := s.configure.Mysql.ReadOnlyDB
		tm, _ := time.ParseDuration(wDBConf.Timeout)
		rm, _ := time.ParseDuration(wDBConf.ReadTimeout)
		wm, _ := time.ParseDuration(wDBConf.WriteTimeout)

		readOnlyDBParam = &MysqlParameter{
			DriverName:   mainDBConf.DriverName,
			User:         mainDBConf.User,
			Password:     mainDBConf.Password,
			Net:          mainDBConf.Net,
			Address:      mainDBConf.Address,
			DBName:       mainDBConf.DBName,
			Timeout:      tm,
			ReadTimeout:  rm,
			WriteTimeout: wm,
			ParseTime:    mainDBConf.ParseTime,
		}
	}

	return NewMySQL(mainDBParam, s.configure.Mysql.ReadWriteSplitting, readOnlyDBParam)
}
func (s *Service) enableMQ() (*MQ, error) {
	if s.configure.MQ == nil {
		return nil, errors.New("mq configure is not initialed")
	}
	mqConf := s.configure.MQ

	return newMQ(mqConf.Url, mqConf, s.logger)
}

func (s *Service) enableHTTP() (*HTTPServer, error) {
	if s.configure.HTTP == nil {
		return nil, errors.New("http configure is not initialed")
	}
	httpConf := s.configure.HTTP

	httpParams := &httpParameter{
		Address:         httpConf.Address,
		Port:            httpConf.Port,
		ReadTimeout:     time.Second * time.Duration(httpConf.ReadTimeoutSecond),
		WriteTimeout:    time.Second * time.Duration(httpConf.WriteTimeoutSecond),
		IdleTimeout:     time.Second * time.Duration(httpConf.IdleTimeoutSecond),
		MaxHeaderBytes:  httpConf.MaxHeaderBytes,
		IsTLS:           httpConf.IsTLS,
		CertificateFile: httpConf.CertificateFile,
		KeyFile:         httpConf.KeyFile,
	}

	return newHTTP(httpParams)
}

func (s *Service) enableRedis() (*Redis, error) {
	if s.configure.Redis == nil {
		return nil, errors.New("redis configure is not initialed")
	}
	conf := s.configure.Redis

	params := &RedisParameter{
		Network:      conf.Network,
		Address:      conf.Address,
		Password:     conf.Password,
		DB:           conf.DB,
		DialTimeout:  time.Second * time.Duration(conf.DialTimeoutSecond),
		ReadTimeout:  time.Second * time.Duration(conf.ReadTimeoutSecond),
		WriteTimeout: time.Second * time.Duration(conf.WriteTimeoutSecond),
		PoolSize:     conf.PoolSize,
	}

	return NewRedis(params)
}

func (s *Service) GetConfigure() (*Configure, error) {
	return s.configure, nil
}

func (s *Service) GetMySQL() (*MySQL, error) {
	if s.mysql != nil {
		return s.mysql, nil
	}
	return nil, errors.New("mysql is not enabled")
}
func (s *Service) GetMQ() (*MQ, error) {
	if s.mq != nil {
		return s.mq, nil
	}
	return nil, errors.New("mq is not enabled")
}
func (s *Service) GetLogger() (*Logger, error) {
	if s.logger != nil {
		return s.logger, nil
	}
	return nil, errors.New("logger is not enabled")
}
func (s *Service) GetHTTP() (*HTTPServer, error) {
	if s.http != nil {
		return s.http, nil
	}
	return nil, errors.New("http is not enabled")
}

func (s *Service) GetRedis() (*Redis, error) {
	if s.redis != nil {
		return s.redis, nil
	}
	return nil, errors.New("redis is not enabled")
}

func (s *Service) GetIDGenerator(nodeNum int64) (*idGenerator, error) {
	if nodeNum < 0 || nodeNum > 1023 {
		return nil, fmt.Errorf("nodeNum out of range want 0-1023 got %d", nodeNum)
	}

	if v, ok := s.idGenerators[nodeNum]; ok {
		return v, nil
	}
	//create gen
	g, err := newIDGenerator(nodeNum)

	if err != nil {
		return nil, fmt.Errorf("get idGenerator error nodeNum %d", nodeNum)
	}

	s.idGenerators[nodeNum] = g

	return s.idGenerators[nodeNum], nil

}
