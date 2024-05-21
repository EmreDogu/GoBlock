package p2p

import "github.com/EmreDogu/GoBlock/internal/blockchain/block"

type Message interface {
	GetType() string
	GetInterval() int64
	GetParent() *block.Block
	Run()
}

type Node interface {
	SendNextBlockMessage()
	ReceiveMessage(*Messages)
	ReceiveBlock(*block.Block)
	GetUseCBR() bool
	GetRegion() int
	GetID() int
}

type Task struct {
	Message
}
