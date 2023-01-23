package xpsql

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"code.olapie.com/sugar/v2/xtype"
	"code.olapie.com/sugar/xpsql/v2/internal/composite"
)

var (
	_ driver.Valuer = (*moneyValuer)(nil)
	_ sql.Scanner   = (*moneyScanner)(nil)
)

type moneyScanner struct {
	v **xtype.Money
}

func (ms *moneyScanner) Scan(src any) error {
	if src == nil {
		return nil
	}

	s, ok := src.(string)
	if !ok {
		b, ok := src.([]byte)
		if !ok {
			return fmt.Errorf("src is not []byte or string")
		}
		s = string(b)
	}

	if len(s) == 0 {
		return nil
	}

	fields, err := composite.ParseFields(s)
	if err != nil {
		return fmt.Errorf("parse composite fields %s: %w", s, err)
	}

	if len(fields) != 2 {
		return fmt.Errorf("parse composite fields %s: got %v", s, fields)
	}
	m := new(xtype.Money)
	m.Currency = fields[0]
	m.Amount = fields[1]
	if err != nil {
		return fmt.Errorf("parse amount %s: %w", fields[1], err)
	}
	*ms.v = m
	return nil
}

type moneyValuer struct {
	v *xtype.Money
}

func (mv *moneyValuer) Value() (driver.Value, error) {
	if mv == nil || mv.v == nil {
		return nil, nil
	}
	return fmt.Sprintf("(%s,%s)", mv.v.Currency, mv.v.Amount), nil
}
