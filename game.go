package main

import (
	"client/network"
	"client/protocol"
	"fmt"
	"time"
)

type LifeGameClient struct {
	PacketChan chan protocol.Packet
}

func ConnectLifeGameServer() {

	client := LifeGameClient{}
	/* 패킷헤더 사이즈 정의 */
	protocol.InitPacketHeaderSize()

	/* 패킷 채널 생성 */
	client.PacketChan = make(chan protocol.Packet, 256)

	snFunctor := network.SessionNetworkFunctor{
		OnConnect:           client.OnConnect,
		OnReceive:           client.OnReceive,
		PacketTotalSizeFunc: network.PacketTotalSize,
		PacketHeaderSize:    protocol.GetPacketHeaderSize(),
	}

	go client.PacketProcess()

	network.ConnectServer(snFunctor)
}

func (client *LifeGameClient) PacketProcess() {
	for {
		select {
		case packet := <-client.PacketChan:
			{
				bodySize := packet.DataSize
				bodyData := packet.Data
				packetId := packet.Id

				if packetId == protocol.PACKET_ID_LOGIN_RES {
					fmt.Println("Login Response")
					ProcessPacketLogin(bodySize, bodyData)
				} else if packetId == protocol.PACKET_ID_JOIN_RES {
					fmt.Println("Join Response")
					ProcessPacketJoin(bodySize, bodyData)
				}
			}
		}
	}
}

func ProcessPacketJoin(bodySize int16, bodyData []byte) {
	var joinRes protocol.LoginResPacket
	result := (&joinRes).Decoding(bodyData)
	if !result {
		fmt.Println("Join Failed")
		return
	}

	if joinRes.ErrorCode != protocol.ERROR_CODE_NONE {
		fmt.Println("Join Failed")
		return
	}

	fmt.Println("Join!")
	SendLogin("tuuna2983", "password")
	fmt.Println("Send Login")
}

func ProcessPacketLogin(bodySize int16, bodyData []byte) {
	var loginRes protocol.LoginResPacket
	result := (&loginRes).Decoding(bodyData)
	if !result {
		fmt.Println("Login Failed")
		return
	}

	if loginRes.ErrorCode != protocol.ERROR_CODE_NONE {
		fmt.Println("Login Failed")
		return
	}

	fmt.Println("Login!")
}

func SendPing() {
	for {
		return
		time.Sleep(1 * time.Millisecond)
		pingReq := protocol.PingReqPacket{
			Ping: protocol.PING,
		}

		packet, packetSize := pingReq.EncodingPacket()
		network.SendToServer(packet, packetSize)
	}
}

func (client *LifeGameClient) OnConnect() {
	fmt.Println("Connect!")
	/* 로그인 패킷 전송 */
	SendJoin("tuuna2983", "password", "tuuna")
	//SendLogin("tuuna2983", "password")
	//fmt.Println("Send Login")
}

func (client *LifeGameClient) OnReceive(packetData []byte) {
	packetID := protocol.PeekPacketID(packetData)
	bodySize, packetBody := protocol.PeekPacketBody(packetData)

	packet := protocol.Packet{
		Id:       packetID,
		DataSize: bodySize,
		Data:     make([]byte, bodySize),
	}

	copy(packet.Data, packetBody)
	client.PacketChan <- packet
}

func SendLogin(userID, userPW string) {
	loginReq := protocol.LoginReqPacket{
		UserID: make([]byte, protocol.MAX_USER_ID_BYTE_LENGTH),
		UserPW: make([]byte, protocol.MAX_USER_PW_BYTE_LENGTH),
	}
	copy(loginReq.UserID[:], []byte(userID))
	copy(loginReq.UserPW[:], []byte(userPW))
	packet, packetSize := loginReq.EncodingPacket()

	network.SendToServer(packet, packetSize)
}

func SendJoin(userID, userPW, userNAME string) {
	joinReq := protocol.JoinReqPacket{
		UserID:   make([]byte, protocol.MAX_USER_ID_BYTE_LENGTH),
		UserPW:   make([]byte, protocol.MAX_USER_PW_BYTE_LENGTH),
		UserName: make([]byte, protocol.MAX_USER_NAME_BYTE_LENGTH),
	}

	copy(joinReq.UserID[:], []byte(userID))
	copy(joinReq.UserPW[:], []byte(userPW))
	copy(joinReq.UserName[:], []byte(userNAME))

	packet, packetSize := joinReq.EncodingPacket()
	network.SendToServer(packet, packetSize)
}
