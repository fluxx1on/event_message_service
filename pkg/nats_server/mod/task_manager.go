package mod

import (
	"log/slog"
	"slices"
	"time"

	"github.com/nats-io/nats.go"
	server "gitlab.com/fluxx1on_group/event_message_service/pkg/nats_server"
)

type Task struct {
	Msg   *nats.Msg
	start time.Time
	end   time.Time
}

// In - method to define how time.Now() compare with basic interval
// Statuses:
// 1 - task Expired; need to clean
// 0 - task is ready to consume
// -1 - it's to early
func (t Task) In() int {
	now := time.Now()
	if t.start.Compare(now) == -1 && t.end.Compare(now) == 1 {
		return 0
	} else if t.end.Compare(now) == -1 {
		return 1
	}
	return -1
}

func (t Task) SetIn(t_start time.Time, t_end time.Time) {
	t.start = t_start
	t.end = t_end
}

// TaskManager consume messages with required time interval.
// Manager wait for Task.start timing and after run MsgHandler function.
type TaskManager struct {
	conn *server.Connection

	// gChan consume nats messages and send it to MsgTimeHadnler function.
	// cChan consume expired messages and Terminate these with MsgTermHadnler function.
	gChan chan *nats.Msg
	cChan chan *nats.Msg

	router map[string]HandlerGroup
	tasks  []Task

	// rate is time.Sleep() duration for call Task.In() in cycle one more time
	rate time.Duration
	stop chan struct{}
}

func NewTaskManager(
	conn *server.Connection,
	router map[string]HandlerGroup,
	stop chan struct{},
) *TaskManager {

	return &TaskManager{
		conn:   conn,
		gChan:  make(chan *nats.Msg, len(router)),
		cChan:  make(chan *nats.Msg, len(router)),
		router: router,
		tasks:  make([]Task, 0),
		rate:   20 * time.Second,
		stop:   stop,
	}
}

func (m *TaskManager) Subscribe() {
	for subj := range m.router {
		ch := make(chan *nats.Msg)
		newSub, err := m.conn.ChanSubscribe(subj, ch)
		if err != nil {
			slog.Error("Subscription failed",
				slog.String("subj", subj),
				slog.String("ErrorMsg", err.Error()),
			)
		}
		defer newSub.Unsubscribe()
		go func() {
			for {
				m.gChan <- <-ch
			}
		}()
	}

	go m.consume()

	for {
		var toClean = make([]int, 0)
		for iter, task := range m.tasks {
			in := task.In()
			if in == 0 {
				m.gChan <- task.Msg
			} else if in == 1 {
				toClean = append(toClean, iter)
			}
		}
		m.clean(toClean...)
		time.Sleep(m.rate)
	}
}

func (m *TaskManager) consume() {
	for {
		select {
		case <-m.stop:
			return
		case msg := <-m.gChan:
			go func() {
				ok, time := m.router[msg.Subject].Get(msg)
				if ok != false {
					m.cChan <- msg
				} else if time.In() == -1 {
					m.tasks = append(m.tasks, time)
				}
			}()
		case msg := <-m.cChan:
			go func() {
				m.router[msg.Subject].Clean(msg)
			}()
		}
	}
}

func (m *TaskManager) clean(toClean ...int) {
	for _, c := range toClean {
		m.tasks = slices.Delete[[]Task, Task](m.tasks, c, c+1)
	}
}
