//go:build windows

// Package engine implements the lid/sleep state machine as a single
// event-driven goroutine. External packages only post a single trigger event;
// the engine queries the current power/lid state itself and decides what to do
// (schedule a re-sleep, cancel, or nothing), owning all state and the debounced
// timer internally.
package engine

import (
	"sync"
	"sync/atomic"
	"time"

	"lidresleep/internal/i18n"
	"lidresleep/internal/log"
	"lidresleep/internal/win32"
)

// delayMS is the debounced delay before re-sleeping after a wake with the lid
// still closed. Atomic: written by the GUI thread, read by the engine loop.
var delayMS atomic.Int32

func init() { delayMS.Store(3000) }

// verifyDelay is how long after a sleep attempt the verify timer fires. The
// timer only runs while the system is awake, so when it fires we compare the
// wall-clock elapsed time against the delay to decide whether the sleep worked:
// a short elapsed time means the system stayed awake (retry), while an elapsed
// time noticeably longer than the delay means the timer was held up by a real
// suspend (sleep took effect, no retry).
const verifyDelay = 5 * time.Second

// verifyConfirmMargin is added to verifyDelay to form the elapsed-time
// threshold (6s) that confirms a sleep actually took effect.
const verifyConfirmMargin = 1 * time.Second

// maxSleepRetries caps consecutive failed sleep attempts, so a machine that
// simply will not sleep does not keep flashing the screen off all night.
const maxSleepRetries = 3

// Delay returns the current re-sleep delay in milliseconds.
func Delay() int { return int(delayMS.Load()) }

// SetDelay updates the re-sleep delay in milliseconds.
func SetDelay(ms int) { delayMS.Store(int32(ms)) }

// stateQuery returns the current system state: whether the system is suspended
// and whether the lid is closed. Injected by the power package; stored in an
// atomic.Value because it is written on the power thread and read on the engine
// loop goroutine.
var stateQuery atomic.Value // func() (suspended, lidClosed bool)

// SetStateQuery injects the state query function from the power package.
func SetStateQuery(fn func() (suspended, lidClosed bool)) { stateQuery.Store(fn) }

// queryState returns the injected query function and whether it is available.
func queryState() (fn func() (suspended, lidClosed bool), ok bool) {
	fn, ok = stateQuery.Load().(func() (suspended, lidClosed bool))
	return fn, ok
}

// event is a single input to the engine loop.
type event struct{ kind eventKind }

type eventKind uint8

const (
	evTrigger  eventKind = iota // any power event: evaluate and act
	evSleepNow                  // manual immediate sleep (test)
)

var (
	started atomic.Bool
	mu      sync.Mutex // serializes Start/Stop so a concurrent Stop can never close a nil/stale stopCh
	stopCh  chan struct{}
	evCh    = make(chan event, 16)
)

// Running reports whether the engine is guarding.
func Running() bool { return started.Load() }

// Start begins guarding: launches the event loop goroutine and immediately
// evaluates the current state so a lid that is already closed is armed right
// away (the initial power notification fires before the engine is started).
func Start() {
	mu.Lock()
	defer mu.Unlock()
	if started.Load() {
		return
	}
	stopCh = make(chan struct{})
	started.Store(true)
	go runLoop()
	postEvent(evTrigger)
}

// Stop stops guarding, cancelling any pending re-sleep; returns whether it was
// actually running.
func Stop() bool {
	mu.Lock()
	defer mu.Unlock()
	if !started.Load() {
		return false
	}
	started.Store(false)
	close(stopCh)
	stopCh = nil
	return true
}

// Trigger posts a power event (suspend, resume, or lid change). The engine
// queries the state and decides what to do. Non-blocking and dropped while the
// engine is stopped, so a full channel can never block a caller (e.g. the power
// message loop or the GUI thread).
func Trigger() {
	if !started.Load() {
		return
	}
	postEvent(evTrigger)
}

// SleepNow posts a manual immediate-sleep request (test). When the engine is
// not guarding, the sleep is executed directly so the display-off mechanism can
// still be verified.
func SleepNow() {
	if !started.Load() {
		log.Print(i18n.T("log.sepTest"))
		doSleep()
		return
	}
	postEvent(evSleepNow)
}

// postEvent enqueues an event without blocking; events are dropped when the
// channel is full (the engine loop is backed up), never stalling callers.
func postEvent(kind eventKind) {
	select {
	case evCh <- event{kind: kind}:
	default:
	}
}

// appState is the engine's internal state, touched only on the loop goroutine.
type appState struct {
	lidClosed bool      // last known lid state (true = closed)
	sleeping  bool      // a sleep transition we initiated is in progress
	pending   bool      // a re-sleep timer is armed and not yet fired
	verify    bool      // a sleep-verify timer is armed (sleep attempt in flight)
	verifyAt  time.Time // wall-clock time when the verify timer was armed
	retries   int       // consecutive sleep attempts that did not take effect
}

var st appState

// runLoop processes events serially and manages the debounced re-sleep timer.
func runLoop() {
	timer := time.NewTimer(time.Hour)
	if !timer.Stop() {
		<-timer.C
	}
	defer timer.Stop()

	for {
		select {
		case <-stopCh:
			return
		case <-timer.C:
			if st.verify {
				onVerifyFired(timer)
			} else {
				onTimerFired(timer)
			}
		case e := <-evCh:
			handle(e, timer)
		}
	}
}

// handle processes a single event.
func handle(e event, timer *time.Timer) {
	switch e.kind {
	case evTrigger:
		onTrigger(timer)
	case evSleepNow:
		log.Print(i18n.T("log.sepTest"))
		sleep(timer)
	}
}

// onTrigger evaluates the current system state: record suspend/resume and the
// lid state, then schedule a re-sleep when the lid is closed or cancel any
// pending one when it is open.
func onTrigger(timer *time.Timer) {
	query, ok := queryState()
	if !ok {
		return
	}
	suspended, closed := query()

	st.sleeping = suspended
	if suspended {
		st.retries = 0 // the sleep attempt succeeded
		stopPending(timer)
		return
	}

	if closed != st.lidClosed {
		st.lidClosed = closed
		if closed {
			log.Event(i18n.T("log.evtLidClosed"))
		} else {
			log.Event(i18n.T("log.evtLidOpened"))
			st.retries = 0
		}
	}

	if closed {
		armTimer(timer)
	} else {
		stopTimer(timer)
	}
}

// onTimerFired runs when the re-sleep delay elapses: sleep only if the lid is
// still closed and no sleep is already in progress.
func onTimerFired(timer *time.Timer) {
	st.pending = false
	if !st.lidClosed || st.sleeping || !started.Load() {
		return
	}
	sleep(timer)
}

// onVerifyFired runs when the sleep-verify timer fires. The timer can only fire
// while the system is awake, so we compare the wall-clock elapsed time since
// the sleep attempt against the verify window:
//
//   - elapsed > verifyDelay+margin (6s): the timer was delayed by an actual
//     suspend, i.e. the sleep took effect; re-run the normal flow (no retry).
//   - elapsed ~= verifyDelay: the system stayed awake, the sleep did not take
//     effect; retry up to maxSleepRetries times.
func onVerifyFired(timer *time.Timer) {
	st.verify = false
	if !st.lidClosed || !started.Load() {
		return
	}

	if time.Since(st.verifyAt) > verifyDelay+verifyConfirmMargin {
		log.Action(i18n.T("log.actSleepConfirmed"))
		st.retries = 0
		onTrigger(timer)
		return
	}

	st.retries++
	if st.retries > maxSleepRetries {
		log.Error(i18n.T("log.actSleepRetryGiveUp"))
		return
	}
	log.Action(i18n.F("log.actSleepRetry", st.retries, maxSleepRetries))
	onTrigger(timer)
}

// armTimer (re)arms the debounced delay timer: an existing pending timer is
// cancelled first, so rapid events restart the delay (debounce).
func armTimer(timer *time.Timer) {
	wasPending := st.pending
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
	timer.Reset(time.Duration(delayMS.Load()) * time.Millisecond)
	st.pending = true
	st.verify = false

	ms := int(delayMS.Load())
	if wasPending {
		log.Action(i18n.F("log.actReschedule", ms))
	} else {
		log.Action(i18n.F("log.actSchedule", ms))
	}
}

// armVerify arms the sleep-verify timer that fires verifyDelay after a sleep
// attempt, to detect a sleep that did not take effect. The armed time is
// recorded with the monotonic clock stripped (Round(0)) so the later elapsed
// comparison uses the wall clock, which keeps advancing while the system is
// suspended.
func armVerify(timer *time.Timer) {
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
	timer.Reset(verifyDelay)
	st.pending = false
	st.verify = true
	st.verifyAt = time.Now().Round(0)
}

// stopPending cancels only the pending re-sleep, leaving a sleep-verify timer
// armed: on suspend the verify timer is deliberately kept so that, when it
// fires after a resume, the elapsed-time check can confirm the sleep took
// effect. The two never share the timer (pending and verify are exclusive).
func stopPending(timer *time.Timer) {
	if !st.pending {
		return
	}
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
	st.pending = false
}

// stopTimer cancels a pending re-sleep or sleep-verify timer.
func stopTimer(timer *time.Timer) {
	if !st.pending && !st.verify {
		return
	}
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
	st.pending = false
	st.verify = false
	log.Action(i18n.T("log.actCancel"))
}

// sleep executes the sleep operation: turn the display off to trigger S0
// Modern Standby. It is re-entrancy-safe: no-op while a sleep is in progress.
func sleep(timer *time.Timer) {
	if !started.Load() {
		return
	}
	if st.sleeping {
		return
	}
	st.sleeping = true
	doSleep()
	armVerify(timer)
}

// doSleep logs and executes the display-off operation. Only the engine loop
// goroutine calls it while guarding; SleepNow calls it directly when not
// guarding (single GUI thread, no state to guard).
func doSleep() {
	log.Action(i18n.T("log.actSleep"))
	win32.SendScreenOff()
}
