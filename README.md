# ftunnel

### 修改说明
* 公用tunnel
* 建立tunnel之前测试MTU
* 发送端各自收集丢包信息，调整自身TCP read buf size，最小为tunnel MTU
* keepAlive，断链自动重连
* 拆分TCP包发送后，阻塞读取，等待确认包，提高链接稳定性，防止乱序和竞争
* 压缩加密，防止gfw，并减小数据包
* ListenAndServeUDP，ListenAndServeTCP由ftunnel包提供，精简server和client的实现，减少出错