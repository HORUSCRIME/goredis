package database

import (
	"sort"
	"sync"
)

type ZSetMember struct {
	Member string
	Score  float64
}

type ZSet struct {
	mu sync.RWMutex

	members []ZSetMember
	index   map[string]float64 
}

func NewZSet() *ZSet {
	return &ZSet{
		members: make([]ZSetMember, 0),
		index:   make(map[string]float64),
	}
}

func (z *ZSet) Type() string {
	return "zset"
}

func (z *ZSet) ZAdd(score float64, member string) int {
	z.mu.Lock()
	defer z.mu.Unlock()

	if oldScore, ok := z.index[member]; ok {
		if oldScore == score {
			return 0
		}
		for i, m := range z.members {
			if m.Member == member {
				z.members = append(z.members[:i], z.members[i+1:]...)
				break
			}
		}
	}
	z.members = append(z.members, ZSetMember{Member: member, Score: score})
	z.index[member] = score

	sort.Slice(z.members, func(i, j int) bool {
		return z.members[i].Score < z.members[j].Score
	})
	return 1
}

func (z *ZSet) ZScore(member string) (float64, bool) {
	z.mu.RLock()
	defer z.mu.RUnlock()
	score, ok := z.index[member]
	return score, ok
}

func (z *ZSet) ZRem(members ...string) int {
	z.mu.Lock()
	defer z.mu.Unlock()

	removedCount := 0
	for _, memberToRemove := range members {
		if _, ok := z.index[memberToRemove]; ok {
			for i, m := range z.members {
				if m.Member == memberToRemove {
					z.members = append(z.members[:i], z.members[i+1:]...)
					break
				}
			}
			delete(z.index, memberToRemove)
			removedCount++
		}
	}
	return removedCount
}

func (z *ZSet) ZCard() int {
	z.mu.RLock()
	defer z.mu.RUnlock()
	return len(z.members)
}