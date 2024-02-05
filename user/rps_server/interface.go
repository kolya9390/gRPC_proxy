package rpcserver

type RPCServer interface {
	StartServer(port string, rcvr ...any) error
}

type User struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
}