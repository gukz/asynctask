package asynctask

import (
	"errors"
)

const NoResultFound string = "No errors found"

type producer struct {
	*AsyncBase
}

func (t *producer) CheckHealth() bool {
	backendHealthy := false
	brokerHealthy := false
	if t.broker == nil {
		logger.Info("The producer's broker is nil")
		brokerHealthy = true
	} else {
		brokerHealthy = t.broker.CheckHealth()
	}
	if t.backend == nil {
		logger.Info("The producer's backend is nil")
		backendHealthy = true
	} else {
		backendHealthy = t.backend.CheckHealth()
	}
	return backendHealthy && brokerHealthy
}

func (t *producer) Send(message *Message) error {
	if msg, err := encode(message); err != nil {
		return err
	} else {
		return t.broker.PushMessage(t.queue, msg)
	}
}

func (t *producer) GetResult(taskId string) (*Result, error) {
	resultBytes, err := t.backend.GetResult(taskId)
	if err != nil {
		return nil, err
	}
	if resultBytes == nil {
		return nil, errors.New(NoResultFound)
	}
	result := &Result{}
	err = decode(resultBytes, result)
	return result, err
}
