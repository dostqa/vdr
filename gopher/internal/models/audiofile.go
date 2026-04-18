package models

type AudioFile struct {
	ID       string `json:"-"`
	Filename string `json:"filename"`
}

func NewAudioFile(id string, filename string) AudioFile {
	return AudioFile{
		ID:       id,
		Filename: filename,
	}
}
