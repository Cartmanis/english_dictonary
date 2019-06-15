package mailer

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/cartmanis/english_dictonary/backend/app/cmd/mailer/proto"
	"google.golang.org/grpc"
)

func ConnectClient(ip string, port int) error {

	//Открываем соединение, grpc.WithInsecure() означает,
	//что шифрование не используется
	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", ip, port), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	/*

	   Создаём нового клиента, используя соединение conn
	   Обратим внимание на название клиента и на название сервиса,
	   которое мы определили в proto-файле:

	   service Mailer {
	   rpc SendMail(MsgRequest) returns (MsgReply) {}
	   rpc RetrievePass(MsgRequest) returns (MsgReply) {}
	   }

	*/

	c := pb.NewMailerClient(conn)

	//Определяем контекст с таймаутом в 1 секунду
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	/*
	   Шлём запрос 1, ожидаем получение true в структуру rply
	   типа MsgReply, определённую в прото-файле как:

	   message MsgReply {
	   bool sent = 1;
	   }

	*/

	rply, err := c.SendMail(ctx, &pb.MsgRequest{To: "VShmelcer@gmail.com", Msg: "Текст письма"})
	if err != nil {
		return err
	}
	log.Println(rply.Sent)

	//Шлём запрос 2, ожидаем false
	rply, err = c.RetrievePass(ctx, &pb.MsgRequest{To: "jhaine1984@gmail.com", Msg: "Текст письма2"})
	if err != nil {
		return err
	}
	log.Println(rply.Sent)
	return nil
}
