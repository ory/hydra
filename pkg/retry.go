package pkg

import "time"

func Retry(maxWait time.Duration, failAfter time.Duration, f func() error) (err error) {
	loopWait := time.Millisecond * 500
	retryStart := time.Now()
	for retryStart.Add(failAfter).Before(time.Now()) {
		if err = f(); err == nil {
			return nil
		}
		time.Sleep(loopWait)

	}
	return err
}
