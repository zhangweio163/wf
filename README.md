# wf
支持参数：  
服务模块：go run server.go -name 服务名称 -port 服务端口号 -reg 注册中心服务地址  
  
注册中心模块:go run RegServer.go -http 注册中心http服务端口 -port 注册中心服务端口  
  
客户端模块 ： go run main.go -name 需调用服务名称(注册中心注册服务名称) -port 客户端http端口 -reg 注册中心服务地址
  
一个可以动态扩展服务的项目，目前情况如下图  
![image](1.png)  
未完成：  
（1）统一认证模块  
（2）统一日志模块  
（3）统一服务网关  

目前实现效果图：  
1、注册中心：  
![image](2.png)  
2、服务启动注册中心发现：  
![image](3.png)  
3、多服务启动注册发现:  
![image](4.png)  
4、客户端请求多服务动态负载响应：
![image](5.png)  
5、只留一个服务业务正常运行:  
![image](6.png)  
6、服务掉线告警:  
![image](7.png)  
7、新服务上线自动注册并提供服务:  
![image](8.png)  
  
  
目前还有很多功能未实现，希望有志之士一起完善