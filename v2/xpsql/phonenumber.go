package xpsql

import (
	"code.olapie.com/sugar/v2/conv"
	"code.olapie.com/sugar/v2/xcontact"
	"code.olapie.com/sugar/v2/xpsql/internal/composite"
	"database/sql"
	"database/sql/driver"
	"fmt"
)

var (
	_ driver.Valuer = (*phoneNumberValuer)(nil)
	_ sql.Scanner   = (*phoneNumberScanner)(nil)
)

type phoneNumberScanner struct {
	v **xcontact.PhoneNumber
}

func (ps *phoneNumberScanner) Scan(src any) error {
	if src == nil {
		return nil
	}

	s, err := conv.ToString(src)
	if err != nil {
		return fmt.Errorf("parse string: %w", err)
	}
	if len(s) == 0 {
		return nil
	}

	fields, err := composite.ParseFields(s)
	if err != nil {
		return fmt.Errorf("parse composite fields %s: %w", s, err)
	}

	if len(fields) != 3 {
		return fmt.Errorf("parse composite fields %s: got %v", s, fields)
	}

	n := new(xcontact.PhoneNumber)
	n.Code, err = conv.ToInt32(fields[0])
	if err != nil {
		return fmt.Errorf("parse code %s: %w", fields[0], err)
	}
	n.Number, err = conv.ToInt64(fields[1])
	if err != nil {
		return fmt.Errorf("parse code %s: %w", fields[1], err)
	}
	//n.Extension = fields[2]
	*ps.v = n
	return nil
}

type phoneNumberValuer struct {
	v *xcontact.PhoneNumber
}

func (pv *phoneNumberValuer) Value() (driver.Value, error) {
	if pv == nil {
		return nil, nil
	}
	//ext := strings.Replace(pv.v.Extension, ",", "\\,", -1)
	ext := ""
	s := fmt.Sprintf("(%d,%d,%s)", pv.v.Code, pv.v.Number, ext)
	return s, nil
}
