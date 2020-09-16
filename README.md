### Pigeond

Systemd服务，监听Unix scoket用于接受pigeon发送的指令。
日志记录于/var/log/pigeond.log并标准输出被journald
服务收集。

### Pigeon

使用cobra实现cli，并将指令通过Unix socket发送给pigeond。