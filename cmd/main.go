package main

import (
	"errors"
	"fmt"
	"notis_filer/internal/app"
	"os"
	"os/signal"
	"syscall"

	common "github.com/RenterRus/text_coler"

	"github.com/sourcegraph/conc/pool"
)

func main() {
	common.ColorPrintln(common.ForegBlack, common.BackDarkGrey, " The start of the application launch ")

	ErrSignalReceiver := errors.New("a completion signal has been received")

	p := pool.New().WithErrors()
	exit := make(chan os.Signal, 1)
	a := app.NewApp(exit)

	p.Go(func() error {
		return a.Run()
	})

	p.Go(func() error {
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		_, ok := <-exit
		fmt.Println()
		if ok {
			a.Stop(ErrSignalReceiver)
		}

		return nil
	})

	err := p.Wait()
	common.ColorPrintln(common.ForegBlack, common.BackRed, " Notes is stopped ")
	if err != nil {
		common.ColorPrintln(common.ForegRed, common.BackStd, " witch error: %v\n", err.Error())
	}
	fmt.Println()
}
