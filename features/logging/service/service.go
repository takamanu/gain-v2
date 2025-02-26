package service

import (
	"gain-v2/features/logging"
)

type LoggingService struct {
	d logging.LoggingDataInterface
}

func NewService(data logging.LoggingDataInterface) logging.LoggingServiceInterface {
	return &LoggingService{
		d: data,
	}
}

func (ls *LoggingService) AddLog(newData logging.LogData) (*logging.LogData, error) {

	res, err := ls.d.AddLog(newData)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ls *LoggingService) ViewLog() ([]logging.LogData, error) {
	res, err := ls.d.ViewLog()

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ls *LoggingService) ViewOneLog(logID string) (*logging.LogData, error) {
	res, err := ls.d.ViewOneLog(logID)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ls *LoggingService) DeleteLog(logID string) (bool, error) {
	_, err := ls.d.DeleteLog(logID)

	if err != nil {
		return false, err
	}

	return true, nil
}
