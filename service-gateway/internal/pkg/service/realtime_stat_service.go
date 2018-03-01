package service

import (
	"sync"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

type RealtimeStatService struct {
	lock                        sync.RWMutex
	currentUserCountSubscribers map[string]chan<- int64
	totalAccessCountSubscribers map[string]chan<- int64
}

func NewRealtimeStatService() *RealtimeStatService {
	return &RealtimeStatService{
		lock:                        sync.RWMutex{},
		currentUserCountSubscribers: map[string]chan<- int64{},
		totalAccessCountSubscribers: map[string]chan<- int64{},
	}
}

func (m *RealtimeStatService) IncreaseCurrentUserCountSubscriber(uuid string, subscriber chan<- int64) error {
	if m == nil {
		return nil
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.currentUserCountSubscribers[uuid]; ok {
		return status.Errorf(codes.AlreadyExists, "The uuid %q is already in use by someone", uuid)
	}

	m.currentUserCountSubscribers[uuid] = subscriber
	count := int64(len(m.currentUserCountSubscribers))
	sugar.Infow("Increased current connection count", "count", count)

	go func() {
		m.lock.Lock()
		defer m.lock.Unlock()

		for k, _ := range m.currentUserCountSubscribers {
			if s, ok := m.currentUserCountSubscribers[k]; ok {
				s <- count
			}
		}
	}()

	return nil
}

func (m *RealtimeStatService) DecreaseCurrentUserCountSubscriber(uuid string) {
	if m == nil {
		return
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	m.lock.Lock()
	defer m.lock.Unlock()

	if c, ok := m.currentUserCountSubscribers[uuid]; ok {
		close(c)
		delete(m.currentUserCountSubscribers, uuid)

		count := int64(len(m.currentUserCountSubscribers))
		sugar.Infow("Decreased current connection count", "count", count)

		go func() {
			m.lock.Lock()
			defer m.lock.Unlock()

			for k, _ := range m.currentUserCountSubscribers {
				if s, ok := m.currentUserCountSubscribers[k]; ok {
					s <- count
				}
			}
		}()

	}

}
