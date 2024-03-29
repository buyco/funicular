// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/client/amqp.go

// Package mock_source is a generated GoMock package.
package mock_source

import (
	tls "crypto/tls"
	net "net"
	reflect "reflect"

	client "github.com/buyco/funicular/pkg/client"
	gomock "github.com/golang/mock/gomock"
	amqp091 "github.com/rabbitmq/amqp091-go"
)

// MockAMQPConnection is a mock of AMQPConnection interface.
type MockAMQPConnection struct {
	ctrl     *gomock.Controller
	recorder *MockAMQPConnectionMockRecorder
}

// MockAMQPConnectionMockRecorder is the mock recorder for MockAMQPConnection.
type MockAMQPConnectionMockRecorder struct {
	mock *MockAMQPConnection
}

// NewMockAMQPConnection creates a new mock instance.
func NewMockAMQPConnection(ctrl *gomock.Controller) *MockAMQPConnection {
	mock := &MockAMQPConnection{ctrl: ctrl}
	mock.recorder = &MockAMQPConnectionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAMQPConnection) EXPECT() *MockAMQPConnectionMockRecorder {
	return m.recorder
}

// Channel mocks base method.
func (m *MockAMQPConnection) Channel() (client.AMQPChannel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Channel")
	ret0, _ := ret[0].(client.AMQPChannel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Channel indicates an expected call of Channel.
func (mr *MockAMQPConnectionMockRecorder) Channel() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Channel", reflect.TypeOf((*MockAMQPConnection)(nil).Channel))
}

// Close mocks base method.
func (m *MockAMQPConnection) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockAMQPConnectionMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockAMQPConnection)(nil).Close))
}

// ConnectionState mocks base method.
func (m *MockAMQPConnection) ConnectionState() tls.ConnectionState {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConnectionState")
	ret0, _ := ret[0].(tls.ConnectionState)
	return ret0
}

// ConnectionState indicates an expected call of ConnectionState.
func (mr *MockAMQPConnectionMockRecorder) ConnectionState() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnectionState", reflect.TypeOf((*MockAMQPConnection)(nil).ConnectionState))
}

// IsClosed mocks base method.
func (m *MockAMQPConnection) IsClosed() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsClosed")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsClosed indicates an expected call of IsClosed.
func (mr *MockAMQPConnectionMockRecorder) IsClosed() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsClosed", reflect.TypeOf((*MockAMQPConnection)(nil).IsClosed))
}

// LocalAddr mocks base method.
func (m *MockAMQPConnection) LocalAddr() net.Addr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LocalAddr")
	ret0, _ := ret[0].(net.Addr)
	return ret0
}

// LocalAddr indicates an expected call of LocalAddr.
func (mr *MockAMQPConnectionMockRecorder) LocalAddr() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LocalAddr", reflect.TypeOf((*MockAMQPConnection)(nil).LocalAddr))
}

// NotifyBlocked mocks base method.
func (m *MockAMQPConnection) NotifyBlocked(receiver chan amqp091.Blocking) chan amqp091.Blocking {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NotifyBlocked", receiver)
	ret0, _ := ret[0].(chan amqp091.Blocking)
	return ret0
}

// NotifyBlocked indicates an expected call of NotifyBlocked.
func (mr *MockAMQPConnectionMockRecorder) NotifyBlocked(receiver interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NotifyBlocked", reflect.TypeOf((*MockAMQPConnection)(nil).NotifyBlocked), receiver)
}

// NotifyClose mocks base method.
func (m *MockAMQPConnection) NotifyClose(receiver chan *amqp091.Error) chan *amqp091.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NotifyClose", receiver)
	ret0, _ := ret[0].(chan *amqp091.Error)
	return ret0
}

// NotifyClose indicates an expected call of NotifyClose.
func (mr *MockAMQPConnectionMockRecorder) NotifyClose(receiver interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NotifyClose", reflect.TypeOf((*MockAMQPConnection)(nil).NotifyClose), receiver)
}

// MockAMQPChannel is a mock of AMQPChannel interface.
type MockAMQPChannel struct {
	ctrl     *gomock.Controller
	recorder *MockAMQPChannelMockRecorder
}

// MockAMQPChannelMockRecorder is the mock recorder for MockAMQPChannel.
type MockAMQPChannelMockRecorder struct {
	mock *MockAMQPChannel
}

// NewMockAMQPChannel creates a new mock instance.
func NewMockAMQPChannel(ctrl *gomock.Controller) *MockAMQPChannel {
	mock := &MockAMQPChannel{ctrl: ctrl}
	mock.recorder = &MockAMQPChannelMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAMQPChannel) EXPECT() *MockAMQPChannelMockRecorder {
	return m.recorder
}

// Ack mocks base method.
func (m *MockAMQPChannel) Ack(tag uint64, multiple bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ack", tag, multiple)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ack indicates an expected call of Ack.
func (mr *MockAMQPChannelMockRecorder) Ack(tag, multiple interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ack", reflect.TypeOf((*MockAMQPChannel)(nil).Ack), tag, multiple)
}

// Cancel mocks base method.
func (m *MockAMQPChannel) Cancel(consumer string, noWait bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cancel", consumer, noWait)
	ret0, _ := ret[0].(error)
	return ret0
}

// Cancel indicates an expected call of Cancel.
func (mr *MockAMQPChannelMockRecorder) Cancel(consumer, noWait interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cancel", reflect.TypeOf((*MockAMQPChannel)(nil).Cancel), consumer, noWait)
}

// Close mocks base method.
func (m *MockAMQPChannel) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockAMQPChannelMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockAMQPChannel)(nil).Close))
}

// Confirm mocks base method.
func (m *MockAMQPChannel) Confirm(noWait bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Confirm", noWait)
	ret0, _ := ret[0].(error)
	return ret0
}

// Confirm indicates an expected call of Confirm.
func (mr *MockAMQPChannelMockRecorder) Confirm(noWait interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Confirm", reflect.TypeOf((*MockAMQPChannel)(nil).Confirm), noWait)
}

// Consume mocks base method.
func (m *MockAMQPChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp091.Table) (<-chan amqp091.Delivery, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Consume", queue, consumer, autoAck, exclusive, noLocal, noWait, args)
	ret0, _ := ret[0].(<-chan amqp091.Delivery)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Consume indicates an expected call of Consume.
func (mr *MockAMQPChannelMockRecorder) Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Consume", reflect.TypeOf((*MockAMQPChannel)(nil).Consume), queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}

// ExchangeBind mocks base method.
func (m *MockAMQPChannel) ExchangeBind(destination, key, source string, noWait bool, args amqp091.Table) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExchangeBind", destination, key, source, noWait, args)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExchangeBind indicates an expected call of ExchangeBind.
func (mr *MockAMQPChannelMockRecorder) ExchangeBind(destination, key, source, noWait, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExchangeBind", reflect.TypeOf((*MockAMQPChannel)(nil).ExchangeBind), destination, key, source, noWait, args)
}

// ExchangeDeclare mocks base method.
func (m *MockAMQPChannel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp091.Table) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExchangeDeclare", name, kind, durable, autoDelete, internal, noWait, args)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExchangeDeclare indicates an expected call of ExchangeDeclare.
func (mr *MockAMQPChannelMockRecorder) ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExchangeDeclare", reflect.TypeOf((*MockAMQPChannel)(nil).ExchangeDeclare), name, kind, durable, autoDelete, internal, noWait, args)
}

// ExchangeDeclarePassive mocks base method.
func (m *MockAMQPChannel) ExchangeDeclarePassive(name, kind string, durable, autoDelete, internal, noWait bool, args amqp091.Table) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExchangeDeclarePassive", name, kind, durable, autoDelete, internal, noWait, args)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExchangeDeclarePassive indicates an expected call of ExchangeDeclarePassive.
func (mr *MockAMQPChannelMockRecorder) ExchangeDeclarePassive(name, kind, durable, autoDelete, internal, noWait, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExchangeDeclarePassive", reflect.TypeOf((*MockAMQPChannel)(nil).ExchangeDeclarePassive), name, kind, durable, autoDelete, internal, noWait, args)
}

// ExchangeDelete mocks base method.
func (m *MockAMQPChannel) ExchangeDelete(name string, ifUnused, noWait bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExchangeDelete", name, ifUnused, noWait)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExchangeDelete indicates an expected call of ExchangeDelete.
func (mr *MockAMQPChannelMockRecorder) ExchangeDelete(name, ifUnused, noWait interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExchangeDelete", reflect.TypeOf((*MockAMQPChannel)(nil).ExchangeDelete), name, ifUnused, noWait)
}

// ExchangeUnbind mocks base method.
func (m *MockAMQPChannel) ExchangeUnbind(destination, key, source string, noWait bool, args amqp091.Table) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExchangeUnbind", destination, key, source, noWait, args)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExchangeUnbind indicates an expected call of ExchangeUnbind.
func (mr *MockAMQPChannelMockRecorder) ExchangeUnbind(destination, key, source, noWait, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExchangeUnbind", reflect.TypeOf((*MockAMQPChannel)(nil).ExchangeUnbind), destination, key, source, noWait, args)
}

// Flow mocks base method.
func (m *MockAMQPChannel) Flow(active bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Flow", active)
	ret0, _ := ret[0].(error)
	return ret0
}

// Flow indicates an expected call of Flow.
func (mr *MockAMQPChannelMockRecorder) Flow(active interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Flow", reflect.TypeOf((*MockAMQPChannel)(nil).Flow), active)
}

// Get mocks base method.
func (m *MockAMQPChannel) Get(queue string, autoAck bool) (amqp091.Delivery, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", queue, autoAck)
	ret0, _ := ret[0].(amqp091.Delivery)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockAMQPChannelMockRecorder) Get(queue, autoAck interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockAMQPChannel)(nil).Get), queue, autoAck)
}

// GetNextPublishSeqNo mocks base method.
func (m *MockAMQPChannel) GetNextPublishSeqNo() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNextPublishSeqNo")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// GetNextPublishSeqNo indicates an expected call of GetNextPublishSeqNo.
func (mr *MockAMQPChannelMockRecorder) GetNextPublishSeqNo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNextPublishSeqNo", reflect.TypeOf((*MockAMQPChannel)(nil).GetNextPublishSeqNo))
}

// IsClosed mocks base method.
func (m *MockAMQPChannel) IsClosed() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsClosed")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsClosed indicates an expected call of IsClosed.
func (mr *MockAMQPChannelMockRecorder) IsClosed() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsClosed", reflect.TypeOf((*MockAMQPChannel)(nil).IsClosed))
}

// Nack mocks base method.
func (m *MockAMQPChannel) Nack(tag uint64, multiple, requeue bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Nack", tag, multiple, requeue)
	ret0, _ := ret[0].(error)
	return ret0
}

// Nack indicates an expected call of Nack.
func (mr *MockAMQPChannelMockRecorder) Nack(tag, multiple, requeue interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Nack", reflect.TypeOf((*MockAMQPChannel)(nil).Nack), tag, multiple, requeue)
}

// NotifyCancel mocks base method.
func (m *MockAMQPChannel) NotifyCancel(c chan string) chan string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NotifyCancel", c)
	ret0, _ := ret[0].(chan string)
	return ret0
}

// NotifyCancel indicates an expected call of NotifyCancel.
func (mr *MockAMQPChannelMockRecorder) NotifyCancel(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NotifyCancel", reflect.TypeOf((*MockAMQPChannel)(nil).NotifyCancel), c)
}

// NotifyClose mocks base method.
func (m *MockAMQPChannel) NotifyClose(c chan *amqp091.Error) chan *amqp091.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NotifyClose", c)
	ret0, _ := ret[0].(chan *amqp091.Error)
	return ret0
}

// NotifyClose indicates an expected call of NotifyClose.
func (mr *MockAMQPChannelMockRecorder) NotifyClose(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NotifyClose", reflect.TypeOf((*MockAMQPChannel)(nil).NotifyClose), c)
}

// NotifyConfirm mocks base method.
func (m *MockAMQPChannel) NotifyConfirm(ack, nack chan uint64) (chan uint64, chan uint64) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NotifyConfirm", ack, nack)
	ret0, _ := ret[0].(chan uint64)
	ret1, _ := ret[1].(chan uint64)
	return ret0, ret1
}

// NotifyConfirm indicates an expected call of NotifyConfirm.
func (mr *MockAMQPChannelMockRecorder) NotifyConfirm(ack, nack interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NotifyConfirm", reflect.TypeOf((*MockAMQPChannel)(nil).NotifyConfirm), ack, nack)
}

// NotifyFlow mocks base method.
func (m *MockAMQPChannel) NotifyFlow(c chan bool) chan bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NotifyFlow", c)
	ret0, _ := ret[0].(chan bool)
	return ret0
}

// NotifyFlow indicates an expected call of NotifyFlow.
func (mr *MockAMQPChannelMockRecorder) NotifyFlow(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NotifyFlow", reflect.TypeOf((*MockAMQPChannel)(nil).NotifyFlow), c)
}

// NotifyPublish mocks base method.
func (m *MockAMQPChannel) NotifyPublish(confirm chan amqp091.Confirmation) chan amqp091.Confirmation {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NotifyPublish", confirm)
	ret0, _ := ret[0].(chan amqp091.Confirmation)
	return ret0
}

// NotifyPublish indicates an expected call of NotifyPublish.
func (mr *MockAMQPChannelMockRecorder) NotifyPublish(confirm interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NotifyPublish", reflect.TypeOf((*MockAMQPChannel)(nil).NotifyPublish), confirm)
}

// NotifyReturn mocks base method.
func (m *MockAMQPChannel) NotifyReturn(c chan amqp091.Return) chan amqp091.Return {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NotifyReturn", c)
	ret0, _ := ret[0].(chan amqp091.Return)
	return ret0
}

// NotifyReturn indicates an expected call of NotifyReturn.
func (mr *MockAMQPChannelMockRecorder) NotifyReturn(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NotifyReturn", reflect.TypeOf((*MockAMQPChannel)(nil).NotifyReturn), c)
}

// Publish mocks base method.
func (m *MockAMQPChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp091.Publishing) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", exchange, key, mandatory, immediate, msg)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish.
func (mr *MockAMQPChannelMockRecorder) Publish(exchange, key, mandatory, immediate, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockAMQPChannel)(nil).Publish), exchange, key, mandatory, immediate, msg)
}

// PublishWithDeferredConfirm mocks base method.
func (m *MockAMQPChannel) PublishWithDeferredConfirm(exchange, key string, mandatory, immediate bool, msg amqp091.Publishing) (*amqp091.DeferredConfirmation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishWithDeferredConfirm", exchange, key, mandatory, immediate, msg)
	ret0, _ := ret[0].(*amqp091.DeferredConfirmation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PublishWithDeferredConfirm indicates an expected call of PublishWithDeferredConfirm.
func (mr *MockAMQPChannelMockRecorder) PublishWithDeferredConfirm(exchange, key, mandatory, immediate, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishWithDeferredConfirm", reflect.TypeOf((*MockAMQPChannel)(nil).PublishWithDeferredConfirm), exchange, key, mandatory, immediate, msg)
}

// Qos mocks base method.
func (m *MockAMQPChannel) Qos(prefetchCount, prefetchSize int, global bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Qos", prefetchCount, prefetchSize, global)
	ret0, _ := ret[0].(error)
	return ret0
}

// Qos indicates an expected call of Qos.
func (mr *MockAMQPChannelMockRecorder) Qos(prefetchCount, prefetchSize, global interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Qos", reflect.TypeOf((*MockAMQPChannel)(nil).Qos), prefetchCount, prefetchSize, global)
}

// QueueBind mocks base method.
func (m *MockAMQPChannel) QueueBind(name, key, exchange string, noWait bool, args amqp091.Table) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueueBind", name, key, exchange, noWait, args)
	ret0, _ := ret[0].(error)
	return ret0
}

// QueueBind indicates an expected call of QueueBind.
func (mr *MockAMQPChannelMockRecorder) QueueBind(name, key, exchange, noWait, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueueBind", reflect.TypeOf((*MockAMQPChannel)(nil).QueueBind), name, key, exchange, noWait, args)
}

// QueueDeclare mocks base method.
func (m *MockAMQPChannel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp091.Table) (amqp091.Queue, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueueDeclare", name, durable, autoDelete, exclusive, noWait, args)
	ret0, _ := ret[0].(amqp091.Queue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueueDeclare indicates an expected call of QueueDeclare.
func (mr *MockAMQPChannelMockRecorder) QueueDeclare(name, durable, autoDelete, exclusive, noWait, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueueDeclare", reflect.TypeOf((*MockAMQPChannel)(nil).QueueDeclare), name, durable, autoDelete, exclusive, noWait, args)
}

// QueueDeclarePassive mocks base method.
func (m *MockAMQPChannel) QueueDeclarePassive(name string, durable, autoDelete, exclusive, noWait bool, args amqp091.Table) (amqp091.Queue, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueueDeclarePassive", name, durable, autoDelete, exclusive, noWait, args)
	ret0, _ := ret[0].(amqp091.Queue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueueDeclarePassive indicates an expected call of QueueDeclarePassive.
func (mr *MockAMQPChannelMockRecorder) QueueDeclarePassive(name, durable, autoDelete, exclusive, noWait, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueueDeclarePassive", reflect.TypeOf((*MockAMQPChannel)(nil).QueueDeclarePassive), name, durable, autoDelete, exclusive, noWait, args)
}

// QueueDelete mocks base method.
func (m *MockAMQPChannel) QueueDelete(name string, ifUnused, ifEmpty, noWait bool) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueueDelete", name, ifUnused, ifEmpty, noWait)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueueDelete indicates an expected call of QueueDelete.
func (mr *MockAMQPChannelMockRecorder) QueueDelete(name, ifUnused, ifEmpty, noWait interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueueDelete", reflect.TypeOf((*MockAMQPChannel)(nil).QueueDelete), name, ifUnused, ifEmpty, noWait)
}

// QueueInspect mocks base method.
func (m *MockAMQPChannel) QueueInspect(name string) (amqp091.Queue, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueueInspect", name)
	ret0, _ := ret[0].(amqp091.Queue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueueInspect indicates an expected call of QueueInspect.
func (mr *MockAMQPChannelMockRecorder) QueueInspect(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueueInspect", reflect.TypeOf((*MockAMQPChannel)(nil).QueueInspect), name)
}

// QueuePurge mocks base method.
func (m *MockAMQPChannel) QueuePurge(name string, noWait bool) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueuePurge", name, noWait)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueuePurge indicates an expected call of QueuePurge.
func (mr *MockAMQPChannelMockRecorder) QueuePurge(name, noWait interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueuePurge", reflect.TypeOf((*MockAMQPChannel)(nil).QueuePurge), name, noWait)
}

// QueueUnbind mocks base method.
func (m *MockAMQPChannel) QueueUnbind(name, key, exchange string, args amqp091.Table) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueueUnbind", name, key, exchange, args)
	ret0, _ := ret[0].(error)
	return ret0
}

// QueueUnbind indicates an expected call of QueueUnbind.
func (mr *MockAMQPChannelMockRecorder) QueueUnbind(name, key, exchange, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueueUnbind", reflect.TypeOf((*MockAMQPChannel)(nil).QueueUnbind), name, key, exchange, args)
}

// Recover mocks base method.
func (m *MockAMQPChannel) Recover(requeue bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recover", requeue)
	ret0, _ := ret[0].(error)
	return ret0
}

// Recover indicates an expected call of Recover.
func (mr *MockAMQPChannelMockRecorder) Recover(requeue interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recover", reflect.TypeOf((*MockAMQPChannel)(nil).Recover), requeue)
}

// Reject mocks base method.
func (m *MockAMQPChannel) Reject(tag uint64, requeue bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reject", tag, requeue)
	ret0, _ := ret[0].(error)
	return ret0
}

// Reject indicates an expected call of Reject.
func (mr *MockAMQPChannelMockRecorder) Reject(tag, requeue interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reject", reflect.TypeOf((*MockAMQPChannel)(nil).Reject), tag, requeue)
}

// Tx mocks base method.
func (m *MockAMQPChannel) Tx() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tx")
	ret0, _ := ret[0].(error)
	return ret0
}

// Tx indicates an expected call of Tx.
func (mr *MockAMQPChannelMockRecorder) Tx() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tx", reflect.TypeOf((*MockAMQPChannel)(nil).Tx))
}

// TxCommit mocks base method.
func (m *MockAMQPChannel) TxCommit() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TxCommit")
	ret0, _ := ret[0].(error)
	return ret0
}

// TxCommit indicates an expected call of TxCommit.
func (mr *MockAMQPChannelMockRecorder) TxCommit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TxCommit", reflect.TypeOf((*MockAMQPChannel)(nil).TxCommit))
}

// TxRollback mocks base method.
func (m *MockAMQPChannel) TxRollback() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TxRollback")
	ret0, _ := ret[0].(error)
	return ret0
}

// TxRollback indicates an expected call of TxRollback.
func (mr *MockAMQPChannelMockRecorder) TxRollback() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TxRollback", reflect.TypeOf((*MockAMQPChannel)(nil).TxRollback))
}
