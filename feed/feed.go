package feed

import (
	"goalgotrade/common"
	"goalgotrade/core"
	lg "goalgotrade/logger"
	"time"
)

type regDs struct {
	key  string
	freq common.Frequency
}

type BaseFeed struct {
	*core.DefaultSubject
	event        common.Event
	dataseries   map[string]map[common.Frequency]common.DataSeries
	registeredDs []regDs
	maxlen       int
}

func NewBaseFeed(maxlen int) *BaseFeed {
	return &BaseFeed{
		DefaultSubject: core.NewDefaultSubject(),
		event:          core.NewEvent(),
		dataseries:     make(map[string]map[common.Frequency]common.DataSeries),
		maxlen:         maxlen,
	}
}

func (b *BaseFeed) Reset() {
	b.dataseries = make(map[string]map[common.Frequency]common.DataSeries)
	for _, v := range b.registeredDs {
		b.RegisterDataSeries(v.key, v.freq)
	}
}

func (b *BaseFeed) CreateDataSeries(key string, maxlen int) common.DataSeries {
	lg.Logger.Error("not implemented")
	panic("not implemented")
}

func (b *BaseFeed) GetNextValues() (*time.Time, common.Bars, common.Frequency, error) {
	lg.Logger.Error("not implemented")
	panic("not implemented")
}

func (b *BaseFeed) GetNextValuesAndUpdateDS() (*time.Time, common.Bars, common.Frequency, error) {
	datetime, values, freq, err := b.GetNextValues()
	if err != nil || datetime == nil {
		keys := values.GetInstruments()
		for _, k := range keys {
			if v, ok := b.dataseries[k]; !ok {
				b.dataseries[k] = make(map[common.Frequency]common.DataSeries)
			} else {
				if v2, ok2 := v[freq]; ok2 {
					v2.AppendWithDateTime(*datetime, values.GetBar(k))
				} else {
					b.dataseries[k][freq] = b.CreateDataSeries(k, b.maxlen)
				}
			}
		}
	}
	return datetime, values, freq, err
}

func (b *BaseFeed) RegisterDataSeries(key string, freq common.Frequency) error {
	if _, ok := b.dataseries[key]; !ok {
		b.dataseries[key] = map[common.Frequency]common.DataSeries{}
	}
	if _, ok := b.dataseries[key][freq]; !ok {
		b.dataseries[key][freq] = b.CreateDataSeries(key, b.maxlen)
		for _, v := range b.registeredDs {
			if v.key == key && v.freq == freq {
				return nil
			}
		}
		b.registeredDs = append(b.registeredDs, regDs{key: key, freq: freq})
	}
	return nil
}

func (b *BaseFeed) GetNewValuesEvent() common.Event {
	return b.event
}

func (b *BaseFeed) Dispatch() (bool, error) {
	// TODO: check if freq here is needed
	datetime, values, _, err := b.GetNextValuesAndUpdateDS()
	if err != nil {
		return false, err
	}
	if datetime != nil {
		b.event.Emit(datetime, values)
	}
	return datetime != nil && err == nil, nil
}

func (b *BaseFeed) Eof() bool {
	panic("not implemented")
}
