package app

import (
	"context"
	"fmt"
	"net"
	"notis_filer/internal/filer"
	"notis_filer/internal/server"
	"path/filepath"
	"strings"

	common "github.com/RenterRus/text_coler"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"os"

	pkg "notis_filer/pkg/api"

	"github.com/sourcegraph/conc/pool"
	"google.golang.org/grpc"
)

type App struct {
	config
	ctx  context.Context
	cncl context.CancelFunc
	exit chan os.Signal
}

var (
	configPath = pflag.StringP("config", "c", "cfg/config.yaml", "path to config file")
	showHelp   = pflag.BoolP("help", "h", false, "Show help message")
)

func NewApp(exit chan os.Signal) *App {
	ctx, cnclfnc := context.WithCancel(context.Background())
	return &App{
		ctx:  ctx,
		cncl: cnclfnc,
		exit: exit,
	}
}

func (a *App) Run() error {
	var errRun error
	common.ColorPrintln(common.ForegBlack, common.BackLightGrey, " Prepare initialized ")
	defer common.ColorPrintln(common.ForegOrange, common.BackStd, "Module Run has been stopped")

	//
	// Config
	//
	common.ColorPrint(common.ForegDarkGreen, common.BackStd, "STAGE 1. Config ")

	pflag.Parse()
	if *showHelp {
		pflag.Usage()
		os.Exit(0)
	}

	dir, file := filepath.Split(*configPath)
	viper.AddConfigPath(dir)
	ext := filepath.Ext(file)
	name := strings.TrimSuffix(file, ext)

	viper.SetConfigName(name)
	viper.SetConfigType(strings.TrimPrefix(ext, "."))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		common.ColorPrintln(common.ForegBlack, common.BackRed, fmt.Errorf("read config: %w", err).Error())
		a.Stop(err)
	}

	if err := viper.Unmarshal(&a.config); err != nil {
		common.ColorPrintln(common.ForegBlack, common.BackRed, fmt.Errorf("unmarshal config: %w", err).Error())
		a.Stop(err)
	}

	if err := Validator.Struct(a.config); err != nil {
		common.ColorPrintln(common.ForegBlack, common.BackRed, fmt.Errorf("validation config: %w", err).Error())
		a.Stop(err)
	}

	common.ColorPrintln(common.ForegBlack, common.BackGreen, " successfully ")

	//
	// Log
	//
	common.ColorPrint(common.ForegDarkGreen, common.BackStd, "STAGE 2. Zero logger ")

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	switch a.config.LogLvl {
	case -1:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case 0:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case 2:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case 3:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case 4:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case 5:
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	common.ColorPrintln(common.ForegBlack, common.BackGreen, " successfully ")

	//
	// gRPC
	//
	common.ColorPrint(common.ForegDarkGreen, common.BackStd, "STAGE 3. gRPC ")

	listener, err := net.Listen("tcp", a.GRPC.Addr)
	if err != nil {
		common.ColorPrintln(common.ForegBlack, common.BackRed, " failed %v", err)
		errRun = a.Stop(err)
	}
	s := grpc.NewServer()

	pkg.RegisterNotesServer(s, server.NewService(
		filer.NewFiler(a.BaseDir),
	))

	common.ColorPrintln(common.ForegBlack, common.BackGreen, " successfully ")

	p := pool.New().WithErrors()

	p.Go(func() error {
		common.ColorPrintln(common.ForegGreen, common.BackStd, "grpc listen at [%v]", listener.Addr())

		if err := s.Serve(listener); err != nil {
			common.ColorPrintln(common.ForegBlack, common.BackRed, " gRPC ERROR: %v", err)
		}
		return nil
	})

	//
	// Stopped
	//

	<-a.ctx.Done()
	s.GracefulStop()
	listener.Close()
	p.Wait()
	return errRun
}

func (a *App) Stop(err error) error {
	common.ColorPrintln(common.ForegDarkYellow, common.BackStd, "Graceful shutdown initialized")
	a.cncl()
	close(a.exit)
	return err
}
