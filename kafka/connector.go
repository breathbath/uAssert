package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"strings"
)

type Connector struct {
	conns map[Address]*kafka.Conn
	generalConn *kafka.Conn
	connStr  string
}

func NewConnector(connStr string) *Connector {
	return &Connector{
		conns: make(map[Address]*kafka.Conn),
		connStr:  connStr,
	}
}

func (kc *Connector) GetConn(addr Address, cont context.Context) (*kafka.Conn, error) {
	conn, ok := kc.conns[addr]
	if ok {
		return conn, nil
	}

	var err error

	kc.conns[addr], err = kafka.DialLeader(
		cont,
		"tcp",
		kc.connStr,
		addr.Topic,
		addr.Partition,
	)

	return kc.conns[addr], err
}

func (kc *Connector) GetGeneralConn() (*kafka.Conn, error) {
	if kc.generalConn != nil {
		return kc.generalConn, nil
	}
	var err error
	kc.generalConn, err = kafka.Dial("tcp", kc.connStr)
	if err != nil {
		return nil, err
	}

	return kc.generalConn, nil
}

func (kc *Connector) Destroy() error {
	errs := []string{}
	for _, conn := range kc.conns {
		err := conn.Close()
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	if kc.generalConn != nil {
		err := kc.generalConn.Close()
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return fmt.Errorf("Disconnection failure: %s", strings.Join(errs, ", "))
}
