package srtgo

/*
#cgo LDFLAGS: -lsrt
#include <srt/srt.h>

int64_t srt_connection_time_wrapped(SRTSOCKET u, int *srterror, int *syserror)
{
	int64_t ret = srt_connection_time(u);
	if (ret < 0) {
		*srterror = srt_getlasterror(syserror);
	}
	return ret;
}


int64_t srt_time_now_wrapped(int *srterror, int *syserror)
{
	int64_t ret = srt_time_now();
	if (ret < 0) {
		*srterror = srt_getlasterror(syserror);
	}
	return ret;
}
*/
import "C"
import (
	"syscall"
	"time"
)

func (s SrtSocket) ConnectionDuration() (time.Duration, error) {
	srterr := C.int(0)
	syserr := C.int(0)
	t := C.srt_connection_time_wrapped(s.socket, &srterr, &syserr)
	if t < 0 {
		srterror := SRTErrno(srterr)
		if syserr < 0 {
			srterror.wrapSysErr(syscall.Errno(syserr))
		}
		return 0, srterror
	}
	return time.Duration(t), nil
}

func Now() (time.Duration, error) {
	srterr := C.int(0)
	syserr := C.int(0)
	t := C.srt_time_now_wrapped(&srterr, &syserr)
	if t < 0 {
		srterror := SRTErrno(srterr)
		if syserr < 0 {
			srterror.wrapSysErr(syscall.Errno(syserr))
		}
		return 0, srterror
	}
	return time.Duration(t), nil
}
