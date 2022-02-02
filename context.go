package xsync

import "context"

func Ctx(ctx context.Context, routine func() error) (err error) {
	errRecoverCh := Go(routine)
	select {
	case errRecover := <- errRecoverCh:
		if errRecover.Err != nil {
			return errRecover.Err
		}
		if errRecover.Recover != nil {
			panic(errRecover.Recover)
		}
	case <- ctx.Done():
		return ctx.Err()
	}
	return
}
