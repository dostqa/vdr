package message

type InputMessage struct {
	RequestID int    `json:"request_id"`
	FilePath  string `json:"file_path"`
}
