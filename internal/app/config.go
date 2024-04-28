package app

import "github.com/go-playground/validator/v10"

var Validator = validator.New()

type GRPC struct {
	Addr string `validate:"required,hostname_port"`
}

type config struct {
	GRPC   GRPC
	LogLvl int `validate:"min=-1,max=5"`
}
