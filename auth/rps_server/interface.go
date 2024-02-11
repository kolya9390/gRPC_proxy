package rpcserver

type RPCServer interface {
	StartServer(port string, rcvr ...any) error
}


type User struct {
	ID       int64
	Email    string
	PassHash []byte
}