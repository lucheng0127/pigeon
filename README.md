### Pigeond

Systemd服务，监听Unix scoket用于接受pigeon发送的指令。
日志记录于/var/log/pigeond.log并标准输出被journald
服务收集。

### Pigeon

使用cobra实现cli，并将指令通过Unix socket发送给pigeond。

#### modules/tasks

pigeon通过Unix socket给pigeond发送指令，tasks模块从
socket server中获取到指令后执行对于的task，并返回task
执行结果。执行结果通过socket connection返回给pigeon。
Task执行结果以json格式返回。结果通过pigeon/cmd/utils.go
中的checkJSONRst将json数据解析成map。
scripts清单数据保存与/var/run/pigeon/scirpt_inventory.csv

指令格式：

> [AutoAck] [Task] [Arg1] [Arg2] ... [END]

* AutoAck 是否自动应答 T/F 设置成T的话无论指令执行成功与否exit_code都会返回0
* Task 执行的任务，需要在modules/tasks/tasks.go里面的taskProxy中注册
* Args 参数
* END 标志指令结束

TODO:

1. 处理Socket连接时，定义一个socket连接池
2. 重构错误处理和日志收集