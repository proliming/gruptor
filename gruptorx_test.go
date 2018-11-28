// Description:
// Author: liming.one@bytedance.com
package gruptor

import (
	"errors"
	"fmt"
)

type AnEventFactory struct {
}

func (AnEventFactory) NewEvent() Event {
	return nil
}

type AnEventHandler struct {
}

func (eh *AnEventHandler) OnEvent(event Event, sequence int64) error {
	if event != sequence {
		warning := fmt.Sprintf("\nRace condition--Sequence: %d, Event: %d\n", sequence, event)
		fmt.Printf(warning)
		return errors.New(warning)
	}
	return nil
}
