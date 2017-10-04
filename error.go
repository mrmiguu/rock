package rock

import "errors"

func (E *Error) makeW() {
	E.p.w.c = make(chan []byte, E.Len)
	go postIfClient(E.p.w.c, Terror, E.Name)
}

func (E *Error) makeR() {
	E.p.r.c = make(chan []byte, E.Len)
	go getIfClient(E.p.r.c, Terror, E.Name)
}

func (E *Error) makeN() {
	E.p.n.c = make(chan int)
}

func (E *Error) add() {
	errorDict.Lock()
	if errorDict.m == nil {
		errorDict.m = map[string]*Error{}
	}
	if _, found := errorDict.m[E.Name]; !found {
		errorDict.m[E.Name] = E
	}
	errorDict.Unlock()
}

func (E *Error) To(e error) {
	go started.Do(getAndOrPostIfServer)

	E.add()

	E.p.w.Do(E.makeW)
	if IsClient {
		E.p.w.c <- []byte(e.Error())
		return
	}

	E.p.n.Do(E.makeN)
	for {
		<-E.p.n.c
		E.p.w.c <- []byte(e.Error())
		if len(E.p.n.c) == 0 {
			break
		}
	}
}

func (E *Error) From() error {
	go started.Do(getAndOrPostIfServer)

	E.add()

	E.p.r.Do(E.makeR)
	return errors.New(string(<-E.p.r.c))
}
