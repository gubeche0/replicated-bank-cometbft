package model

import "errors"

type QueryTx struct {
	Method string `json:"method"`
	Value  string `json:"value,omitempty"`
}

func (q *QueryTx) Validate() error {
	if q.Method != "listAll" && q.Method != "getByName" {
		return errors.New("return must be listAll or getByName")
	}

	if q.Method == "getByName" && q.Value == "" {
		return errors.New("value must be non-empty for getByName")
	}

	return nil
}
