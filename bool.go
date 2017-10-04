package rock

func (B *Bool) makeW() {
	B.p.w.c = make(chan []byte, B.Len)
	go postIfClient(B.p.w.c, Tbool, B.Name)
}

func (B *Bool) makeR() {
	B.p.r.c = make(chan []byte, B.Len)
	go getIfClient(B.p.r.c, Tbool, B.Name)
}

func (B *Bool) makeN() {
	B.p.n.c = make(chan int)
}

func (B *Bool) add() {
	boolDict.Lock()
	if boolDict.m == nil {
		boolDict.m = map[string]*Bool{}
	}
	if _, found := boolDict.m[B.Name]; !found {
		boolDict.m[B.Name] = B
	}
	boolDict.Unlock()
}

func (B *Bool) To(b bool) {
	go started.Do(getAndOrPostIfServer)

	B.add()

	B.p.w.Do(B.makeW)
	if IsClient {
		B.p.w.c <- bool2bytes(b)
		return
	}

	B.p.n.Do(B.makeN)
	for {
		<-B.p.n.c
		B.p.w.c <- bool2bytes(b)
		if len(B.p.n.c) == 0 {
			break
		}
	}
}

func (B *Bool) From() bool {
	go started.Do(getAndOrPostIfServer)

	B.add()

	B.p.r.Do(B.makeR)
	return bytes2bool(<-B.p.r.c)
}
