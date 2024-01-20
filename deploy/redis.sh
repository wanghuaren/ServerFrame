#CentOS 8 添加 mirror
$ sudo sed -i -e "s|mirrorlist=|#mirrorlist=|g" /etc/yum.repos.d/CentOS-*;sudo sed -i -e "s|#baseurl=http://mirror.centos.org|baseurl=http://vault.centos.org|g" /etc/yum.repos.d/CentOS-*

#添加EPEL仓库
sudo yum install epel-release
#更新yum源
sudo yum update

#Cent 8安装redis
yum install redis
#设置开机启动
systemctl enable redis.service
#启动
systemctl start redis

#systemctl start redis.service #启动redis服务
#systemctl stop redis.service #停止redis服务
#systemctl restart redis.service #重新启动服务
#systemctl status redis.service #查看服务当前状态
#systemctl enable redis.service #设置开机自启动
#systemctl disable redis.service #停止开机自启动

#编辑配置文件
vim /etc/redis.conf

#开启aof持久化, 开启可能会报 systemd错误
appendonly yes
## no：开机不会自启动；systemed：开机自启动
supervised systemd
#开启修改访问密码
requirepass password
#绑定主机
bind 0.0.0.0

#slave设置
# 设置为后台运行
daemonize yes
# 保存pid的文件，如果是在一台机器搭建主从，需要区分一下
pidfile /var/run/redis_6379.pid
# 指定日志文件
logfile redis.log

# slave中添加内容
#replicaof <masterip> <masterport>
replicaof 192.168.111.129 6379
#master中如果有密码,slave在设置相同
masterauth root2023

#同步延迟和同步速率
info replication
#master_repl_offset和slave_repl_offset
#主节点输出的数据的偏移量（offset）大于等于从节点的偏移量时，从节点就和主节点数据同步完成

#master_last_io_seconds_ago和slave_last_io_seconds_ago

#ubuntu apt 安装
#sudo apt install redis-server

# 堡塔命令行:
#bt

#mysql
# 删除老版本
#yum remove -y mysql
#find / -name mysql //找到残留的文件,再通过rm -rf去删除对应的文件

#禁用CentOS8自带mysql模块
#yum module disable mysql // 禁用命令

#yum install mysql-community-server 这一步的时候可能很多人安装不上，
#因为是yum安装库的问题，错误（Error: GPG check FAILED），可以将--nogpgcheck添加到后面：
#yum install mysql-community-server --nogpgcheck

#格式：mysql> set password for 用户名@localhost = password('新密码'); 
#mysql> set password for root@localhost = password('123');

