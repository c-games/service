package service

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	gorillaWebsocket "github.com/gorilla/websocket"
)

/*
   Packet Structure=head+type+body
   head 4 bytes
   type 1 bytes
   body n bytes
   minimum packet length = 6 bytes
*/

const (
	//Packet structure 4bytes head(packet length) +1bytes type +body

	//PacketTypeJSON is const of packet data type json
	PacketTypeJSON int = 0
	//PacketTypeClass is const of packet data type class
	PacketTypeClass int = 1
	//PacketTypeHeartbeat is const of packet data type heartbeat
	PacketTypeHeartbeat int = 2

	packetLengthMin  int = 4 + 1 + 1
	packetLengthMax  int = 40960
	packetHeadLength int = 4
	packetTypeLength int = 1

	//TextMessage the same as gorilla message type
	TextMessage int = 1
	//BinaryMessage the same as gorilla message type
	BinaryMessage int = 2
)

var (
	//pingPongPacket for connData pingPong
	pingPongPacket = composePacket(uint16(PacketTypeHeartbeat), []byte{1})

	//gorillaWebsocketUpgrader gorillaWebsocket.Upgrader
)

type WSParameter struct {
	Address         string
	Port            string
	ConnPoolSize    int
	ChanPoolSize    int
	AcceptTimeout   time.Duration
	AliveTimeout    time.Duration
	ReadBufferSize  int
	WriteBufferSize int
}

//ReceivePacketData contain parsed date which is read from connection
//Without packet head
//With data type and packet body
type ReceivePacketData struct {
	Packet      []byte
	ConnID      string
	MessageType int
}

//SendPacketData store receiveIDs and packet to send
type SendPacketData struct {
	receiverIDs []string
	packetBody  []byte
}

type Websocket struct {
	active  bool
	address string
	port    string

	//外部
	//connect 後 id 放這裡
	connectIDChan chan string

	//receivePacketChan 收到packet放這裡，
	receivePacketChan chan *ReceivePacketData

	//connID in this chan is for others to get
	lostConnIDChan chan string

	connectDataTable map[string]*WSConn
	mutex            *sync.Mutex

	acceptTimeout time.Duration
	aliveTimeout  time.Duration
	httpServer    *http.Server

	upgrader gorillaWebsocket.Upgrader
}

//NewService make a new SocketManager
func NewWebsocket(param *WSParameter) (retWS *Websocket, retErr error) {

	defer func() {
		if r := recover(); r != nil {
			retErr = fmt.Errorf("panic %s", r)
			retWS = nil
		}
	}()

	s := &Websocket{
		address: param.Address,
		port:    param.Port,
		mutex:   &sync.Mutex{},
		active:  false,
	}

	s.connectIDChan = make(chan string, param.ConnPoolSize)
	s.connectDataTable = make(map[string]*WSConn)
	s.receivePacketChan = make(chan *ReceivePacketData, param.ChanPoolSize)
	s.lostConnIDChan = make(chan string, param.ChanPoolSize)
	s.acceptTimeout = param.AcceptTimeout
	s.aliveTimeout = param.AliveTimeout

	s.upgrader = gorillaWebsocket.Upgrader{
		// HandshakeTimeout specifies the duration for the handshake to complete.
		//HandshakeTimeout time.Duration

		// ReadBufferSize and WriteBufferSize specify I/O buffer sizes. If a buffer
		// size is zero, then a default value of 4096 is used. The I/O buffer sizes
		// do not limit the size of the messages that can be sent or received.
		ReadBufferSize:  param.ReadBufferSize,
		WriteBufferSize: param.WriteBufferSize,

		// Subprotocols specifies the server's supported protocols in order of
		// preference. If this field is set, then the Upgrade method negotiates a
		// subprotocol by selecting the first match in this list with a protocol
		// requested by the client.
		//Subprotocols []string

		// Error specifies the function for generating HTTP error responses. If Error
		// is nil, then http.Error is used to generate the HTTP response.
		//Error func(w http.ResponseWriter, r *http.Request, status int, reason error)

		// CheckOrigin returns true if the request Origin header is acceptable. If
		// CheckOrigin is nil, the host in the Origin header must not be set or
		// must match the host of the request.
		//CheckOrigin func(r *http.Request) bool
		//不檢查origin
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		// EnableCompression specify if the server should attempt to negotiate per
		// message compression (RFC 7692). Setting this value to true does not
		// guarantee that compression will be supported. Currently only "no context
		// takeover" modes are supported.
		EnableCompression: false,
	}

	return s, nil
}

/*
//Initial to init package variable
//acceptTimeout is socket accept timout
//aliveTimeout is used to check if connection is alive. If timeout, connection is closed.
func (s *Websocket) initial(connectPoolCapacity int,
	packetChanCapacity int,
	connectionLostChanCapacity int,
	acceptTimeout time.Duration,
	aliveTimeout time.Duration) (reErr error) {

	s.connectIDChan = make(chan string)
	s.connectDataTable = make(map[string]*WSConn)
	s.receivePacketChan = make(chan *ReceivePacketData, packetChanCapacity)
	s.lostConnIDChan = make(chan string, connectionLostChanCapacity)
	s.acceptTimeout = acceptTimeout
	s.aliveTimeout = aliveTimeout

	gorillaWebsocketUpgrader = gorillaWebsocket.Upgrader{
		// HandshakeTimeout specifies the duration for the handshake to complete.
		//HandshakeTimeout time.Duration

		// ReadBufferSize and WriteBufferSize specify I/O buffer sizes. If a buffer
		// size is zero, then a default value of 4096 is used. The I/O buffer sizes
		// do not limit the size of the messages that can be sent or received.
		ReadBufferSize:  10240,
		WriteBufferSize: 10240,

		// Subprotocols specifies the server's supported protocols in order of
		// preference. If this field is set, then the Upgrade method negotiates a
		// subprotocol by selecting the first match in this list with a protocol
		// requested by the client.
		//Subprotocols []string

		// Error specifies the function for generating HTTP error responses. If Error
		// is nil, then http.Error is used to generate the HTTP response.
		//Error func(w http.ResponseWriter, r *http.Request, status int, reason error)

		// CheckOrigin returns true if the request Origin header is acceptable. If
		// CheckOrigin is nil, the host in the Origin header must not be set or
		// must match the host of the request.
		//CheckOrigin func(r *http.Request) bool
		//不檢查origin
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		// EnableCompression specify if the server should attempt to negotiate per
		// message compression (RFC 7692). Setting this value to true does not
		// guarantee that compression will be supported. Currently only "no context
		// takeover" modes are supported.
		EnableCompression: false,
	}

	reErr = nil
	return reErr
}
*/

func (s *Websocket) handleAccept(w http.ResponseWriter, r *http.Request) {

	if s.active {

		conn, err := s.upgrader.Upgrade(w, r, nil)

		if err != nil {
			if conn != nil {
				conn.Close()
			}

			//if HandshakeError
			if _, ok := err.(gorillaWebsocket.HandshakeError); ok {
				//http.Error(w, "Not a websocket handshake", 400)
				//common.PrintInfo(s.handlerAccept, "handshakeError", err.Error())
				//此時未有connection
			}
			return
		}

		//create new connData for receive and send
		cd, _ := NewWSConn(conn, s.lostConnIDChan, s.receivePacketChan, s.aliveTimeout)
		s.addConnData(cd)
		cd.Start()
		s.connectIDChan<-cd.connID
	}
}

//Start start read packet action and others
func (s *Websocket) Start() error {

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", s.handleAccept)

		s.httpServer = &http.Server{
			Addr:           s.address + ":" + s.port,
			Handler:        mux,
			ReadTimeout:    20 * time.Second,
			WriteTimeout:   20 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		if err := s.httpServer.ListenAndServe(); err != nil {
			panic(fmt.Sprintf("httpServer listenAndServer error =%s", err.Error()))
		}

	}()
	s.active = true
	return nil
}

//Stop stop read packet action and others
func (s *Websocket) Stop() error {
	s.active = false
	s.stopConns()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

//Add add new conn data to pool
func (s *Websocket) addConnData(cd *WSConn) {

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.connectDataTable[cd.GetConnID()] = cd
}

//Remove remove conn data in pool
func (s *Websocket) removeConnData(id string) (isFound bool) {
	isFound = false

	if ok := s.containConnData(id); !ok {
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.connectDataTable, id)
	isFound = true
	return
}

//StartConn starts connection action
//func (s *Websocket) startConnData(id string) {
//	if cd, ok := s.getConnData(id); ok {
//		cd.Start()
//	}
//}


////StopConn stops action of connection
//func (s *Websocket) stopConnData(id string) {
//	if cd, ok := s.getConnData(id); ok {
//		cd.Stop()
//	}
//}

func (s *Websocket) stopConns() {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	for _, cd := range s.connectDataTable {
		cd.Stop()
	}
}

func (s *Websocket) getConnData(id string) (*WSConn, bool) {
	s.mutex.Lock()
	c, ok := s.connectDataTable[id]
	s.mutex.Unlock()
	if ok {
		return c, true
	}
	return nil, false
}

//Contains return if this id is already have
func (s *Websocket) containConnData(id string) bool {
	s.mutex.Lock()
	_, ok := s.connectDataTable[id]
	s.mutex.Unlock()
	return ok
}

//GetPacketChan return packetChan
func (s *Websocket) GetConnectConnIDChan() <-chan string {
	return s.connectIDChan
}

//GetPacketChan return packetChan
func (s *Websocket) GetReceivePacketChan() <-chan *ReceivePacketData {
	return s.receivePacketChan
}

//GetConnectionLostChan returns connectionLostChan containing connID
//Others using connid to do his work
func (s *Websocket) GetLostConnIDChan() <-chan string {
	return s.lostConnIDChan
}

func (s *Websocket) SendAll(packet []byte) error {
	defer func() {
		if r := recover(); r != nil {
			//todo
		}
	}()

	for _, v := range s.connectDataTable {

		v.send(packet, gorillaWebsocket.TextMessage)
	}

	return nil
}

//Send send packet to client.
func (s *Websocket) Send(receiverConnIDs []string, sentPacket []byte) error {

	defer func() {
		if r := recover(); r != nil {
			//todo
		}
	}()

	if !s.active || len(receiverConnIDs) < 1 {
		return errors.New("send error")
	}

	//idLen := len(receiverConnIDs)
	//for i := 0; i < idLen; i++ {
	//	cd, ok := s.getConnData(receiverConnIDs[i])
	//	if ok {
	//		cd.Send(sentPacket)
	//	}
	//}

	for _, v := range receiverConnIDs {
		cd, ok := s.getConnData(v)
		if ok {
			cd.send(sentPacket, gorillaWebsocket.TextMessage)
		}
	}

	return nil
}

//Disconnect disconnects target user's socket connection.
func (s *Websocket) Disconnect(id string) error {
	if cd, ok := s.getConnData(id); ok {
		cd.Stop()
		// s.disconnectChan <- id
		//不能先close，因為這時conn可能在waiting,且有sendChan
		//只有在最後close的地方才能close
		return nil
	}
	return errors.New("websocket.disconnect errors")
}

//listen wait and listen connect
//accept handshake begin
//func (s *Websocket) listen() error {
//
//	mux := http.NewServeMux()
//	mux.HandleFunc("/", s.handlerAccept)
//
//	s.httpServer = &http.Server{
//		Addr:           s.address + ":" + s.port,
//		Handler:        mux,
//		ReadTimeout:    20 * time.Second,
//		WriteTimeout:   20 * time.Second,
//		MaxHeaderBytes: 1 << 20,
//	}
//
//	s.httpServer.ListenAndServe()
//
//	return nil
//}

//ComposePacket add head and type to packet
//return a complete packet
func composePacket(packetType uint16, packetBody []byte) []byte {
	pl := uint32(len(packetBody))

	headBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(headBytes, pl)

	typeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(typeBytes, packetType)

	res := make([]byte, pl+4+1)
	copy(res[0:4], headBytes[:])
	copy(res[4:5], typeBytes[1:])
	copy(res[5:], packetBody[:])

	return res

}
