package explorer

import "time"

func (t *ExplorerHealth) SetSampleRate(sampleRate int) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.sampleRate = sampleRate
}

func (t *ExplorerHealth) GetSampleRate() int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.sampleRate
}

func (t *ExplorerHealth) SetErrors(errors int64) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.errors = errors
}

func (t *ExplorerHealth) GetErrors() int64 {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.errors
}

func (t *ExplorerHealth) SetReceived(received int64) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.received = received
}

func (t *ExplorerHealth) GetReceived() int64 {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.received
}

func (t *ExplorerHealth) SetStartTime(startTime time.Time) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.startTime = startTime
}

func (t *ExplorerHealth) GetStartTime() time.Time {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.startTime
}

func (t *ExplorerHealth) SetUpdatedAt(updatedAt time.Time) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.updatedAt = updatedAt
}

func (t *ExplorerHealth) GetUpdatedAt() time.Time {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.updatedAt
}
