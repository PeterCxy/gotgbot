package utils

import (
	"time"

	"github.com/PeterCxy/gotgbot/support/types"
)

// Input "grabber"
var grabbers map[int64]types.Grabber
var grabberChannel chan types.Grabber
var grabberDelChannel chan int64
var okChannel chan bool
var states map[int64]map[string]interface{}

func GrabberDaemon() {
	grabbers = make(map[int64]types.Grabber)
	grabberChannel = make(chan types.Grabber)
	grabberDelChannel = make(chan int64)
	okChannel = make(chan bool)
	states = make(map[int64]map[string]interface{})

	for {
		select {
			case grabber := <-grabberChannel:
				grabbers[grabber.Uid + grabber.Chat] = grabber
				states[grabber.Uid + grabber.Chat] = make(map[string]interface{})
				okChannel <- true
			case id := <-grabberDelChannel:
				delete(grabbers, id)
				delete(states, id)
		}
	}
}

// Whether the current session is grabbed by a processor
func HasGrabber(uid int64, chat int64) bool {
	grabber := grabbers[uid + chat]
	return grabber.Name != "" && grabber.Processor != nil
}

// Get the name and the processor of the current session
func Grabber(uid int64, chat int64) (string, types.CommandProcessor) {
	grabber := grabbers[uid + chat]
	return grabber.Name, grabber.Processor
}

// Grab all the non-command inputs in the current session
// Returns the state storage object
func SetGrabber(cmd types.Grabber) *map[string]interface{} {
	grabberChannel <- cmd
	<-okChannel // Wait
	return GrabberState(cmd.Uid, cmd.Chat)
}

// Release the current session
func ReleaseGrabber(uid int64, chat int64) {
	if HasGrabber(uid, chat) {
		grabberDelChannel <- uid + chat
	}
}

// The state storage
func GrabberState(uid int64, chat int64) *map[string]interface{} {
	state := states[uid + chat]
	return &state
}

// The scheduler
var schedule map[time.Time]func()

func init() {
	schedule = make(map[time.Time]func())
}

func PostDelayed(callback func(), duration time.Duration) {
	schedule[time.Now().Add(duration)] = callback
}

// Start scheduled functions
// Called from main loop
func ScheduleDaemon() {
	for {
		now := time.Now()
		for k, v := range schedule {
			if now.Equal(k) || now.After(k) {
				go v()
				delete(schedule, k)
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}
