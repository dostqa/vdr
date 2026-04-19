package message

type ObjectPdn struct {
	Text      string  `json:"text"`
	Type      string  `json:"type"`
	StartTime float64 `json:"start_time"`
	EndTime   float64 `json:"end_time"`
}

type OutputMessage struct {
	RequestID    int         `json:"request_id"`
	OldFilePath  string      `json:"old_file_path"`
	NewFilePath  string      `json:"new_file_path"`
	OriginalText string      `json:"original_text"`
	AnonText     string      `json:"anon_text"`
	ObjectsPdns  []ObjectPdn `json:"objects_pdns"`
}
