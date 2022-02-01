package xsync

type ErrorRecover struct {
	Err error
	Recover interface{}
}
func CoreGo(routine func() (err error)) (errRecoverCh chan ErrorRecover) {
	// 使用 1 缓存通道防止routine 泄露
	errRecoverCh = make(chan ErrorRecover, 1)
	go func() {
		errRecover := ErrorRecover{}
		defer func() {
			r := recover()
			if r != nil {
				errRecover.Recover = r
			}
			errRecoverCh <- errRecover
		}()
		routineErr := routine()
		if routineErr != nil {
			errRecover.Err = routineErr
		}
	}()
	return
}