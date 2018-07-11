package grpcbase
import(
	"context"
	grpc2 "google.golang.org/grpc"
	"github.com/processout/grpc-go-pool"
	"time"
	"whale-market/library/logger"
	"go.uber.org/zap"
	"sync"
)

var(
	serviceMap = make(map[string] *pollService, 0)
	lock = &sync.Mutex{}
)

type ServiceConfig struct {
	serviceName string
	address string
	init int
	capacity int
}

type pollService struct{
	pool *grpcpool.Pool
}

func (service *ServiceConfig)createPoll() (*grpcpool.Pool, error) {
	p, errPoll := grpcpool.New(func() (*grpc2.ClientConn, error) {
		return grpc2.Dial(service.address, grpc2.WithInsecure())
	}, service.init, service.capacity, time.Second)
	if errPoll != nil{
		logger.ZapError.Error("init poll fail ", zap.Error(errPoll))
	}
	return p, errPoll
}

func (service *ServiceConfig)GetServiceConn() (*grpcpool.ClientConn, error) {
	if poolGet, ok := serviceMap[service.serviceName]; ok {
		if checkPoolValid(poolGet.pool) {
			conn, err :=  poolGet.pool.Get(context.Background())
			return conn, err
		}else{
			service.closePool()
		}
	}
	lock.Lock()
	defer lock.Unlock()
	if _, ok := serviceMap[service.serviceName]; !ok {
		poolCreated, poolErr := service.createPoll()
		if poolErr != nil {
			return nil, poolErr
		}
		serviceMap[service.serviceName] = &pollService{poolCreated}
	}
	conn, err :=  serviceMap[service.serviceName].pool.Get(context.Background())
	return conn, err
}

func (service *ServiceConfig)closePool() {
	if _, ok := serviceMap[service.serviceName]; ok {
		serviceMap[service.serviceName].pool.Close()
	}
	delete(serviceMap, service.serviceName)
}

func CreateService(serviceName, address string, init, capacity int) *ServiceConfig{
	return &ServiceConfig{serviceName, address, init, capacity}
}

func checkPoolValid(pool *grpcpool.Pool) bool {
	return !pool.IsClosed()
}

