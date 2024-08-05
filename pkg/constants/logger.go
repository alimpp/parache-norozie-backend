package constants

import (
	"github.com/rs/zerolog/log"
)

var Logger = log.With().Str("service", ServiceName).Logger()
