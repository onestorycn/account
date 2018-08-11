package main

import (
	"context"
	"fmt"
	grpc2 "google.golang.org/grpc"
	"github.com/processout/grpc-go-pool"
	"time"
	account "account/grpc/proto/account"
)

func main() {
	p, errPoll := grpcpool.New(func() (*grpc2.ClientConn, error) {
		return grpc2.Dial("127.0.0.1:9994", grpc2.WithInsecure())
	}, 1, 100, time.Second)

	if errPoll != nil {
		fmt.Println(errPoll)
	}

	client, err := p.Get(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	cl := account.NewAccountServiceClient(client.ClientConn)
	//Info := &account.FullAccount{PassId:"", Phone:13512415778, Email:"605714010@qq.com", NickName:"tom", Ext:"", OpenId:"", Password:"333333333333", Avatar:""}
	//rsp, err := cl.InsertAccountInfo(context.TODO(), &account.Request{Info:Info})
	rsp, err := cl.LoginAccount(context.TODO(), &account.RequestLogin{Email:"605714010@qq.com", Password:"333333333333"})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("code : %v \n " , rsp.Code)
	fmt.Printf("message : %v \n " , rsp.Message)
	fmt.Printf("data : %v \n " , rsp.Data)
}

