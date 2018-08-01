#### 执行单元
![image](https://chen-wj.oss-cn-hangzhou.aliyuncs.com/youdao-markdown/corm%E6%89%A7%E8%A1%8C%E5%8D%95%E5%85%83.png)
- 一个执行单元对应一个module文件，其中可以包含多个struct
- 通过每一个struct都可以找到所在的执行单元
- 通过执行单元可以找到其包含的每个struct的结构
- 一个执行单元对应一个table
- 一个执行单元对应一个缓冲cache
- 可以以执行单元为单位设置是否开启缓冲机制
- 执行单元中有修改表的操作时清空其缓冲

#### 结构
- builder
- corm
- errorHandle
- exec
- logs
- struct_utils