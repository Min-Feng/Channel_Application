# Channel_Application
個人利用golang中的特性channel,所做的練習程式

## [runner](runner/runner.go)
使用通道來監視程序的執行時間,如果執行時間太長,則終止程序

## [pool](pool/pool.go)
利用有緩衝的channel來實現資源池,管理任意數量的goroutine之間  
共享使用的資源,在共享資料庫連接或記憶體緩衝下非常有用

## [golimit](golimit/golimit.go)
控制goroutine的數量在指定的範圍內,數量太多,會有問題  
例如系統資源耗盡導致panic，或者CPU使用率過高

## [constantNumCall](constantNumCall/main.go)
練習某一個函數最多只能被多少人使用,原理類似 golimit