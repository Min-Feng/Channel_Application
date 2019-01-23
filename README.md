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
使用**有緩衝**的channel的方式實現  
內部沒有gorountine leaking的問題

## [work](work/work.go)
控制goroutine的數量在指定的範圍內,功能如同golimit  
使用**無緩衝**的channel來實現,需要另外開work goroutine執行Task   
由於內部另外開work goroutine,所以使用完work Pool必須關閉Close()  
不然內部會有gorountine leaking的問題  
Close()方法,直到所有工作都完成才會返回

## [constantNumCall](constantNumCall/main.go)
練習某一個函數最多只能被多少人使用,原理類似 golimit