package inerr

import (
	"fmt"

	"github.com/AsaHero/just-ask-bot/pkg/logger"
	"github.com/sirupsen/logrus"
)

func Err(err error, msg, scope, caller, callee string) error {
	if err == nil {
		err = fmt.Errorf("unknown error")
	}

	logger.Error(err.Error(), logrus.Fields{
		"summary": msg,
		"scope":   scope,
		"caller":  caller,
		"callee":  callee,
	})

	return err
}

func New(err error, msg, scope, caller, callee string) error {
	if err == nil {
		err = fmt.Errorf(msg)
	}

	logger.Error(err.Error(), logrus.Fields{
		"summary": msg,
		"scope":   scope,
		"caller":  caller,
		"callee":  callee,
	})

	return fmt.Errorf(msg)
}

func ErrAlertError(err error, msg, scope, caller, callee string) error {
	if err == nil {
		err = fmt.Errorf("unknown error")
	}

	logger.AlertError(err.Error(), logrus.Fields{
		"summary": msg,
		"scope":   scope,
		"caller":  caller,
		"callee":  callee,
	})

	return err
}

func ErrAlertWarn(err error, msg, scope, caller, callee string) error {
	if err == nil {
		err = fmt.Errorf("unknown error")
	}

	logger.AlertWarn(err.Error(), logrus.Fields{
		"summary": msg,
		"scope":   scope,
		"caller":  caller,
		"callee":  callee,
	})

	return err
}
