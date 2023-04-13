package main

import (
	"TrackMaster/initializer"
	"TrackMaster/model"
	"TrackMaster/model/task"
	"TrackMaster/pkg"
	"TrackMaster/pkg/worker"
	"TrackMaster/third_party/slack"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

func init() {
	initializer.LoadEnvVariables()
	initializer.ConnectDB()
}

func NewConsumer() (sarama.Consumer, error) {
	brokerAddr, ok := os.LookupEnv("BROKER")
	var brokerList []string
	if !ok {
		brokerAddr = "localhost:9092"
		brokerList = []string{brokerAddr}
	} else {
		brokerList = strings.Split(brokerAddr, ",")
	}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_0_0_0
	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

func NewLock(projectID string) *redsync.Mutex {
	redisAddr, ok := os.LookupEnv("REDIS")
	if !ok {
		redisAddr = "localhost:6379"
	}
	client := goredislib.NewClient(&goredislib.Options{
		Addr: redisAddr,
	})

	pool := goredis.NewPool(client)
	rs := redsync.New(pool)

	mutexName := "pikachu_job_" + projectID
	mutex := rs.NewMutex(mutexName)
	return mutex
}

func DoJob(msg []byte, wp *worker.Pool) *pkg.Error {
	schedule := model.Schedule{}
	err := json.Unmarshal(msg, &schedule)
	if err != nil {
		wp.Errors <- pkg.NewError(pkg.ServerError, err.Error()).WithDetails("发生在unmarshal msg时")

	}

	// 判断schedule的最后执行时间，比interval大才会执行
	if time.Since(schedule.LastExecuted) > schedule.Interval {
		// 获取锁
		mutex := NewLock(schedule.ProjectID)
		log.Println("获取锁：", mutex.Name())
		if err := mutex.Lock(); err != nil {
			return pkg.NewError(pkg.ServerError, err.Error()).WithDetails("获取锁失败")
		}

		defer func() {
			if ok, err := mutex.Unlock(); !ok || err != nil {
				panic("unlock failed")
			}
			log.Println("释放锁：", mutex.Name())
		}()

		_ = slack.SendMessage("开始执行任务...")
		// 先把schedule的最后执行时间修改了
		schedule.LastExecuted = time.Now()
		initializer.DB.Save(&schedule)

		project := model.Project{
			ID: schedule.ProjectID,
		}

		job := task.UpdateStoryTask{
			Project:  &project,
			Interval: schedule.Interval,
			WP:       wp,
		}
		wp.Jobs <- &job
	} else {
		_ = slack.SendMessage("未到执行周期，放弃任务")
	}
	return nil
}

func main() {
	consumer, err := NewConsumer()
	if err != nil {
		log.Fatalln("Failed to create consumer:", err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln("Failed to close consumer:", err)
		}
	}()

	go func() {
		http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("pong"))
		})

		err := http.ListenAndServe(":8000", nil)
		if err != nil {
			log.Fatalln("failed to start http server")
		}
		log.Println("starting http server on port 8000")
	}()

	topic, ok := os.LookupEnv("TOPIC")
	if !ok {
		topic = "pikachu-track"
	}

	partitionConsumer, err := consumer.ConsumePartition(topic, 1, sarama.OffsetNewest)
	if err != nil {
		log.Fatalln("Failed to create partitionConsumer:", err)
	}

	errorCh := make(chan *pkg.Error, worker.MaxQueue)
	wp := worker.NewWorkerPool(errorCh)
	wp.Begin()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Println("Received: ", string(msg.Value))
			_ = slack.SendMessage("Received message: " + string(msg.Value))
			err := DoJob(msg.Value, wp)
			wp.Errors <- err

		case err := <-partitionConsumer.Errors():
			log.Println("Error: ", err.Error())
			_ = slack.SendMessage("Received Error: " + err.Error())

		case <-errorCh:
			log.Println("Error: ", err.Error())
			_ = slack.SendMessage("Error: " + err.Error())

		case <-signals:
			return
		}
	}
}
