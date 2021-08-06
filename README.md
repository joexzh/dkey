# dkey 重复模拟键盘按键

**主要用于个人玩Diablo3**, ~~已弃坑好多年~~

按`5`开始, <code>\`</code>停止, `ctrl+shift+q`退出程序

配置文件放在程序当前目录, 以`config_`开头

例`config_example.yaml`:
```yaml
1:
  mi: 100 # 每100毫秒一次
  action: down # 按键按下
4:
  mi: 300
  action: downup # 按键按下并松开
2:
  mi: 110000
  action: downup
3:
  mi: 110000
  action: downup
```

配置了自动按键`1`,`2`,`3`,`4`, 所有按键同步进行.

创建快捷方式, 目标设置成`/path/to/dkey.exe -f example`使用配置文件`config_example.yaml`启动程序