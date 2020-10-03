package entity

type (
	// Message 消息体
	Message struct {
		Data []byte
		Name string
		End  bool
		MD5  string
		Size int64
	}
)
