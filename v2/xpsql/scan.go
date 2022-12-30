package xpsql

import (
	"database/sql"
	"database/sql/driver"

	"code.olapie.com/sugar/v2/xtype"
)

type supportedxtype interface {
	*xtype.Point | *xtype.Place | *xtype.PhoneNumber | *xtype.FullName | *xtype.Money | map[string]string
}

func Scan[T supportedxtype](v *T) sql.Scanner {
	switch val := any(v).(type) {
	case **xtype.Point:
		return &pointScanner{v: val}
	case **xtype.PhoneNumber:
		return &phoneNumberScanner{v: val}
	case **xtype.Place:
		return &placeScanner{v: val}
	case **xtype.Money:
		return &moneyScanner{v: val}
	case **xtype.FullName:
		return &fullNameScanner{v: val}
	case *map[string]string:
		return &hstoreScanner{m: val}
	default:
		panic("unsupported type")
	}
}

func Value[T supportedxtype](v T) driver.Valuer {
	switch val := any(v).(type) {
	case *xtype.Point:
		return &pointValuer{v: val}
	case *xtype.PhoneNumber:
		return &phoneNumberValuer{v: val}
	case *xtype.Place:
		return &placeValuer{v: val}
	case *xtype.Money:
		return &moneyValuer{v: val}
	case *xtype.FullName:
		return &fullNameValuer{v: val}
	case map[string]string:
		return mapToHstore(val)
	default:
		panic("unsupported type")
	}
}
