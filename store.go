package main

import "sync"

type failureInfo struct {
	numFailures int
	failedURLs  []string
}

type store struct {
	lock         sync.RWMutex
	failureCount map[int]*failureInfo // Key: build number
}

func newStore() *store {
	return &store{
		failureCount: make(map[int]*failureInfo),
	}
}

func (s *store) get(buildNumber int) *failureInfo {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.failureCount[buildNumber]
}

func (s *store) push(buildNumber int, jobURL string) *failureInfo {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.failureCount[buildNumber] == nil {
		s.failureCount[buildNumber] = &failureInfo{}
	}
	s.failureCount[buildNumber].numFailures++
	s.failureCount[buildNumber].failedURLs = append(s.failureCount[buildNumber].failedURLs, jobURL)
	return s.failureCount[buildNumber]
}
