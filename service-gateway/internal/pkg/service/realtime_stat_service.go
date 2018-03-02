package service

import (
	"sync"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"github.com/jinzhu/gorm"
	"github.com/1ambda/go-ref/service-gateway/internal/pkg/model"
)

type RealtimeStatService struct {
	lock                        sync.RWMutex
	db                          *gorm.DB
	currentUserCountSubscribers map[string]chan<- int64
	totalAccessCountSubscribers map[string]chan<- int64
}

func NewRealtimeStatService(db *gorm.DB) *RealtimeStatService {
	return &RealtimeStatService{
		lock:                        sync.RWMutex{},
		db:                          db,
		currentUserCountSubscribers: map[string]chan<- int64{},
		totalAccessCountSubscribers: map[string]chan<- int64{},
	}
}

func (r *RealtimeStatService) IncreaseCurrentUserCountSubscriber(uuid string, subscriber chan<- int64) error {
	if r == nil {
		return nil
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	r.lock.Lock()
	defer r.lock.Unlock()

	if _, ok := r.currentUserCountSubscribers[uuid]; ok {
		return status.Errorf(codes.AlreadyExists, "The uuid %q is already in use by someone", uuid)
	}

	r.currentUserCountSubscribers[uuid] = subscriber
	count := int64(len(r.currentUserCountSubscribers))
	sugar.Infow("Increased current connection count", "count", count)

	go func() {
		r.lock.Lock()
		defer r.lock.Unlock()

		for k, _ := range r.currentUserCountSubscribers {
			if s, ok := r.currentUserCountSubscribers[k]; ok {
				s <- count
			}
		}
	}()

	return nil
}

func (r *RealtimeStatService) DecreaseCurrentUserCountSubscriber(uuid string) {
	if r == nil {
		return
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	r.lock.Lock()
	defer r.lock.Unlock()

	if c, ok := r.currentUserCountSubscribers[uuid]; ok {
		close(c)
		delete(r.currentUserCountSubscribers, uuid)

		count := int64(len(r.currentUserCountSubscribers))
		sugar.Infow("Decreased current connection count", "count", count)

		go func() {
			r.lock.Lock()
			defer r.lock.Unlock()

			for k, _ := range r.currentUserCountSubscribers {
				if s, ok := r.currentUserCountSubscribers[k]; ok {
					s <- count
				}
			}
		}()
	}
}

func (r *RealtimeStatService) AddTotalAccessCountSubscriber(uuid string, subscriber chan<- int64) error {
	if r == nil {
		return nil
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	sugar.Infow("Adding TotalAccessCount Subscriber", "uuid", uuid)

	r.lock.Lock()
	defer r.lock.Unlock()

	if _, ok := r.totalAccessCountSubscribers[uuid]; ok {
		return status.Errorf(codes.AlreadyExists, "The uuid %q is already in use by someone", uuid)
	}

	r.totalAccessCountSubscribers[uuid] = subscriber

	return nil
}

func (r *RealtimeStatService) DeleteTotalAccessCountSubscriber(uuid string) {
	if r == nil {
		return
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	sugar.Infow("Deleting TotalAccessCount Subscriber", "uuid", uuid)

	r.lock.Lock()
	defer r.lock.Unlock()

	if c, ok := r.totalAccessCountSubscribers[uuid]; ok {
		close(c)
		delete(r.totalAccessCountSubscribers, uuid)
	}
}

func (r *RealtimeStatService) BroadcastToTalAccessCount() error {
	if r == nil {
		return nil
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	var totalAccessCount int64
	err := r.db.Table(model.AccessTable).Count(&totalAccessCount).Error
	if err != nil {
		sugar.Errorf("Failed to get total access count", "error", err)
		return status.Errorf(codes.Internal, "Failed to get total access count")
	}

	go func() {
		sugar.Infow("Broadcasting Total Access Count", "count", totalAccessCount)

		r.lock.Lock()
		defer r.lock.Unlock()

		for k, _ := range r.totalAccessCountSubscribers {
			if s, ok := r.totalAccessCountSubscribers[k]; ok {
				s <- totalAccessCount
			}
		}
	}()


	return nil
}
