package handler

import (
	"context"
	"go.uber.org/zap"
	account "account/grpc/proto/account"
	accModel "account/page/model/account"
	"account/library/logger"
	"account/library"
	"github.com/grpc/grpc-go/status"
	"google.golang.org/grpc/codes"
	"fmt"
	"errors"
	"time"
)

type Account struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Account) InsertAccountInfo(ctx context.Context, req *account.Request, rsp *account.ResponseSafe) (err error) {
	logger.ZapInfo.Info("Received Account.InsertAccountInfo request", zap.String("request", req.String()))
	defer func() {
		if err := recover(); err != nil {
			rsp.Code = library.InternalError
			isError, ok := err.(error)
			if ok {
				rsp.Message = isError.Error()
				err = status.Error(codes.Aborted, isError.Error())
			} else {
				rsp.Message = library.CodeString(library.InternalError)
				err = fmt.Sprintf("%v", err)
			}
		}
	}()
	accountService := accModel.LoadAccountService()
	if accountService == nil {
		panic(errors.New("load account service fail"))
	}
	if library.IsEmpty(req.Info.Password) {
		panic(errors.New("password invalid"))
	}
	userProfile := new(accModel.UserProfile)
	userProfile.Openid = req.Info.OpenId
	userProfile.Passid = req.Info.PassId
	userProfile.Email = req.Info.Email
	userProfile.Avatar = req.Info.Avatar
	userProfile.Password = library.EncodeMd5(req.Info.Password)
	userProfile.Update_time = time.Now().Unix()
	userProfile.Phone = req.Info.Phone
	userProfile.Nick_name = req.Info.NickName

	inserErr := accountService.InsertNewUser(userProfile)
	if inserErr != nil {
		panic(inserErr)
	}
	safeAccount := account.SafeAccount{
		OpenId:   userProfile.Openid,
		PassId:   userProfile.Passid,
		Phone:    userProfile.Phone,
		Ext:      userProfile.Ext,
		Avatar:   userProfile.Avatar,
		NickName: userProfile.Nick_name,
		Email:    userProfile.Email,
	}
	rsp.Data = &safeAccount
	rsp.Code = library.CodeSucc
	rsp.Message = library.CodeString(library.CodeSucc)
	return status.Error(codes.OK, "success")
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Account) UpdateAccountInfo(ctx context.Context, req *account.Request, rsp *account.Response) (err error) {
	logger.ZapInfo.Info("Received Account.UpdateAccountInfo request", zap.String("request", req.String()))
	defer func() {
		if err := recover(); err != nil {
			rsp.Code = library.InternalError
			isError, ok := err.(error)
			if ok {
				rsp.Data = ""
				rsp.Message = isError.Error()
				err = status.Error(codes.Aborted, isError.Error())
			} else {
				rsp.Data = ""
				rsp.Message = library.CodeString(library.InternalError)
				err = fmt.Sprintf("%v", err)
			}
		}
	}()
	accountService := accModel.LoadAccountService()
	if accountService == nil {
		panic(errors.New("load account service fail"))
	}

	if len(req.Info.PassId) != 1 {
		panic(errors.New("get passid empty"))
	}

	affectedRow, inserErr := accountService.UpdateUser(req.Info.Password, req.Info.Email, req.Info.NickName, req.Info.Avatar, req.Info.Ext, req.Info.PassId, req.Info.Phone, true)
	if inserErr != nil {
		panic(inserErr)
	}
	if affectedRow != 1 {
		rsp.Message = "no update record"
	} else {
		rsp.Message = library.CodeString(library.CodeSucc)
	}
	rsp.Code = library.CodeSucc
	return status.Error(codes.OK, "success")
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Account) GetAccountInfo(ctx context.Context, req *account.RequestQuery, rsp *account.ResponseSafe) (err error) {
	logger.ZapInfo.Info("Received Account.GetAccountInfo request", zap.String("request", req.String()))
	defer func() {
		if err := recover(); err != nil {
			rsp.Code = library.InternalError
			isError, ok := err.(error)
			if ok {
				rsp.Message = isError.Error()
				err = status.Error(codes.Aborted, isError.Error())
			} else {
				rsp.Message = library.CodeString(library.InternalError)
				err = fmt.Sprintf("%v", err)
			}
		}
	}()
	accountService := accModel.LoadAccountService()
	if accountService == nil {
		panic(errors.New("load account service fail"))
	}

	var userProfile *accModel.UserProfile
	var userErr error
	if len(req.PassId) > 0 {
		userProfile, userErr = accountService.GetUserById(req.PassId, "")
	} else if len(req.OpenId) > 0 {
		userProfile, userErr = accountService.GetUserById("", req.OpenId)
	} else if req.Phone > 0 {
		userProfile, userErr = accountService.GetUserByPhone(req.Phone)
	} else {
		userErr = errors.New("no valid query option")
	}
	if userErr != nil {
		panic(userErr)
	}
	safeAccount := account.SafeAccount{
		OpenId:   userProfile.Openid,
		PassId:   userProfile.Passid,
		Phone:    userProfile.Phone,
		Ext:      userProfile.Ext,
		Avatar:   userProfile.Avatar,
		NickName: userProfile.Nick_name,
		Email:    userProfile.Email,
	}
	rsp.Data = &safeAccount
	rsp.Code = library.CodeSucc
	rsp.Message = library.CodeString(library.CodeSucc)
	return status.Error(codes.OK, "success")
}

func (e *Account) LoginAccount(ctx context.Context, req *account.RequestLogin, rsp *account.ResponseSafe) (err error) {
	logger.ZapInfo.Info("Received Account.LoginAccount request", zap.String("request", req.String()))
	defer func() {
		if err := recover(); err != nil {
			rsp.Code = library.InternalError
			isError, ok := err.(error)
			if ok {
				rsp.Message = isError.Error()
				err = status.Error(codes.Aborted, isError.Error())
			} else {
				rsp.Message = library.CodeString(library.InternalError)
				err = fmt.Sprintf("%v", err)
			}
		}
	}()
	if len(req.Password) < 1 {
		panic(errors.New("get Password empty"))
	}
	var userProfile *accModel.UserProfile
	var userErr error

	accountService := accModel.LoadAccountService()
	if accountService == nil {
		panic(errors.New("load account service fail"))
	}
	password := library.EncodeMd5(req.Password)
	if req.Phone > 0 {
		userProfile, userErr = accountService.CheckUserLoginByPhone(req.Phone, password)
	}else if len(req.Email) > 0 {
		userProfile, userErr = accountService.CheckUserLoginByEmail(req.Email, password)
	}else{
		panic(errors.New("get login data empty"))
	}
	if userErr != nil {
		panic(userErr)
	}
	safeAccount := account.SafeAccount{
		OpenId:   userProfile.Openid,
		PassId:   userProfile.Passid,
		Phone:    userProfile.Phone,
		Ext:      userProfile.Ext,
		Avatar:   userProfile.Avatar,
		NickName: userProfile.Nick_name,
		Email:    userProfile.Email,
	}
	rsp.Data = &safeAccount
	rsp.Code = library.CodeSucc
	rsp.Message = library.CodeString(library.CodeSucc)
	return status.Error(codes.OK, "success")
}