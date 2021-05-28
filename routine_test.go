package xsync_test

//
// func TestRoutine_Go(t *testing.T) {
// 	// no err no panic
// 	{
// 		exec := false
// 		xsync.Go(func() error {
// 			exec = true
// 			return nil
// 		})
// 		err, recoverValue := routine.Wait()
// 		assert.Equal(t, exec, true)
// 		assert.NoError(t, err)
// 		assert.Nil(t, recoverValue)
// 	}
// 	// has err no panic
// 	{
// 		exec := false
// 		routine := new(xsync.Routine)
// 		routine.Go(func() error {
// 			exec = true
// 			return errors.New("abc")
// 		})
// 		err, recoverValue := routine.Wait()
// 		assert.Equal(t, exec, true)
// 		assert.EqualError(t, err, "abc")
// 		assert.Nil(t, recoverValue)
// 	}
// 	// no err has panic
// 	{
// 		exec := false
// 		routine := new(xsync.Routine)
// 		routine.Go(func() error {
// 			exec = true
// 			panic("xyz")
// 			return nil
// 		})
// 		err, recoverValue := routine.Wait()
// 		assert.Equal(t, exec, true)
// 		assert.NoError(t, err)
// 		assert.Equal(t, recoverValue, "xyz")
// 	}
// }
//
