//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package utils

import (
	"bytes"
	"time"

	"github.com/go-sql-driver/mysql"
)

type NullTime struct {
	mysql.NullTime
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return nt.Time.MarshalJSON()
	} else {
		return []byte("null"), nil
	}
}

func (nt *NullTime) UnmarshalJSON(data []byte) error {
	if bytes.Compare(data, []byte("null")) == 0 {
		nt.Valid = false
		return nil
	}

	t := time.Now()
	err := t.UnmarshalJSON(data)

	if err != nil {
		return err
	}

	nt.Valid = true
	nt.Time = t

	return nil
}
