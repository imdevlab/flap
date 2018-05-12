package main

import (
	"fmt"
	"net"
	"sync"

	"github.com/meqio/meq/proto"
)

func pub(conns []net.Conn) {
	wg := &sync.WaitGroup{}
	wg.Add(len(conns))
	for i, conn := range conns {
		go func(i int, conn net.Conn) {
			defer wg.Done()
			n := 1
			cache := make([]proto.Message, 0, 100)
			for {
				if n > 100000 {
					break
				}
				// 27
				m := proto.Message{
					ID:      []byte(fmt.Sprintf("%d-%010d", i, n)),
					Topic:   []byte(topic),
					Payload: []byte("123456789"),
					Type:    1,
				}
				if len(cache) < 500 {
					cache = append(cache, m)
				} else {
					cache = append(cache, m)
					msg := proto.PackMsgs(cache, proto.MSG_PUB)
					_, err := conn.Write(msg)
					if err != nil {
						panic(err)
					}
					cache = cache[:0]
				}
				n++
			}
		}(i, conn)
	}

	wg.Wait()
}

func pubTimer(conn net.Conn) {

	m := proto.TimerMsg{
		ID:      []byte(fmt.Sprintf("%010d", 1)),
		Topic:   []byte(topic),
		Payload: []byte("1234567891234567"),
		Delay:   30,
		Count:   2,
		Base:    10,
		Power:   1,
	}
	msg := proto.PackTimerMsg(&m)
	_, err := conn.Write(msg)
	if err != nil {
		panic(err)
	}

}