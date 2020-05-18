package syslog

import (
	"log/syslog"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/pkg/errors"
)

var (
	ErrUnknownFacility = errors.New("unknown syslog facility")
)

var (
	facilityNames = map[string]syslog.Priority{
		"kern":     syslog.LOG_KERN,
		"user":     syslog.LOG_USER,
		"mail":     syslog.LOG_MAIL,
		"daemon":   syslog.LOG_DAEMON,
		"auth":     syslog.LOG_AUTH,
		"syslog":   syslog.LOG_SYSLOG,
		"lpr":      syslog.LOG_LPR,
		"news":     syslog.LOG_NEWS,
		"uucp":     syslog.LOG_UUCP,
		"cron":     syslog.LOG_CRON,
		"authpriv": syslog.LOG_AUTHPRIV,
		"ftp":      syslog.LOG_FTP,
		"local0":   syslog.LOG_LOCAL0,
		"local1":   syslog.LOG_LOCAL1,
		"local2":   syslog.LOG_LOCAL2,
		"local3":   syslog.LOG_LOCAL3,
		"local4":   syslog.LOG_LOCAL4,
		"local5":   syslog.LOG_LOCAL5,
		"local6":   syslog.LOG_LOCAL6,
		"local7":   syslog.LOG_LOCAL7,
	}
)

func parseFacility(s string) (syslog.Priority, error) {
	if s == "" {
		return syslog.LOG_LOCAL0, nil
	}
	if s[0] < '0' || s[0] > '9' {
		s = strings.ToLower(s)
		for name, v := range facilityNames {
			if name == s {
				return v, nil
			}
		}
		return syslog.LOG_LOCAL0, ErrUnknownFacility
	}
	id, err := strconv.Atoi(s)
	if err != nil {
		return syslog.LOG_LOCAL0, ErrUnknownFacility
	}
	for _, v := range facilityNames {
		if int(v) == id {
			return v, nil
		}
	}
	return syslog.LOG_LOCAL0, ErrUnknownFacility
}

func bytes2Str(b []byte) string {
	bytesHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	strHeader := reflect.StringHeader{
		Data: bytesHeader.Data,
		Len:  bytesHeader.Len,
	}
	return *(*string)(unsafe.Pointer(&strHeader))
}

// bytesToStringInplace converts string to byte slice.
func str2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}
