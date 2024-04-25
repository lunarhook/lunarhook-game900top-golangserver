# Tetris
这是一个给AWS的ELB功能设置的游戏例子</br>
是从网络上一个单机的H5 俄罗斯修改过来的，用来证明AWS的ELB可以支持TCP和Websocket功能 </br>
</br>
先介绍下工程</br>
/views/h5Russia_single.html 单机例子</br>
/views/h5Russia_server.html 游戏主客户端，websocket的操作发起，可以多发起，可以多操控，会出现交叉管理游戏的可能</br>
/views/h5Russia_client.html 游戏从客户端，webscoket的多端接入，只能观看</br>
</br>
游戏进入以后</br>
localhost:8088/game 映射主客户端</br>
localhost:8088/watch 映射从客户端</br>
localhost:8088 映射单机</br>
</br>
系统是beego完成的，因为考虑未来容器化，选择go作为后台语言</br>
</br>
开发过程如下：</br>
aws Assignment 3任务如下：</br>
1、寻找一个h5单机俄罗斯方块</br>
2、修改客户端和服务器代码，服务器用go语言实现，所有运算转入服务器完成（这是游戏持有状态的主要原因，因为需要中央分发）</br>
3、建立服务，建立EC2主机，建立ELB模式设置端口支持http和websocket模式（主要就是用ELB支持websocket，就是同时支持TCP_UDP模式，贯穿8088端口）</br>
4、登陆服务器，安装环境和下载程序编译，执行</br>
5、可以访问游戏，可以通过左右和空格控制游戏，进入画面3秒进入游戏，死亡结束可以刷新，未死亡则需要继续游戏</br>
 </br>
游戏多平台移植工程化是我硕士毕业的方向，具体的游戏还有一些bug，都和网络属性有关，比如TCP粘包，包限制控流，异常处理，比单机复杂，但是总体来说功能基本实现</br>
我会在月底前清理所有相关的云主机配置信息</br>
 </br>
游戏者模式（可以多开，但是只支持第一个用户，后面的都是watch模式但是可以有控制能力）</br>
http://111-5351c565dbdeac4f.elb.ap-southeast-1.amazonaws.com:8088/game  （2020.7.30关闭）</br>
观察者模式（可以多开，无控制能力）</br>
http://111-5351c565dbdeac4f.elb.ap-southeast-1.amazonaws.com:8088/watch （2020.7.30关闭）</br>
单机原始文件（单html）</br>
http://111-5351c565dbdeac4f.elb.ap-southeast-1.amazonaws.com:8088/      （2020.7.30关闭）</br>
github代码</br>
https://github.com/lunarhook/lunarhook-game900top-golangserver</br>
视频录像测试</br>
https://v.youku.com/v_show/id_XNDc2MzI5NzU1Ng==.html</br>
</br>
各位出海的游戏咖门，可以测试下，使用的新加坡方面的主机，要做好还必须优化网络配置什么的，</br>
enjoy:)</br>

