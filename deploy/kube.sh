#CentOS需要添加源
sudo sed -i -e "s|mirrorlist=|#mirrorlist=|g" /etc/yum.repos.d/CentOS-* && sudo sed -i -e "s|#baseurl=http://mirror.centos.org|baseurl=http://vault.centos.org|g" /etc/yum.repos.d/CentOS-*

# 安装docker所需的工具
yum install -y yum-utils device-mapper-persistent-data lvm2
# 配置阿里云的docker源
yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
#安装docker
curl -s https://get.docker.com/ | sh
# 启动docker
systemctl enable docker && systemctl start docker
#关闭防火墙
systemctl stop firewalld && systemctl disable firewalld
# 关闭selinux
# 临时禁用selinux
# setenforce 0
# 永久关闭 修改/etc/sysconfig/selinux文件设置
setenforce 0 && sed -i 's/SELINUX=permissive/SELINUX=disabled/' /etc/sysconfig/selinux && sed -i "s/SELINUX=enforcing/SELINUX=disabled/g" /etc/selinux/config
# 禁用交换分区
# swapoff -a
# 永久禁用，打开/etc/fstab注释掉swap那一行。
swapoff -a && sed -i 's/.*swap.*/#&/' /etc/fstab
# 修改内核参数
cat <<EOF >  /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF
sysctl --system
# 执行配置k8s阿里云源
cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64/
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg https://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
EOF
#安装kubelet、kubeadm
yum install -y kubelet-1.23.0 kubeadm-1.23.0
# 启动kubelet服务
systemctl enable kubelet && systemctl start kubelet

#修改docker配置 
#docker驱动查看
#docker info | grep "Storage Driver"
#停止docker
#systemctl stop docker
#更改docker驱动前先备份
#sudo cp -rf /var/lib/docker /var/lib/docker.backup
#修改存储驱动
# sudo vim /etc/docker/daemon.json
# {
#   "storage-driver": "<new_storage_driver>"
# }
# systemctl daemon-reload && systemctl restart docker
#更改docker驱动与kubeadm统一
cat <<EOF > /etc/docker/daemon.json
{
     "exec-opts": [
         "native.cgroupdriver=systemd"
     ]
}
EOF
systemctl daemon-reload && systemctl restart docker
#更改驱动后清理旧驱动
#sudo rm -rf /var/lib/docker
#sudo mv /var/lib/docker.backup /var/lib/docker

#修改kubeadm配置
# vim /etc/sysconfig/kubelet
#修改kub启动驱动
#打开后添加参数 KUBELET_EXTRA_ARGS=--cgroup-driver=cgroupfs
#解决initial timeout 40s(1.23不支持)
#打开后添加参数 KUBELET_EXTRA_ARGS=--feature-gates=SupportPodPidsLimit=false,SupportNodePidsLimit=false
# systemctl daemon-reload && systemctl restart kubelet

cat <<EOF > /etc/hosts
192.168.111.129 master.xhhy.com kube_master
192.168.111.130 node1.xhhy.com kube_node1
192.168.111.131 node2.xhhy.com kube_node2
EOF
#使host生效
systemctl restart NetworkManager



#如果有问题,重启一下
#reboot



#以下master主机需要安装
#安装kubectl-master
yum install -y kubectl-1.23.0
# 下载管理节点中用到的6个docker镜像，你可以使用docker images查看到
# 重要：这里的--apiserver-advertise-address使用的是master和node间能互相ping通的ip
# 这里需要大概两分钟等待，会卡在[preflight] You can also perform this action in beforehand using ''kubeadm config images pull
kubeadm init --apiserver-advertise-address=192.168.111.129 --image-repository=registry.aliyuncs.com/google_containers --kubernetes-version v1.23.0 --service-cidr=10.10.10.0/24 --pod-network-cidr=10.20.20.0/24 --token-ttl 0 --ignore-preflight-errors=all
#systemctl status kubelet
#journalctl -xeu kubelet
#重置一下,再次进行init
#kubeadm reset
#rm -rf /etc/cni/net.d



#拷贝k8s认证文件
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config && sudo chown $(id -u):$(id -g) $HOME/.kube/config

#master节点创建token,kubeadm init时会生成token
#手动创建可以加入集群的token
#kubeadm token create --print-join-command
#其它node节点加入K8S集群(token有效期24小时)命令
# 例如:
#kubeadm join 192.168.111.129:6443 --token zfuvyw.f2j9i81atjcmrm6o --discovery-token-ca-cert-hash sha256:5826b43bcdcd6943bc281e3d134ea0844d40bffed0856a41137a5c79ea6828e3 


# Calico
wget https://docs.projectcalico.org/v3.19/manifests/calico.yaml --no-check-certificate
# wget https://docs.projectcalico.org/manifests/calico.yaml --no-check-certificate
#修改为国内镜像加速站
cat calico.yaml |grep 'image:'
sed -i 's#docker.io/##g' calico.yaml
#修改calico.yaml文件中CALICO_IPV4POOL_CIDR的value与--pod-network-cidr相同
#例如:--pod-network-cidr=10.20.20.0/24
#部署
kubectl apply -f calico.yaml
#删除
#kubectl delete -f calico.yaml
#查看状态，执行完上一条命令需要等一会才全部running
kubectl get pods -n kube-system


#设置时区
# timedatectl set-timezone Asia/Shanghai

# vim /etc/chrony.conf
# # pool 2.centos.pool.ntp.org iburst
# server 210.72.145.44 iburst
# server ntp.aliyun.com iburst
# systemctl restart chronyd.service && systemctl enable chronyd.service
# chronyc sources -v && date


#登陆Token
#eyJhbGciOiJSUzI1NiIsImtpZCI6IkREaW1TVWNTVU9UWldtaDNuZkhpc1RqS09hZEo2NFhVa0JEVEZqa0I3S1kifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJkYXNoYm9hcmQtYWRtaW4tdG9rZW4tZnJ0azYiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGFzaGJvYXJkLWFkbWluIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQudWlkIjoiMzY0NzNmNWMtZDk4Ny00MDk4LTg2YzItNWEyNmM5MThiYjY4Iiwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50Omt1YmUtc3lzdGVtOmRhc2hib2FyZC1hZG1pbiJ9.IASOTQCa16-cJRh7Tn7V8ZFyiDomT4jT2_GLel0Xe0pCJaivJUQFEKkrRI7P8uQqKviI5HlqKJvJCltjmsXkqe7M8Mx7AhsCqv-QO0xtUgJ4amcWu644C-IzMF9bKH6Gvu_s38NMO0DTkZdQkTidA8KqQfkCgomgOTzkgy7BphX2pk5ibjtIjo-c-ViOwfxCQB6-7-7kRlzsZ8OfsPvxQ9W526NufZchbdGzhCQKLpkhKJh_X1rK3l46ErLwwN3CLw2L459RhSuQLumI9g2oju8SG_7wDu3rwKo2T3eTfeKnTN8AplmbMhgTAHMqWpTGEhs6t-o8HrdcD6QL5kx-ZQ
#部署Dashboard
#报错就多试几下
wget https://raw.githubusercontent.com/kubernetes/dashboard/v2.6.1/aio/deploy/recommended.yaml
#修改39行 nodePort范围 30000-32767
#spec
#   ports:
#       ...
#增加   nodePort: 31000
#增加   type:NodePort
kubectl apply -f recommended.yaml
#kubectl delete -f recommended.yaml

#在master节点创建service account并绑定默认cluster-admin管理员集群角色
# 创建用户
kubectl create serviceaccount dashboard-admin -n kube-system
# 用户授权
kubectl create clusterrolebinding dashboard-admin --clusterrole=cluster-admin --serviceaccount=kube-system:dashboard-admin
# 获取用户Token
kubectl describe secrets -n kube-system $(kubectl -n kube-system get secret | awk '/dashboard-admin/{print $1}')

#访问
# https://(masterip/nodeip):31000



#修改kube名字
#修改系统主机名
hostname kube_master