package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"sync"
	"time"
)

const (
	CGResMessageErrorCodeTimeout   int = -1
	CGResMessageErrorCodeChanError int = -2
)

type MQ struct {
	url                 string
	mqConf              *mqConfig
	connection          *amqp.Connection
	channel             *amqp.Channel
	consumeCommandChan  <-chan amqp.Delivery //consume command to me
	consumeResponseChan <-chan amqp.Delivery //consume response message to me

	//receivers  client使用publishProcedure 一起送來自己的接收chan，以serial當key放在這裡。
	//client使用publishProcedure 一起送來自己的接收chan，以serial當key放在這裡。
	//收到response時，用serial找出對應的chan，把CGResMessage傳下去
	//sync.Map 的key,value都是interface{}
	receivers sync.Map

	logger          *Logger
	mqHandlerRouter *MQHandlerRouter
	timeoutSecond   int

	notifyConnectionClose chan *amqp.Error
	notifyConnectionBlock chan amqp.Blocking
	notifyChannelClose    chan *amqp.Error
	notifyPublishConfirm  chan amqp.Confirmation
}

type QueueDefinition struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWaite    bool
	Args       map[string]interface{}
}

//CGMessage is structure for consume command message
type CGMessage struct {
	Serial        int64             `json:"serial"`
	ResponseQueue string          `json:"response_queue"`
	Command       string          `json:"command"`
	Data          json.RawMessage `json:"data"`
	WaitResponse  bool            `json:"wait_response"`
}

//CGResMessage is struct for consume response message
type CGResMessage struct {
	Serial       int64           `json:"serial"`
	Command      string          `json:"command"`
	ErrorMessage string          `json:"error_message"`
	ErrorCode    int             `json:"error_code"`
	Data         json.RawMessage `json:"data"`
}

//PublishParam 是傳送的資料，
//包括實際傳送的資料 []byte
//發送者收response的chan
//mq紀錄誰送的sn
type PublishParam struct {
	QueueName      string
	PublishingBody []byte
	ReceiveChan    chan []byte
}

type receiver struct {
	Serial      int64
	ReceiveChan chan []byte
	beginTime   time.Time
}

//Adder is a type with
//adds new function to MQ
type Adder interface {
	ServeFunc()
}

//queueDefine is a slice contains QueueDefine
//these queues will be declare in InitAllProcedure
func newMQ(url string, mqCon *mqConfig, lg *Logger) (*MQ, error) {

	mh := &MQHandlerRouter{
		handlers: sync.Map{},
	}
	m := &MQ{
		url:             url,
		logger:          lg,
		mqHandlerRouter: mh,
		timeoutSecond:   mqCon.ResponseTimeoutSecond,

		notifyConnectionClose: make(chan *amqp.Error),
		notifyConnectionBlock: make(chan amqp.Blocking),
		notifyChannelClose:    make(chan *amqp.Error),
	}

	m.mqConf = mqCon

	//init all
	_, err := m.InitAllProcedure()
	if err != nil {
		return nil, err
	}

	go m.receiveNotifyConnectionClose()
	go m.receiveNotifyConnectionBlock()
	go m.receiveNotifyChannelClose()
	go m.consumeCommandMessage()
	go m.consumeResponseMessage()

	return m, nil
}

//*****************************
//AddQueue adds queue name and declares with default params
//before InitAllProcedure

//InitAllProcedure init all
//must QueueDefined() first
func (m *MQ) InitAllProcedure() (retMsg string, retErr error) {
	//init connection
	err := m.initConnection()
	if err != nil {
		retMsg = fmt.Sprintf("MQ InitAllProcedure()  err %s", err.Error())
		retErr = errors.New(retMsg)
		m.logger.LogFile(LevelError, retMsg)
		return
	}

	//init channel
	err = m.initChannel()
	if err != nil {
		retMsg = fmt.Sprintf("MQ InitAllProcedure()  err %s", err.Error())
		retErr = errors.New(retMsg)
		m.logger.LogFile(LevelError, retMsg)
		return
	}

	//init queue
	err = m.initQueue()
	if err != nil {
		retMsg = fmt.Sprintf("MQ InitAllProcedure() err %s", err.Error())
		retErr = errors.New(retMsg)
		m.logger.LogFile(LevelError, retMsg)
	}
	return
}

//*****************************
//connection
func (m *MQ) initConnection() error {

	defer func() {
		if r := recover(); r != nil {
			errMsg := fmt.Sprintf("MQ initConnection() panic %s", r)
			m.logger.LogFile(LevelPanic, errMsg)
		}
	}()

	if err := m.CloseConnection(); err != nil {
		errMsg := fmt.Sprintf("MQ initConnection CloseConnection() error %s", err.Error())
		m.logger.LogFile(LevelError, errMsg)
	}

	var conn *amqp.Connection
	var err error
	var i = 0
	for i = 0; i < 3; i++ {

		conn, err = amqp.Dial(m.url)

		if err != nil {
			//sleep
			time.Sleep(time.Second * 1)
		} else {
			break
		}
	}
	if err != nil {
		return fmt.Errorf("initConnection() err %s", err.Error())
	}

	m.connection = conn
	m.notifyConnectionClose = make(chan *amqp.Error)
	m.notifyConnectionBlock = make(chan amqp.Blocking)
	m.connection.NotifyClose(m.notifyConnectionClose)
	m.connection.NotifyBlocked(m.notifyConnectionBlock)
	return nil
}

func (m *MQ) receiveNotifyConnectionClose() {
	defer func() {
		if r := recover(); r != nil {
			errMsg := fmt.Sprintf("MQ receiveConnectionClosedChan() panic %s", r)
			m.logger.LogFile(LevelPanic, errMsg)
		}
	}()

	for {
		e, ok := <-m.notifyConnectionClose
		if ok {
			errMsg := fmt.Sprintf("MQ receiveConnectionClosedChan() connection closed code %d reason %s server %v recover %v", e.Code, e.Reason, e.Server, e.Recover)
			m.logger.LogFile(LevelError, errMsg)

			//全部init
			m.InitAllProcedure()
		}
	}
}
func (m *MQ) receiveNotifyConnectionBlock() {
	//todo block 要停住，等解除再開始

	defer func() {
		if r := recover(); r != nil {
			errMsg := fmt.Sprintf("MQ receiveNotifyConnectionBlock() panic %s", r)
			m.logger.LogFile(LevelPanic, errMsg)
		}
	}()

	for {
		b, ok := <-m.notifyConnectionBlock

		if ok {
			if b.Active {
				errMsg := fmt.Sprintf("MQ receiveNotifyConnectionBlock() connection blocking active reason %s", b.Reason)
				m.logger.LogFile(LevelError, errMsg)
			}
		}
	}
}

//*****************************
//channel
func (m *MQ) initChannelProcedure() {

	err := m.initChannel()
	if err != nil {
		errMsg := fmt.Sprintf("MQ initChannelProcedure()  err %s", err.Error())
		m.logger.LogFile(LevelError, errMsg)
	}

	err = m.initQueue()
	if err != nil {
		errMsg := fmt.Sprintf("MQ initChannelProcedure()  err %s", err.Error())
		m.logger.LogFile(LevelError, errMsg)
	}
}

//channel
func (m *MQ) initChannel() error {

	defer func() {
		if r := recover(); r != nil {
			errMsg := fmt.Sprintf("MQ initChannel() panic %s", r)
			m.logger.LogFile(LevelPanic, errMsg)
		}
	}()

	if err := m.CloseChannel(); err != nil {
		errMsg := fmt.Sprintf("MQ initChannel closeChannel() error %s", err.Error())
		m.logger.LogFile(LevelError, errMsg)
	}

	channel, err := m.connection.Channel()
	if err != nil {
		return fmt.Errorf("initChannel() err %s", err.Error())
	}
	m.channel = channel
	m.notifyChannelClose = make(chan *amqp.Error)
	m.notifyPublishConfirm = make(chan amqp.Confirmation)
	m.channel.NotifyClose(m.notifyChannelClose)
	m.channel.NotifyPublish(m.notifyPublishConfirm)

	return nil
}

func (m *MQ) receiveNotifyChannelClose() {
	defer func() {
		if r := recover(); r != nil {
			errMsg := fmt.Sprintf("MQ receiveChannelClosedChan() panic %s", r)
			m.logger.LogFile(LevelPanic, errMsg)
		}
	}()

	for {
		e, ok := <-m.notifyChannelClose
		if ok {
			errMsg := fmt.Sprintf("MQ receiveChannelClosedChan() channel closed code %d reason %s server %v recover %v", e.Code, e.Reason, e.Server, e.Recover)
			m.logger.LogFile(LevelError, errMsg)

			//init channel
			m.initChannelProcedure()
		}
	}
}

//*****************************

func (m *MQ) CloseConnection() error {
	if m.connection != nil {
		return m.connection.Close()
	}
	return nil
}

func (m *MQ) CloseChannel() error {
	if m.channel != nil {
		return m.channel.Close()
	}
	return nil
}

//*****************************
func (m *MQ) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {

	return m.channel.QueueDeclare(
		name,       // name
		durable,    // durable
		autoDelete, // delete when unused
		exclusive,  // exclusive
		noWait,     // no-wait
		args,       // arguments
	)
}

func (m *MQ) queueDeclare(chann *amqp.Channel, name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {

	return chann.QueueDeclare(
		name,       // name
		durable,    // durable
		autoDelete, // delete when unused
		exclusive,  // exclusive
		noWait,     // no-wait
		args,       // arguments
	)
}

//func (m *MQ) consumeDeclare(queueName, consumerName string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
//
//	return m.channel.Consume(
//		queueName,    // queue
//		consumerName, // consumer
//		autoAck,      // auto-ack
//		exclusive,    // exclusive
//		noLocal,      // no-local
//		noWait,       // no-wait
//		args,         // args
//	)
//}

func (m *MQ) qos(channel *amqp.Channel, prefetchCount, prefetchSize int, global bool) error {
	return channel.Qos(prefetchCount, prefetchSize, global)
}

func (m *MQ) initQueue() error {

	mqCon := m.mqConf
	channel := m.channel

	err := m.qos(channel, mqCon.Qos.PrefetchCount, mqCon.Qos.PrefetchSize, mqCon.Qos.Global)
	if err != nil {
		return err
	}

	//init queue
	for _, v := range m.mqConf.QueueDefinition {
		_, err := m.queueDeclare(channel, v.Name, v.Durable, v.AutoDelete, v.Exclusive, v.NoWaite, v.Args)
		if err != nil {
			return err
		}
	}

	//init command consumer
	if m.consumeCommandChan, err = channel.Consume(
		m.mqConf.CommandConsumerParam.Queue,     // target queue
		m.mqConf.CommandConsumerParam.Consumer,  // consumer
		m.mqConf.CommandConsumerParam.AutoAck,   // auto-ack
		m.mqConf.CommandConsumerParam.Exclusive, // exclusive
		m.mqConf.CommandConsumerParam.NoLocal,   // no-local
		m.mqConf.CommandConsumerParam.NoWait,    // no-wait
		m.mqConf.CommandConsumerParam.Args,      // args
	); err != nil {
		return err
	}

	//init response consumer
	if m.consumeResponseChan, err = channel.Consume(
		m.mqConf.ResponseConsumerParam.Queue,     // target queue
		m.mqConf.ResponseConsumerParam.Consumer,  // consumer
		m.mqConf.ResponseConsumerParam.AutoAck,   // auto-ack
		m.mqConf.ResponseConsumerParam.Exclusive, // exclusive
		m.mqConf.ResponseConsumerParam.NoLocal,   // no-local
		m.mqConf.ResponseConsumerParam.NoWait,    // no-wait
		m.mqConf.ResponseConsumerParam.Args,      // args
	); err != nil {
		return err
	}

	return nil
}

func (m *MQ) consumeCommandMessage() {
	defer func() {
		if r := recover(); r != nil {
			errMsg := fmt.Sprintf("MQ consumeCommandMessage() panic %s", r)
			m.logger.LogFile(LevelPanic, errMsg)
		}
	}()

	for {
		d, ok := <-m.consumeCommandChan

		if ok {

			//chan closed
			if !ok {
				log.Printf("" +
					"MQ consumeCommandMessage() chan not ok")
				continue
			}

			msg := &CGMessage{}
			err := json.Unmarshal(d.Body, msg)
			if err != nil {
				log.Printf("MQ consumeCommandMessage() error %s msg %v", err.Error(), msg)
			} else {
				if v, ok := m.mqHandlerRouter.getHandler(msg.Command); ok {
					go m.mqHandlerRouter.runHandler(msg, v, func() {
						err = d.Ack(false) //delete message
						if err != nil {
							errMsg := fmt.Sprintf("MQ consumeCommandMessage() d.ACK() err %s", err.Error())
							m.logger.LogFile(LevelError, errMsg)
						}
					})
				}
			}

		}
	}
}

//Publish send bode(message) to target queue.
//This is work queue model. exchange is not needed.
func (m *MQ) Publish(queueName string, body []byte) error {

	//chann:=m.channel

	defer func() {
		//if err := chann.Close(); err != nil {
		//
		//	errMsg := fmt.Sprintf(fmt.Sprintf("MQ PublishProcedure() chann close error %s", err.Error()))
		//	m.logger.LogFile(LevelError, errMsg)
		//}

		if r := recover(); r != nil {
			m.logger.LogFile(LevelPanic, fmt.Sprintf("MQ Publish() panic %s", r))
		}
	}()

	return m.channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})
}

//PublishProcedure 處理 publish,
//紀錄publish者的 chan
//處理timeout
//處理收到後傳給chan

var sequence uint8 = 0
var pid int = os.Getpid()

func (m *MQ) PublishProcedure(param *PublishParam) (retErr error) {

	//chann:=m.channel

	defer func() {
		/*
			if err := chann.Close(); err != nil {
				errMsg := fmt.Sprintf(fmt.Sprintf("MQ PublishProcedure() chann close error %s", err.Error()))
				m.logger.LogFile(LevelError, errMsg)
				retErr = fmt.Errorf(errMsg)

			}
		*/
		if r := recover(); r != nil {
			errMsg := fmt.Sprintf("MQ PublishProcedure() panic %s", r)
			m.logger.LogFile(LevelPanic, errMsg)
			retErr = fmt.Errorf(errMsg)
		}
	}()

	cm := &CGMessage{}
	if err := json.Unmarshal(param.PublishingBody, cm); err != nil {
		errMsg := fmt.Sprintf("MQ PublishProcedure() unmarshal body error %s", err.Error())
		m.logger.LogFile(LevelError, errMsg)
		retErr = fmt.Errorf(errMsg)
		return
	}

	timestamp := time.Now().UnixNano()
	serial := timestamp | int64(pid)<<22
	serial |= int64(sequence)
	sequence += 1

	if cm.WaitResponse {
		if ok := m.storeReceiver(serial, &receiver{serial, param.ReceiveChan, GetNowWithDefaultLocation()}); !ok {
			errMsg := fmt.Sprintf("MQ PublishProcedure() putReceiveChan close error serial already used %d queue %s", serial, param.QueueName)
			m.logger.LogFile(LevelError, errMsg)
			return fmt.Errorf(errMsg)
		}
	}
	//todo 可刪
	//rcv:=m.loadReceiver(param.Serial)
	//fmt.Sprintf("mq receiver %v \n",rcv)

	//publish
	err := m.channel.Publish(
		"",
		param.QueueName,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         param.PublishingBody,
		})
	//error
	if err != nil {
		return err
	}

	return

}

//consumeResponseTimeout 每次 publish 出去後啟動 goroutine 去判斷收 response msg 是否 timeout，
//如果 timeout 就自動送一個 timeout message 回 等待的 chan
func (m *MQ) consumeResponseMessage() {

	defer func() {
		if r := recover(); r != nil {
			errMsg := fmt.Sprintf("MQ consumeResponseMessage() panic %s", r)
			m.logger.LogFile(LevelPanic, errMsg)
		}
	}()

	var tsec = int64(m.timeoutSecond)
	for {
		select {

		//每n秒檢查一次receivers裡面所有的receiver的時間
		//超過timeout時間，就送入timeout response
		case <-time.After(time.Duration(tsec) * time.Second):

			now := GetNowWithDefaultLocation()
			var tsec = int64(m.timeoutSecond)

			var f = func(key, value interface{}) bool {

				//sn,_:=key.(int)
				rcv, _ := value.(*receiver)

				if rcv != nil {

					sec, _ := SubSecond(rcv.beginTime, now)
					//超過timeout秒數
					if sec >= tsec {

						//送入timeout message
						if rcv.ReceiveChan != nil {
							rcv.ReceiveChan <- m.getTimeoutCGResMessage(rcv.Serial) //timeout msg
						}
						//delete receiver from receives
						m.receivers.Delete(key)
					}
				}

				return true
			}

			m.rangeReceiver(f)

			//取德resMsg，並依據serial取得對應的chan，
			//回傳msg給chan
		case resMsg, ok := <-m.consumeResponseChan:

			ackErr := resMsg.Ack(false) //delete message
			if ok {
				if ackErr != nil {
					errMsg := fmt.Sprintf("MQ consumeResponseMessage() resMsg.ACK() err %s", ackErr.Error())
					m.logger.LogFile(LevelError, errMsg)
				}
				//chan closed
				if !ok {
					errMsg := fmt.Sprintf("MQ consumeResponseMessage() chan not ok")
					m.logger.LogFile(LevelError, errMsg)
				}

				s := &CGResMessage{}
				err := json.Unmarshal(resMsg.Body, s)
				if err != nil {
					errMsg := fmt.Sprintf("MQ consumeResponseMessage() unmarshal body error %s got body %+v ", err.Error(), resMsg.Body)
					m.logger.LogFile(LevelError, errMsg)
				} else {
					//取receiver，用serial當key
					rcv := m.loadReceiver(s.Serial)
					//取出後，無論是否取到，都delete
					m.deleteReceiver(s.Serial)
					//取到receiver，把msg放進去
					if rcv != nil && rcv.ReceiveChan != nil {
						rcv.ReceiveChan <- resMsg.Body
					}
				}
			}
		}
	}
}

func (m *MQ) getTimeoutCGResMessage(sn int64) []byte {
	//CGResMessage is struct for consume response message
	c := &CGResMessage{
		sn,
		"",
		"timeout",
		CGResMessageErrorCodeTimeout,
		json.RawMessage([]byte("[{}]")),
	}

	b, _ := json.Marshal(c)
	return b
}

//*******************************
//receivers is a sync.Map used to store receiver's chan with serial as key.
//loadReceiver gets receiver's chan and delete it from map
//如果已經被取走，有可能取不到
func (m *MQ) loadReceiver(sn int64) *receiver {

	v, ok := m.receivers.Load(sn)
	if ok {
		//value 轉型成 *receiver
		if r, rok := v.(*receiver); rok {
			return r
		}
		return nil
	}
	return nil
}
func (m *MQ) deleteReceiver(sn int64) {
	m.receivers.Delete(sn)
}

//func (m *MQ) storeReceiver(param *PublishParam) {
func (m *MQ) storeReceiver(key int64, rcv *receiver) (ok bool) {

	//loadOrStore
	//如果有就回傳value，且回傳true
	//如果沒有就儲存，且回傳false
	_, ok = m.receivers.LoadOrStore(key, rcv)

	// if ok =true =load
	// if ok =false =store

	//判斷是否put成功，回傳!ok
	return !ok
}
func (m *MQ) rangeReceiver(f func(key, value interface{}) bool) {
	m.receivers.Range(f)
}

//*******************************
func (m *MQ) GetMQHandlerRouter() *MQHandlerRouter {
	return m.mqHandlerRouter

}

//command handler (CGMessage handler) 負責 register handler with pattern
// 註冊的 handler 需實作 MQHandler interface
// MQ 收到 command mq 時，呼叫 mqHandler 檢查此 command 是否有註冊 MQHandler, 有則呼叫
type MQHandler interface {
	ServeMQ(msg *CGMessage)
}
type MQHandlerRouter struct {
	handlers sync.Map
}

//handle 註冊cmd及MQHandler
func (c *MQHandlerRouter) Handle(cmd string, h MQHandler) {
	c.handlers.Store(cmd, h)
}

func (c *MQHandlerRouter) getHandler(cmd string) (MQHandler, bool) {
	v, ok := c.handlers.Load(cmd)
	if ok {
		if m, ok := v.(MQHandler); ok {
			return m, ok
		}
		return nil, false
	}
	return nil, false
}
func (c *MQHandlerRouter) runHandler(msg *CGMessage, h MQHandler, cb func()) {
	h.ServeMQ(msg)
	cb()
}
