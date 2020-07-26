package timewheel

import (
	"container/list"
	"time"
)

var (
	Round1Bits  = 8
	Round1Size  = 1 << Round1Bits
	Round1Mask  = Round1Size - 1

	RoundNBits  = 6
	RoundNSize  = 1 << RoundNBits
	RoundNMask  = RoundNSize - 1

	TotalRound  = 5
)

type TimeWheel struct {
	currTick int64
	interval time.Duration

	addTimerChannel chan func()
	rounds [][]*list.List
}

func NewTimeWheel(interval time.Duration) *TimeWheel {
	tw := &TimeWheel{
		interval:interval,
		currTick:time.Now().UnixNano() / int64(interval),
		addTimerChannel:make(chan func(), 1024),
	}

	tw.rounds = make([][]*list.List, TotalRound, TotalRound)
	for i:=0;i<TotalRound;i++{
		slotNum := RoundNSize
		if i == 0 {
			slotNum = Round1Size
		}
		round := make([]*list.List, slotNum, slotNum)
		for k:=0;k<slotNum;k++{
			round[k] = list.New()
		}

		tw.rounds[i] = round
	}
	go tw.run()

	return tw
}

func (tw *TimeWheel) AddTimer(duration time.Duration, callback func(interface{}), arg interface{}) {
	tw.addTimerChannel<-func(){
		timeout := int64(duration / tw.interval)
		if timeout <= 0 {
			timeout = 1
		}
		expire := tw.currTick + timeout
		n := newNode(expire, callback, arg)
		tw.addNode(n)
	}
}

func (tw *TimeWheel) AddTimerDeadline(t time.Time, callback func(interface{}), arg interface{}) {
	tw.addTimerChannel<-func(){
		expire := tw.GetTick(t)
		if expire <= tw.currTick {
			expire = tw.currTick + 1
		}
		n := newNode(expire, callback, arg)
		tw.addNode(n)
	}
}

func (tw *TimeWheel) AddTimerTick(timeout int64, callback func(interface{}), arg interface{}) {
	tw.addTimerChannel<-func(){
		expire := tw.currTick + timeout
		if expire <= tw.currTick {
			expire = tw.currTick + 1
		}
		n := newNode(expire, callback, arg)
		tw.addNode(n)
	}
}

func (tw *TimeWheel) AddTimerTickDeadline(deadline int64, callback func(interface{}), arg interface{}) {
	tw.addTimerChannel<-func(){
		expire := deadline
		if expire <= tw.currTick {
			expire = tw.currTick + 1
		}
		n := newNode(expire, callback, arg)
		tw.addNode(n)
	}
}

func (tw *TimeWheel) GetInterval() time.Duration {
	return tw.interval
}

func (tw *TimeWheel) run() {
	tk := time.NewTicker(tw.interval)
	for {
		select {
		case now := <-tk.C:
			now = time.Now()
			nowTick := tw.GetTick(now)
			for i:=tw.currTick+1; i <= nowTick; i++ {
				tw.tick()
			}
		case addFunc := <-tw.addTimerChannel:
			addFunc()
		}
	}
}

func (tw *TimeWheel) tick() {
	tw.currTick++
	index := tw.getRoundIndex(tw.currTick, 0)
	if index == 0 {
		idx := 0
		for round:=1;idx==0&&round<TotalRound;idx, round = tw.getRoundIndex(tw.currTick, round), round+1{
			idx = tw.getRoundIndex(tw.currTick, round)
			tw.cascade(round, idx)
		}
	}

	tList := tw.rounds[0][index]
	for e := tList.Front(); e != nil; e = tList.Front() {
		n := e.Value.(*node)
		tList.Remove(e)
		if n.callback != nil {
			n.callback(n.arg)
		}
		n.release()
	}
}

func (tw *TimeWheel) cascade(round int, idx int) {
	tList := tw.rounds[round][idx]
	for e :=tList.Front(); e != nil; e = tList.Front() {
		n := e.Value.(*node)
		tList.Remove(e)
		tw.addNode(n)
	}
}

func (tw *TimeWheel) addNode(n *node) {
	idx := n.expire - tw.currTick
	var tList *list.List
	if idx < int64(Round1Size) {
		tList = tw.rounds[0][tw.getRoundIndex(n.expire, 0)]
	} else {
		for round:=1;round<5;round++{
			size := int64(1)<<(Round1Bits+round*RoundNBits)
			if idx < size {
				tList = tw.rounds[round][tw.getRoundIndex(n.expire, round)]
				break
			}
		}
	}
	if tList != nil {
		tList.PushBack(n)
	}
}

func (tw *TimeWheel) GetTick(t time.Time) int64 {
	return t.UnixNano() / int64(tw.interval)
}

func (tw *TimeWheel) getRoundIndex(tick int64, round int) int {
	if round == 0 {
		return int(tick & int64(Round1Mask))
	}
	return int((tick >> int64(Round1Bits + (round-1)*RoundNBits)) & int64(RoundNMask))
}




