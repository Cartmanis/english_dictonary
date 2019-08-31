package mailer

import (
	"context"
	"fmt"
	"github.com/cartmanis/english_dictonary/backend/app/cmd/lg"
	pb "github.com/cartmanis/english_dictonary/backend/app/cmd/mailer/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

//Структура нашего gRPC сервера
type server struct {
	port int
}

func NewGrpcServer(port int) *server {
	return &server{port: port}
}

func (s *server) SendMail(ctx context.Context, input *pb.MsgRequest) (*pb.Result, error) {

	//В переменную m считываем MsgRequest(смотрим в mail.proto, чтобы вспомнить, что это).
	m := Message{fmt.Sprintf("%v", cnf.from),
		input.To, input.Subject, input.Msg, tplActivate}

	if err := messageLoop(m); err != nil {
		return &pb.Result{Sent: true}, err
	}

	//Ну, а если все хорошо,  отвечаем клиенту true
	return &pb.Result{Sent: true}, nil
}

func (s *server) RetrievePass(ctx context.Context, input *pb.MsgRequest) (*pb.Result, error) {
	m := Message{fmt.Sprintf("%v", cnf.from),
		input.To, input.Subject, input.Msg, tplRetrieve}

	if err := messageLoop(m); err != nil {
		return &pb.Result{Sent: true}, err
	}

	//Ну, а если все хорошо,  отвечаем клиенту true
	return &pb.Result{Sent: true}, nil
}

func (s *server) StartGrpcSever() {
	if s == nil || s.port == 0 {
		lg.Errorf("Не удалось запустить grpc сервер, так как порт не проиницилизирован")
		return
	}
	listner, err := net.Listen("tcp", fmt.Sprintf(":%v", s.port))
	if err != nil {
		lg.Errorf("Не удалось запустить grpc сервер на порту %v. Ошибка:%v", s.port, err)
		return
	}

	grpcServ := grpc.NewServer()

	//Регистрируем связку сервер + listener
	pb.RegisterMailerServer(grpcServ, &server{})
	reflection.Register(grpcServ)

	lg.Info("запущен grpc сервер на порту:", s.port)
	if err = grpcServ.Serve(listner); err != nil {
		lg.Errorf("Не удалось запустить grpc сервер на порту %v. Ошибка:%v", s.port, err)
		return
	}

}
