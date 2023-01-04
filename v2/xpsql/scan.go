package xpsql

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"code.olapie.com/sugar/v2/xtype"
)

type supportedScanTypes interface {
	*xtype.Point | *xtype.Place | *xtype.PhoneNumber | *xtype.FullName | *xtype.Money | map[string]string
}

func Scan[T supportedScanTypes](v *T) sql.Scanner {
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
		panic(fmt.Sprintf("unsupported scan type: %T", v))
	}
}

func Value[T supportedScanTypes](v T) driver.Valuer {
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
		panic(fmt.Sprintf("unsupported scan type: %T", v))
	}
}
