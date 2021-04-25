// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/PutOrders.proto

package PutOrders

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/asim/go-micro/v3/api"
	client "github.com/asim/go-micro/v3/client"
	server "github.com/asim/go-micro/v3/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for PutOrders service

func NewPutOrdersEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for PutOrders service

type PutOrdersService interface {
	PutOrders(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	Stream(ctx context.Context, in *StreamingRequest, opts ...client.CallOption) (PutOrders_StreamService, error)
	PingPong(ctx context.Context, opts ...client.CallOption) (PutOrders_PingPongService, error)
}

type putOrdersService struct {
	c    client.Client
	name string
}

func NewPutOrdersService(name string, c client.Client) PutOrdersService {
	return &putOrdersService{
		c:    c,
		name: name,
	}
}

func (c *putOrdersService) PutOrders(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "PutOrders.PutOrders", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *putOrdersService) Stream(ctx context.Context, in *StreamingRequest, opts ...client.CallOption) (PutOrders_StreamService, error) {
	req := c.c.NewRequest(c.name, "PutOrders.Stream", &StreamingRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &putOrdersServiceStream{stream}, nil
}

type PutOrders_StreamService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*StreamingResponse, error)
}

type putOrdersServiceStream struct {
	stream client.Stream
}

func (x *putOrdersServiceStream) Close() error {
	return x.stream.Close()
}

func (x *putOrdersServiceStream) Context() context.Context {
	return x.stream.Context()
}

func (x *putOrdersServiceStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *putOrdersServiceStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *putOrdersServiceStream) Recv() (*StreamingResponse, error) {
	m := new(StreamingResponse)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *putOrdersService) PingPong(ctx context.Context, opts ...client.CallOption) (PutOrders_PingPongService, error) {
	req := c.c.NewRequest(c.name, "PutOrders.PingPong", &Ping{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return &putOrdersServicePingPong{stream}, nil
}

type PutOrders_PingPongService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Ping) error
	Recv() (*Pong, error)
}

type putOrdersServicePingPong struct {
	stream client.Stream
}

func (x *putOrdersServicePingPong) Close() error {
	return x.stream.Close()
}

func (x *putOrdersServicePingPong) Context() context.Context {
	return x.stream.Context()
}

func (x *putOrdersServicePingPong) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *putOrdersServicePingPong) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *putOrdersServicePingPong) Send(m *Ping) error {
	return x.stream.Send(m)
}

func (x *putOrdersServicePingPong) Recv() (*Pong, error) {
	m := new(Pong)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for PutOrders service

type PutOrdersHandler interface {
	PutOrders(context.Context, *Request, *Response) error
	Stream(context.Context, *StreamingRequest, PutOrders_StreamStream) error
	PingPong(context.Context, PutOrders_PingPongStream) error
}

func RegisterPutOrdersHandler(s server.Server, hdlr PutOrdersHandler, opts ...server.HandlerOption) error {
	type putOrders interface {
		PutOrders(ctx context.Context, in *Request, out *Response) error
		Stream(ctx context.Context, stream server.Stream) error
		PingPong(ctx context.Context, stream server.Stream) error
	}
	type PutOrders struct {
		putOrders
	}
	h := &putOrdersHandler{hdlr}
	return s.Handle(s.NewHandler(&PutOrders{h}, opts...))
}

type putOrdersHandler struct {
	PutOrdersHandler
}

func (h *putOrdersHandler) PutOrders(ctx context.Context, in *Request, out *Response) error {
	return h.PutOrdersHandler.PutOrders(ctx, in, out)
}

func (h *putOrdersHandler) Stream(ctx context.Context, stream server.Stream) error {
	m := new(StreamingRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.PutOrdersHandler.Stream(ctx, m, &putOrdersStreamStream{stream})
}

type PutOrders_StreamStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*StreamingResponse) error
}

type putOrdersStreamStream struct {
	stream server.Stream
}

func (x *putOrdersStreamStream) Close() error {
	return x.stream.Close()
}

func (x *putOrdersStreamStream) Context() context.Context {
	return x.stream.Context()
}

func (x *putOrdersStreamStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *putOrdersStreamStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *putOrdersStreamStream) Send(m *StreamingResponse) error {
	return x.stream.Send(m)
}

func (h *putOrdersHandler) PingPong(ctx context.Context, stream server.Stream) error {
	return h.PutOrdersHandler.PingPong(ctx, &putOrdersPingPongStream{stream})
}

type PutOrders_PingPongStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Pong) error
	Recv() (*Ping, error)
}

type putOrdersPingPongStream struct {
	stream server.Stream
}

func (x *putOrdersPingPongStream) Close() error {
	return x.stream.Close()
}

func (x *putOrdersPingPongStream) Context() context.Context {
	return x.stream.Context()
}

func (x *putOrdersPingPongStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *putOrdersPingPongStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *putOrdersPingPongStream) Send(m *Pong) error {
	return x.stream.Send(m)
}

func (x *putOrdersPingPongStream) Recv() (*Ping, error) {
	m := new(Ping)
	if err := x.stream.Recv(m); err != nil {
		return nil, err
	}
	return m, nil
}
