```bash
$ cat /proc/sys/kernel/seccomp/actions_logged
kill_process kill_thread trap errno user_notif trace log
```
/proc/sys/kernel/seccomp/actions_logged 是 Linux 内核中用于控制 ‌Seccomp（安全计算模式）‌ 日志记录行为的一个可读写配置文件。

它的主要功能和特性如下：

1. 核心作用
该文件定义了一个‌允许被记录到系统日志中的 Seccomp 返回值列表‌。
当进程触发了 Seccomp 过滤器规则时，内核会根据此文件中配置的返回值类型，决定是否将相关事件记录到内核日志（如 dmesg 或 syslog）中。这主要用于调试、监控和安全审计。

2. 文件格式与排序
内容格式‌：文件中保存的是 Seccomp 返回值宏（如 SECCOMP_RET_KILL、SECCOMP_RET_ERRNO 等）对应的字符串名称。
读取顺序‌：从该文件读取内容时，返回的列表是‌有序‌的。排序规则与 /proc/sys/kernel/seccomp/actions_avail 一致，即按照“最少许可”到“最多许可”的顺序排列（例如，阻止系统调用的动作排在允许系统调用的动作之前）。
写入顺序‌：向该文件写入数据时，不需要严格遵循特定顺序，内核会自动处理。
3. 重要限制
不支持记录 "allow"‌：你‌不能‌将 allow 字符串写入此文件。
原因‌：SECCOMP_RET_ALLOW 表示允许系统调用正常执行，记录所有允许的系统调用会产生巨大的性能开销和日志噪音，因此内核明确禁止记录此类动作。
后果‌：如果尝试向该文件写入 allow，操作会失败并返回 EINVAL（无效参数）错误。
4. 常见使用场景
调试 Seccomp 策略‌：当应用程序因 Seccomp 规则被阻断时，管理员可以通过检查此配置和内核日志，确认哪些系统调用被拦截以及拦截的原因。
安全监控‌：在生产环境中，可以配置记录特定的高风险动作（如 kill 或 trap），以便及时发现潜在的攻击行为或程序错误。


```bash
# grep "sig=31" /var/log/syslog
```

```bash
$ grep "sig=31" /var/log/audit/audit.log
```
