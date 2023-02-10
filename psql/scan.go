package psql

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"code.olapie.com/sugar/contacts"
	"code.olapie.com/sugar/v2/types"
)

type supportedScanTypes interface {
	*types.Point | *types.Place | *contacts.PhoneNumber | *types.FullName | *types.Money | map[string]string
}

func Scan[T supportedScanTypes](v *T) sql.Scanner {
	switch val := any(v).(type) {
	case **types.Point:
		return &pointScanner{v: val}
	case **contacts.PhoneNumber:
		return &phoneNumberScanner{v: val}
	case **types.Place:
		return &placeScanner{v: val}
	case **types.Money:
		return &moneyScanner{v: val}
	case **types.FullName:
		return &fullNameScanner{v: val}
	case *map[string]string:
		return &hstoreScanner{m: val}
	default:
		panic(fmt.Sprintf("unsupported scan type: %T", v))
	}
}

func Value[T supportedScanTypes](v T) driver.Valuer {
	switch val := any(v).(type) {
	case *types.Point:
		return &pointValuer{v: val}
	case *contacts.PhoneNumber:
		return &phoneNumberValuer{v: val}
	case *types.Place:
		return &placeValuer{v: val}
	case *types.Money:
		return &moneyValuer{v: val}
	case *types.FullName:
		return &fullNameValuer{v: val}
	case map[string]string:
		return mapToHstore(val)
	default:
		panic(fmt.Sprintf("unsupported scan type: %T", v))
	}
}
