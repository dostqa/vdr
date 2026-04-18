package models

// type Request struct {
// 	ID        int
// 	Status    bool
// 	CreatedAt time.Timer
// }

type File struct {
	ID        int
	RequestID int
	FileName  string
}

type PDn struct {
	ID        int
	FileID    int
	TypeOf    string
	StartTime float64
	EndTime   float64
}

type Transcription struct {
	ID           int
	RequestID    int
	OriginalText string
	AnonText     string
}
