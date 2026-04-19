syscall.Unshare(syscall.CLONE_NEWUSER) 
Unshare CLONE_NEWUSER 失败: invalid argument

```bash
$ unshare --user --map-root-user echo "OK"
unshare: write failed /proc/self/uid_map: Operation not permitted
```

```bash
$ cat /etc/os-release
PRETTY_NAME="Ubuntu 24.04.4 LTS"
NAME="Ubuntu"
VERSION_ID="24.04"
VERSION="24.04.4 LTS (Noble Numbat)"
VERSION_CODENAME=noble
ID=ubuntu
ID_LIKE=debian
```

原因：Ubuntu24.04引入AppArmor保护机制，默认会限制非特权用户使用User Namespace

解决方法：

临时设置
```bash
$ sudo sysctl -w kernel.apparmor_restrict_unprivileged_userns=0
kernel.apparmor_restrict_unprivileged_userns = 0
```

```bash
$ unshare --user --map-root-user echo "OK"
OK
```

永久配置
```bash
$ echo "kernel.apparmor_restrict_unprivileged_userns=0" | sudo tee /etc/sysctl.d/99-userns.conf
kernel.apparmor_restrict_unprivileged_userns=0
```

```bash
$ sudo sysctl --system
* Applying /usr/lib/sysctl.d/10-apparmor.conf ...
* Applying /etc/sysctl.d/10-bufferbloat.conf ...
* Applying /etc/sysctl.d/10-console-messages.conf ...
* Applying /etc/sysctl.d/10-ipv6-privacy.conf ...
* Applying /etc/sysctl.d/10-kernel-hardening.conf ...
* Applying /etc/sysctl.d/10-magic-sysrq.conf ...
* Applying /etc/sysctl.d/10-map-count.conf ...
* Applying /etc/sysctl.d/10-network-security.conf ...
* Applying /etc/sysctl.d/10-ptrace.conf ...
* Applying /etc/sysctl.d/10-zeropage.conf ...
* Applying /usr/lib/sysctl.d/30-tracker.conf ...
* Applying /usr/lib/sysctl.d/50-bubblewrap.conf ...
* Applying /usr/lib/sysctl.d/50-pid-max.conf ...
* Applying /usr/lib/sysctl.d/99-protect-links.conf ...
* Applying /etc/sysctl.d/99-sysctl.conf ...
* Applying /etc/sysctl.d/99-userns.conf ...
* Applying /etc/sysctl.conf ...
kernel.apparmor_restrict_unprivileged_userns = 1
net.core.default_qdisc = fq_codel
kernel.printk = 4 4 1 7
net.ipv6.conf.all.use_tempaddr = 2
net.ipv6.conf.default.use_tempaddr = 2
kernel.kptr_restrict = 1
kernel.sysrq = 176
vm.max_map_count = 1048576
net.ipv4.conf.default.rp_filter = 2
net.ipv4.conf.all.rp_filter = 2
kernel.yama.ptrace_scope = 1
vm.mmap_min_addr = 65536
fs.inotify.max_user_watches = 65536
kernel.unprivileged_userns_clone = 1
kernel.pid_max = 4194304
fs.protected_fifos = 1
fs.protected_hardlinks = 1
fs.protected_regular = 2
fs.protected_symlinks = 1
kernel.apparmor_restrict_unprivileged_userns = 0
```
