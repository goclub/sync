package xsync

type ErrorRecover struct {
	Err error
	Recover interface{}
}
func Go(routine func() (err error)) (errRecoverCh chan ErrorRecover) {
	errRecover := ErrorRecover{}
	errRecoverCh = make(chan ErrorRecover)
	go func() {
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
