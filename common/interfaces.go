package common

import (
	"time"
)

type Dispatcher interface {
	AddSubject(subject Subject) error
	GetSubjects() []Subject

	GetStartEvent() Event
	GetIdleEvent() Event
	GetCurrentDateTime() *time.Time

	Stop() error
	Run() error
}

type Event interface {
	Subscribe(handler EventHandler) error
	Unsubscribe(handler EventHandler) error
	Emit(args ...interface{}) []error
}

type Subject interface {
	Start() error
	Stop() error
	Join() error
	Eof() bool
	Dispatch() (bool, error)
	PeekDateTime() *time.Time

	GetDispatchPriority() int
	SetDispatchPriority(priority int)

	OnDispatcherRegistered(dispatcher Dispatcher) error
}

type Broker interface {
	Subject
	GetOrderUpdatedEvent() Event
	NotifyOrderEvent(orderEvent *OrderEvent)
	CancelOrder(order Order) error
}

type Order interface {
	GetId() uint64
	IsActive() bool
	IsFilled() bool
	GetExecutionInfo() OrderExecutionInfo
	AddExecutionInfo(info OrderExecutionInfo) error

	GetRemaining() int

	SwitchState(newState OrderState) error
}

type Bar interface {
	SetUseAdjustedValue(useAdjusted bool) error
	GetUseAdjValue() bool

	GetDateTime() *time.Time
	Open(adjusted bool) float64
	High(adjusted bool) float64
	Low(adjusted bool) float64
	Close(adjusted bool) float64
	Volume() int
	AdjClose() float64
	Frequency() Frequency
	Price() float64
}

type Bars interface {
	GetDateTime() *time.Time
	GetInstruments() []string
	GetBar(instrument string) Bar
	GetFrequencies() []Frequency

	AddBar(instrument string, bar Bar) error
}

type BarFeed interface {
	Subject
	GetNewValueEvent() Event
	GetCurrentBars() []Bar
}