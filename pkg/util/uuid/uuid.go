package uuid

import (
	"encoding/hex"
	"github.com/google/uuid"
	"github.com/rs/xid"
)

//GetUUID return UUID v4 encoded hexadecimal to 32 chars.
func GetUUID() string {
	u := uuid.New()
	return hex.EncodeToString(u[:])
}

//GetXID return 20-chars-length unique id.
func GetXID() string {
	return xid.New().String()
}
