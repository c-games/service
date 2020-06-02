# Service
Service is basic service package used to build your own service.

### Services List
* Logger go built-in log
* HTTP go built-in http using [gorilla mux](https://github.com/gorilla/mux)
* MQ  [RabbitMQ](https://www.rabbitmq.com/) 
* MySQL [go-mysql](https://github.com/go-sql-driver/mysql)
* Redis [go-redis](https://github.com/go-redis/redis)

### ServiceConfig.json
```
{
	"Version":"1.0",
	"Logger":{
	    "Enable":true,              // enable logger
	    "FullPath":false,           // log with full path or not
		"FileName":"lottery",	//logger file name
		"Level":"info"	 	//default level:"debug", "info", "warning", "error", "panic"
	},
	"Mysql":{
		"Enable":true,  //true=enable, false=disable
		"MainDB":{
			"DriverName":"mysql",
			"User":"user",						//account
			"Password":"pass",				        //password
			"Net":"tcp",   						//net string
			"Address":"localhost:3306",		  	        //ip:port
			"DBName":"",						//db name
			"Timeout":"1m2s",
			"ReadTimeout":"2s",
			"WriteTimeout":"2s"
		},
		"ReadWriteSplitting":false,
		"ReadOnlyDB":{
			"DriverName":"mysql",
			"User":"user",
			"Password":"pass",
			"Net":"tcp",
			"Address":"localhost:3306",
			"DBName":"",
			"Timeout":"1m2s",
			"ReadTimeout":"2s",
			"WriteTimeout":"2s"
		}
	},
	"MQ":{
		"Enable":true,
		"Url":"amqp://account:password@localhost:5672/",  //account:password@ip:port
		"ResponseTimeoutSecond":10,
		"Qos":{
			"PrefetchCount":100,    //comsumer 一次取幾個 message
			"PrefetchSize":0,       //message 的 size, 0=無限制
			"Global":false
		},
		//目前架構至少要有兩個queue
		//一個收 command, 一個收 response
		 "QueueDefinition": [
              {
                "Name": "command",  //queue name
                "Durable": true,      //durable
                "AutoDelete": false,  //if auto delete queue
                "Exclusive": false,   //if exclusive
                "NoWaite": false,     //if no wait
                "Args": {             //queue args a map[string]interface{}
                  "A": 10,
                  "B": 20
                }
              },
              {
                "Name": "response",
                "Durable": true,
                "AutoDelete": false,
                "Exclusive": false,
                "NoWaite": false,
                "Args": {}
              }
            ],
          // command consumer param, 收 command 的 command consumer 參數  
          "CommandConsumerParam": {
                "Queue": "command",
                "Consumer": "",
                "AutoAck": false,
                "Exclusive": false,
                "NoLocal": false,
                "NoWaite": false,
                "Args": {}
              },
          // 收 response 的 consumer 參數    
          "ResponseConsumerParam": {
                "Queue": "response",
                "Consumer": "",
                "AutoAck": false,
                "Exclusive": false,
                "NoLocal": false,
                "NoWaite": false,
                "Args": {}
          }  
	},
	"HTTP":{
		"Enable":true,
		"Address":"",
		"Port":"20000",
		"ReadTimeoutSecond":10,
		"WriteTimeoutSecond":10,
		"IdleTimeoutSecond":10,
		"MaxHeaderBytes":1048576,
		"IsTLS":false,
		"CertificateFile":"lc8168.com.crt",
		"KeyFile":"lc8168.com.key"
	},
	"Redis":{
		"Enable":true,                 //true=enable, false=disable
		"Network":"tcp",               //net default "tcp"
		"Address":"localhost:6379",    //ip:port
		"Password":"",                 //password
		"DB":0,                        //db index, 0-15
		"DialTimeoutSecond":5,
		"ReadTimeoutSecond":3,
		"WriteTimeoutSecond":3,
		"PoolSize":10
	},
	"Others": {
            "str": "abc",
            "int": 123123,
            "float": 123123.45,
            "arr": [1, 2, 3, 4, 5],
            "map": {"test": "1787"}
      }
}
```

 Others Example

```
    go
	testString, err := cfg.GetString("str")
	if err != nil {
		panic("unknown key")
	}
	fmt.Printf("%s\n", testString)

	i, err := cfg.GetInt64("int")
	if err != nil {
		panic("unknown key")
	}
	fmt.Printf("%v\n", i)

	f, err  := cfg.GetFloat64("float")
	if err != nil {
		panic("unknown key")
	}
	fmt.Printf("%v\n", f)

	var arr []int
	err = cfg.GetArray("arr", &arr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", arr)

	s := struct{
		Test string
	}{}
	err = cfg.GetMap("map", &s)
	if err != nil {
		panic("unknown key")
	}
	fmt.Printf("%v\n", s)

```

output

```
abc123123
123123.45
[1 2 3 4 5]
```
