package main

import (
	"time"
	"fmt"
)

func main() {
	dateStr := "15/Mar/2018:03:14:53 +0800"
	t,_:=time.Parse("2/Jan/2006:15:04:05 -0700",dateStr)
	//print(t)
	fmt.Println(t.Format("2006-01-02 15:04:05"))
	ts := TestSt{}
	fmt.Println(ts.Test)
	fmt.Println(ts.CurrentTime)
}

type TestSt struct {
	Test bool
	CurrentTime time.Time
}

func NewTestSt() *TestSt {
	return &TestSt{Test:true,CurrentTime:time.Now()}
}