package service

import (
	"errors"
	//"sync"
	"time"

	"fmt"
	gorillaWebsocket "github.com/gorilla/websocket"
)

type WSConn struct {
	active bool
	conn   *gorillaWebsocket.Conn
	//sendChan          chan []byte //chan contains data to send
	receiveError      bool // receive got error, conn inactivate
	sendError         bool // send got error, conn inactivate
	//mutex             *sync.Mutex
	connID            string                  //conn uuid for get specific connData
	connectLostChan   chan string             //contains lost connID  send to websocket
	receivePacketChan chan *ReceivePacketData //send to websocket
	aliveTimeout      time.Duration
}

//newConnData returns a new connData.
func NewWSConn(co *gorillaWebsocket.Conn,
	connLostChan chan string,
	receivePacketChan chan *ReceivePacketData,
	aliveTimeout time.Duration) (*WSConn, error) {

	conUUID:= uuidGetV4()

	c := &WSConn{
		active: false,
		conn:   co,
		//sendChan:          make(chan []byte, 300),
		//mutex:             &sync.Mutex{},
		sendError:         false,
		receiveError:      false,
		connID:            conUUID,
		connectLostChan:   connLostChan,
		receivePacketChan: receivePacketChan,
		aliveTimeout:      aliveTimeout,
	}

	return c, nil
}

func (cd *WSConn) Start() {
	cd.active = true
	cd.receive()
}

func (cd *WSConn) Stop() {
	cd.active = false
	//set read timeout 0
	if cd.conn != nil {
		cd.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 1))
	}
}


func (cd *WSConn) GetConnID() string {
	return cd.connID
}

//func (cd *WSConn) closeSendChan() {
//	close(cd.sendChan)
//}

//after accept and after add user
//Close conn 1. aliveTimeout 2.readPacket error
func (cd *WSConn) receive() {

	go func() {

		defer func(connID string) {

			//panic
			if r := recover(); r != nil {
				//todo

				//cd.closeSendChan() //只在這裡close
				cd.receiveError = true
				cd.active = false

				//沒處理過才處理
				if nil != cd.conn && !cd.sendError {
					cd.connectLostChan <- connID // 只在這裡加入
					cd.conn.Close()
					cd.conn = nil
				}
			}

		}(cd.connID)

		for cd != nil && cd.conn!=nil && cd.active && cd.connID != "" {

			//Set aliveTimeout, add duration from now
			if err := cd.conn.SetReadDeadline(time.Now().Add(cd.aliveTimeout)); err != nil {
				//if err := cd.conn.SetReadDeadline(time.Now().Add(time.Second * 10)); err != nil {

				cd.active = false
				//close timeout
				cd.conn.SetReadDeadline(time.Time{})
				return
			}

			//read from gorilla websocket conn
			msgType, resBuff, err2 := cd.conn.ReadMessage()

			//After readPacket,maybe timeout already, so return
			if !cd.active {
				cd.conn.SetReadDeadline(time.Time{})
				return
			}

			//disable readDeadline time.Time{} =zero
			//After readPacket, diable aliveTimeout
			if deadLineErr := cd.conn.SetReadDeadline(time.Time{}); deadLineErr != nil {
				cd.active = false
				return
			}

			//fmt.Printf("receive resBuff=%v \n",resBuff)

			//Check if connection lost or error
			if err2 != nil {
				fmt.Printf("receive error, resBuff=%v \n", resBuff)

				//check close error
				if _, ok := err2.(*gorillaWebsocket.CloseError); ok {
					cd.active = false
					return
				}
				// common.PrintColor(common.ColorGreen, cd.receive, "read error userID :", strconv.Itoa(cd.userID), err.Error())
				cd.active = false
				return
			}

			//text message
			if msgType == gorillaWebsocket.TextMessage {
				cd.receivePacketChan <- &ReceivePacketData{Packet: resBuff, ConnID: cd.connID, MessageType: TextMessage}
			} else {
				cd.receivePacketChan <- &ReceivePacketData{Packet: resBuff, ConnID: cd.connID, MessageType: BinaryMessage}
			}
		}
	}()
}

//Send write packet body into sendChan
//and go routine do real send
//func (cd *WSConn) Send(sentPacket []byte) {
//	if cd.active && cd.sendChan != nil {
//		cd.sendChan <- sentPacket
//	}
//}

func (cd *WSConn) send(packet []byte, messageType int) error {

	defer func(connID string) {

		if r := recover(); r != nil {
			//todo

			cd.active = false
			cd.sendError = true

			//沒處理過，才處理
			if nil != cd.conn && !cd.receiveError {
				cd.connectLostChan <- connID // 只在這裡加入
				cd.conn.Close()
				cd.conn = nil

			}
		}
	}(cd.connID)

	if nil == cd {
		return errors.New("WSConn is nil")
	}
	if nil == cd.conn {
		return errors.New("WSConn.conn is nil")
	}
	if len(packet) < 1 {
		return errors.New(fmt.Sprintf("packet length error got %d", len(packet)))
	}

	//check
	if nil != cd && cd.active && nil != cd.conn {

		//err := cd.conn.WriteMessage(gorillaWebsocket.BinaryMessage, b)
		//err := cd.conn.WriteMessage(gorillaWebsocket.TextMessage, packet)
		err := cd.conn.WriteMessage(messageType, packet)
		if err != nil {
			cd.active = false
		}
	}
	return nil
}
