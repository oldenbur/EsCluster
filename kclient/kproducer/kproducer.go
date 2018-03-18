package kproducer

import (
	"context"
	"github.com/cihub/seelog"
	"golang.org/x/time/rate"
	"time"
)

func Produce() {

	seelog.Info("Rater starting")
	wait()
	seelog.Info("Rater complete")
}

func wait() {

	last := time.Now()
	lim := rate.NewLimiter(rate.Every(500*time.Millisecond), 5)
	t := time.NewTimer(3 * time.Second)

	for {
		select {
		case <-t.C:
			seelog.Info("Timer fired")
			return

		default:
			seelog.Info("Waiting")
			ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
			if err := lim.Wait(ctx); err != nil {
				seelog.Infof("Wait returned error: %v", err)
			}
			seelog.Infof("Wait returned - delay: %v", time.Since(last))
			last = time.Now()
		}
	}

}
