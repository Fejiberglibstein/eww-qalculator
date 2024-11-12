package message

const Port = "/tmp/eww-calc"

type Request uint8

type Message struct {
	Header uint8
	Data   string
}

type ListenerPort string

const (
	ExprResult ListenerPort = "expr_result"
)

const (
	Listen Request = iota
	SendExpr
)
