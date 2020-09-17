## 项目目录结构
```
gin-items/
|── api/            ———————— 控制器入口
|   |—— http/       ———————— http控制器
|—— cmd/            ———————— cmd命令入口
|—— conf/           ———————— 配置文件
|—— dao/            ———————— 数据库操作层
|—— helper/         ———————— 辅助方法
|—— library/        ———————— 类库
|   |—— app
|   |—— constant
|   |—— rabbitmq
|   |—— setting
|   |—— token
|—— middleware/     ———————— 中间件
|   |—— jwt
|   |—— log 
|—— migration/      ———————— 数据库迁移文件
|—— model/          ———————— 结构体抽象层
|—— service/        ———————— 业务处理层
|── makefile        ———————— 构建工具
```

## 功能说明
- [x] 项目目录划分
- [x] 如何接收复杂参数？()
  - 如果用传统form-data方式，go无法处理
  - 查了一下，java默认接收方式也是json，通过requestbdy获取
  - 所以直接统一使用json提交，但是无法给默认参数，只能在实例化的时候给默认参数
- [x] 统一返回格式
- [x] 状态码处理
- [x] mysql业务
  - [x] 支持不定查询字段（已废弃）
    - [x] 查询接口是map\[string\]interface{}是无法使用的
  - [x] 获取上一次创建id（会赋值到原来结构体主键里）
  - [x] 批量插入
  - [x] 无法合并参数默认值，比如添加item初始状态（不支持，可以通过初始化方法处理）
  - [x] 未获取数据的状态码处理
  - [x] 协程并发读读取多个数据库
  - [x] 批量更新(不支持类似CI那种updateBatch方法)
- [x] 中间件鉴权
  - [x] jwt
- [x] 多数据库连接
  - [x] 主从分库 
  - [x] 跨库可以通过tablename指定库名
  - [ ] 未测试出连接池的开启的效果
- [x] 区分测试生产环境配置
  - [x] viper
  - [x] 热更新
- [ ] 数据库迁移migration
  - [x] golang-migrate
  - [ ] 无法多数据库区公用一个版本
- [ ] rabbitmq封装
  - [x] 消费者处理
  - [x] 生产者处理
  - [x] 消息绑定处理
  - [x] 协程发送消息
  - [ ] 生产者连接池(初始化无需产生新的连接)
- [x] mq消息处理
  - [x] SyncSkuInsert（批量处理）
  - [x] SyncSkuUpdate
  - [x] SyncIemInsert
  - [x] SyncIemUpdate
- [x] 日志处理（zap OR logrus? 选择zap的优势是性能更高）
  - [x] 访问日志
    - [x] 日志切割（按文件大小lumberjack，日期file-rotatelogs）
  - [x] 业务日志
    - [x] 区分错误级别
- [ ] docker部署
- [ ] 重启(优雅关闭服务)
- [x] 原有bug/待优化
  - [x] 查询商品列表接口，只返回item_id，会导致查询某一个sku的条码，返回了全部sku
  - [x] mq消息发送重复（由于在dao层处理发送消息，会出现同时调用多个dao层），比如删除商品，即发送了同步item消息，也发送了同步sku消息，回调处理方法差不多
- [ ] 多版本
- [ ] 协程mq消费者如何优雅退出
- [ ] 入口文件调整
- [ ] qps性能对比swoft
- [x] makefile构建

