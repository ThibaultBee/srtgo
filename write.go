package srtgo

/*
#cgo LDFLAGS: -lsrt
#include <srt/srt.h>

int srt_sendmsg2_wrapped(SRTSOCKET u, const char* buf, int len, SRT_MSGCTRL *mctrl, int *srterror, int *syserror)
{
	int ret = srt_sendmsg2(u, buf, len, mctrl);
	if (ret < 0) {
		*srterror = srt_getlasterror(syserror);
	}
	return ret;
}


*/
import "C"
import (
	"errors"
	"syscall"
	"time"
	"unsafe"
)

func srtSendMsg2Impl(u C.SRTSOCKET, buf []byte, msgctrl *C.SRT_MSGCTRL) (n int, err error) {
	srterr := C.int(0)
	syserr := C.int(0)
	n = int(C.srt_sendmsg2_wrapped(u, (*C.char)(unsafe.Pointer(&buf[0])), C.int(len(buf)), msgctrl, &srterr, &syserr))
	if n < 0 {
		srterror := SRTErrno(srterr)
		if syserr < 0 {
			srterror.wrapSysErr(syscall.Errno(syserr))
		}
		err = srterror
		n = 0
	}
	return
}

// Write data to the SRT socket
func (s SrtSocket) Write(b []byte) (n int, err error) {

	//Fastpath:
	if !s.blocking {
		s.pd.reset(ModeWrite)
	}
	n, err = srtSendMsg2Impl(s.socket, b, nil)

	for {
		if !errors.Is(err, error(EAsyncSND)) || s.blocking {
			return
		}
		s.pd.wait(ModeWrite)
		n, err = srtSendMsg2Impl(s.socket, b, nil)
	}
}

// WriteWithSrcTime data to the SRT socket with msgctrl
func (s SrtSocket) WriteWithSrcTime(b []byte, srcTime time.Duration) (n int, err error) {
	msgctrl := C.SRT_MSGCTRL{
		flags:        0, // no flags set
		msgttl:       C.SRT_MSGTTL_INF,
		inorder:      0, // false - not in order (matters for msg mode only)
		boundary:     0,
		srctime:      C.longlong(srcTime / 1000), // srctime: from ns to us
		pktseq:       C.SRT_SEQNO_NONE,
		msgno:        C.SRT_MSGNO_NONE,
		grpdata:      nil, // grpdata not supplied
		grpdata_size: 0,
	}

	//Fastpath:
	if !s.blocking {
		s.pd.reset(ModeWrite)
	}
	n, err = srtSendMsg2Impl(s.socket, b, &msgctrl)

	for {
		if !errors.Is(err, error(EAsyncSND)) || s.blocking {
			return
		}
		s.pd.wait(ModeWrite)
		n, err = srtSendMsg2Impl(s.socket, b, &msgctrl)
	}
}
