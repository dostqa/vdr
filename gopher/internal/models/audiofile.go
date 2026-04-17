package models

type AudioFile struct {
	ID       string
	Filename string `json:"filename"`
	Data     []byte `json:"-"`
}

func NewAudioFile(id string, filename string, data []byte) AudioFile {
	return AudioFile{
		ID:       id,
		Filename: filename,
		Data:     data,
	}
}
