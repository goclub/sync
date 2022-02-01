package xsync

import "context"

func Go(ctx context.Context, routine func() (err error)) (errRecoverCh chan ErrorRecover) {
	errRecoverCh = make(chan ErrorRecover, 1)
	routineResultCh := CoreGo(routine)
	go func() {
		select {
		case routineResult := <- routineResultCh:
			errRecoverCh <- routineResult
		case <- ctx.Done():
			errRecoverCh <- ErrorRecover{
				Err:  ctx.Err(),
			}
		}
	}()
	return
}