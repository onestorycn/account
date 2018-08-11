package config

import (
	"github.com/micro/go-config/source/file"
	"github.com/micro/go-config"
	"fmt"
	"errors"
	"sync"
)

var cacheConf sync.Map

func GetConfigMap(confName string) (mapRes map[string]interface{}) {
	defer func(){
		if err:=recover();err!=nil{
			mapRes = nil
		}
	}()

	finderKey := fmt.Sprintf("conf-%s", confName)
	if res, ok := cacheConf.Load(finderKey); ok {
		if back, okMap := res.(map[string]interface{}); okMap {
			return back
		}
	}
	confPath := fmt.Sprintf("config/%s.json", confName)
	errLoad := config.Load(file.NewSource(
		file.WithPath(confPath),
	))
	if errLoad != nil {
		fmt.Sprintf("load log log fail =====> %s", errLoad.Error())
		return nil
	}
	confGet := config.Map()
	if len(confGet) > 0 {
		cacheConf.Store(finderKey, confGet)
	}
	return confGet
}

func GetConfigMapObj(confName string, obj interface{}, path ...string) (err error) {
	defer func(){
		if err:=recover();err!=nil{
			err = errors.New("load db config fail")
		}
	}()
	confPath := fmt.Sprintf("config/%s.json", confName)
	errorLoad := config.Load(file.NewSource(
		file.WithPath(confPath),
	))
	if errorLoad != nil {
		fmt.Println("load db config fail " + errorLoad.Error())
		return errors.New("load config fail")
	}
	err = config.Get(path...).Scan(&obj)
	return err
}

func GetServiceEnv(envString string) interface{} {
	confGet := GetConfigMap("env")
	if getVar, ok := confGet[envString]; ok {
		return getVar
	}
	return nil
}


func GetDefaultGrpcConfig(envString, grpcName string, obj interface{}) (err error) {
	defer func(){
		if err:=recover();err!=nil{
			err = errors.New("load db config fail")
		}
	}()
	confPath := "config/grpc.json"
	errorLoad := config.Load(file.NewSource(
		file.WithPath(confPath),
	))
	if errorLoad != nil {
		fmt.Println("load db config fail " + errorLoad.Error())
		return errors.New("load config fail")
	}
	err = config.Get(grpcName, envString).Scan(&obj)
	return err
}