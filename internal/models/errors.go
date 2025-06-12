package models

import "errors"

var (
	ErrConfigPathNotProvided   = errors.New("config path didn't provide")
	ErrBadConfigPort           = errors.New("port must be upper than 0")
	ErrBadConfigFactor         = errors.New("factor must be upper than 1")
	ErrBadConfigURL            = errors.New("destination url is required")
	ErrBadConfigMinInterval    = errors.New("min interval must be upper than 0")
	ErrBadConfigMaxInterval    = errors.New("min interval must be a positive num")
	ErrBadConfigMinMaxInterval = errors.New("min interval must be lower or equal than max interval")
)
