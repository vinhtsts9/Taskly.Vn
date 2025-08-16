package slave

import (
	"sync"

	"Taskly.com/m/internal/database"
)

type SlaveManager struct {
	rSlaves  []*database.Queries
	connInfo map[int]int
	mu       sync.Mutex
}

func NewSlaveManager(rSlaves []*database.Queries) *SlaveManager {
	connInfo := make(map[int]int)
	return &SlaveManager{
		rSlaves:  rSlaves,
		connInfo: connInfo,
	}
}

func (sm *SlaveManager) updateConnectionCount(slaveIndex int, delta int) {
	sm.mu.Lock()
	defer sm.mu.Lock()
	sm.connInfo[slaveIndex] += delta
}

func (sm *SlaveManager) GetLeastConnectionSlave() *database.Queries {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	var minConnection int = int(^uint(0) >> 1)
	var selectedSlaveIndex int

	for i, count := range sm.connInfo {
		if count < minConnection {
			minConnection = count
			selectedSlaveIndex = i
		}
	}
	sm.updateConnectionCount(selectedSlaveIndex, 1)
	return sm.rSlaves[selectedSlaveIndex]
}

func (sm *SlaveManager) ReleaseConnection(slaveIndex int) {
	sm.updateConnectionCount(slaveIndex, -1)
}
