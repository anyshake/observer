package miniseed

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/anyshake/observer/drivers/explorer"
)

func (m *MiniSeedService) saveSequence() error {
	dataBytes, err := json.Marshal(m.miniseedSequence)
	if err != nil {
		return err
	}

	err = os.WriteFile(m.getSequenceFilePath(), dataBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (m *MiniSeedService) readSequence() (sequence map[string]int, err error) {
	data, err := os.ReadFile(m.getSequenceFilePath())
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &sequence)
	if err != nil {
		return nil, err
	}

	for _, channel := range []string{
		explorer.EXPLORER_CHANNEL_CODE_Z,
		explorer.EXPLORER_CHANNEL_CODE_E,
		explorer.EXPLORER_CHANNEL_CODE_N,
	} {
		_, ok := sequence[channel]
		if !ok {
			return nil, fmt.Errorf("sequence is missing for channel %s", channel)
		}
	}

	return sequence, nil
}
