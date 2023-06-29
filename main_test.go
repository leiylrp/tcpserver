package main

import (
	"encoding/base64"
	"fmt"
	"testing"
)

//var scheduler *timer.TimerScheduler

//func init() {
//	scheduler = timer.NewAutoExecTimerScheduler()
//}

//func TestTimerScheduler(t *testing.T) {
//	f1 := timer.NewDelayFunc(func(v ...interface{}) {
//		log.Println("test timer scheduler aaa")
//	}, []interface{}{})
//	f2 := timer.NewDelayFunc(func(v ...interface{}) {
//		log.Println("test timer scheduler bbb")
//	}, []interface{}{})
//
//	timerId1, _ := scheduler.CreateTimerAfter(f1, time.Duration(3)*time.Second)
//	timerId2, _ := scheduler.CreateTimerAfter(f2, time.Duration(5)*time.Second)
//	log.Printf("timerId=%d, timerId2=%d\n", timerId1, timerId2)
//	scheduler.CancelTimer(timerId1)
//}

func TestGOMAXPROCS(t *testing.T) {
	//fmt.Println("p numbers:", runtime.GOMAXPROCS(0))
	fmt.Println("hello world1")
}

func TestURL64Encoding(t *testing.T) {
	data := "hello world"
	enc_str := base64.URLEncoding.EncodeToString([]byte(data))
	fmt.Println("encode:", enc_str)

	dec_str, _ := base64.URLEncoding.DecodeString(enc_str)
	fmt.Println("decode:", dec_str, string(dec_str))
}

