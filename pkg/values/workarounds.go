package values

import sdkTypes "github.com/kairos-io/kairos-sdk/types"

var WorkaroundsMap = map[Distro]map[Architecture]map[string][]Workaround{
	Ubuntu: {
		ArchAMD64: {
			"20.04": {},
			"24.04": {TestWorkAround},
		},
		ArchARM64: {
			"20.04": {},
		},
	},
}

func TestWorkAround(s *System, l sdkTypes.KairosLogger) error {
	l.Logger.Info().Msg("Running TestWorkAround")
	return nil
}
