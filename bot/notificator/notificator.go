package notificator

import (
	"time"

	"github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"
	"github.com/H1ghN0on/go-tgbot-engine/logger"
)

type TickerData struct {
	isTicking bool
	ticker    *time.Ticker
	done      chan bool
}

type Notificationer interface {
	GetMessages() []bottypes.Message
	GetUsers() []bottypes.User
	GetTimeoutSec() int
}

type Notification struct {
	users      func() []bottypes.User
	timeoutSec int
}

type StaticNotification struct {
	Notification
	messages []bottypes.Message
}

type DynamicNotification struct {
	Notification
	messages func() []bottypes.Message
}

func (nf StaticNotification) GetMessages() []bottypes.Message {
	return nf.messages
}

func (nf DynamicNotification) GetMessages() []bottypes.Message {
	return nf.messages()
}

func (nf Notification) GetUsers() []bottypes.User {
	return nf.users()
}

func (nf Notification) GetTimeoutSec() int {
	return nf.timeoutSec
}

type Notificator struct {
	notifications   map[*TickerData]Notificationer
	timeoutCallback func(notification Notificationer)
}

func (nf Notificator) AddNotification(notification Notificationer) {
	ticker := &TickerData{}
	nf.notifications[ticker] = notification
}

func (nf *Notificator) startTimer(ticker *TickerData, notification Notificationer) {
	if ticker.isTicking {
		return
	}

	for len(ticker.done) > 0 {
		<-ticker.done
	}

	ticker.done = make(chan bool)
	ticker.ticker = time.NewTicker(time.Duration(notification.GetTimeoutSec() * int(time.Second)))
	ticker.isTicking = true

	go func() {
		for {
			select {
			case <-ticker.done:
				logger.Notificator().Info("done written")
			case <-ticker.ticker.C:
				nf.timeoutCallback(notification)
			}
		}
	}()
}

func (nf *Notificator) stopTimer(ticker *TickerData) {
	if !ticker.isTicking {
		return
	}

	ticker.done <- true
	ticker.ticker.Stop()
	ticker.isTicking = false
}

func (nf *Notificator) Start() {
	logger.Notificator().Info("starting all notifications")

	for key, notification := range nf.notifications {
		nf.startTimer(key, notification)
	}
}

func (nf *Notificator) Stop() {
	logger.Notificator().Info("stopping all notifications")

	for key := range nf.notifications {
		nf.stopTimer(key)
	}
}

func NewStaticNotification(messages []bottypes.Message, users func() []bottypes.User, timeoutSec int) *StaticNotification {
	notification := &Notification{
		users:      users,
		timeoutSec: timeoutSec,
	}
	return &StaticNotification{
		messages:     messages,
		Notification: *notification,
	}
}

func NewDynamicNotification(messages func() []bottypes.Message, users func() []bottypes.User, timeoutSec int) *DynamicNotification {
	notification := &Notification{
		users:      users,
		timeoutSec: timeoutSec,
	}
	return &DynamicNotification{
		messages:     messages,
		Notification: *notification,
	}
}

type AnyNotificationInterface interface {
	StaticNotification | DynamicNotification
}

func NewNotificator(notifications []Notificationer, cb func(Notificationer)) *Notificator {

	notificationMap := make(map[*TickerData]Notificationer)

	for _, notification := range notifications {
		ticker := &TickerData{}
		notificationMap[ticker] = notification
	}

	return &Notificator{
		timeoutCallback: cb,
		notifications:   notificationMap,
	}
}
