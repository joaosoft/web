package errors

type ListErr []*Err

type Err struct {
	Previous *Err   `json:"previous,omitempty"`
	Code     string `json:"code"`
	Err      string `json:"error"`
}
