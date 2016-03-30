package structure

type Refresh struct {
	Type int   `json:"type"`
	Id   int   `json:"uid"`
	St   int64 `json:"st"`
	Et   int64 `json:"et"`
}
