```bash
# python prescript.py /var/sandbox/sandbox-python 3<json_print.py
allow syscall :262,16,8,217,1,3,257,0,202,9,12,10,11,15,25,105,106,102,39,110,186,60,231,234,13,16,24,273,274,334,228,96,35,291,233,230,270,201,14,131,318
Setgid============
{"name": "John", "age": 30, "city": "New York"}
```

```bash
strace -c python prescript.py /var/sandbox/sandbox-python 3<json_print.py
allow syscall :262,16,8,217,1,3,257,0,202,9,12,10,11,15,25,105,106,102,39,110,186,60,231,234,13,16,24,273,274,334,228,96,35,291,233,230,270,201,14,131,318
Setgid============
{"name": "John", "age": 30, "city": "New York"}
% time     seconds  usecs/call     calls    errors syscall
------ ----------- ----------- --------- --------- ------------------
 38.80    0.002421         242        10         1 futex
 10.14    0.000633           3       180        32 newfstatat
  8.41    0.000525         525         1           execve
  6.52    0.000407           5        68         4 openat
  6.30    0.000393           7        53           mmap
  6.06    0.000378           3        95           read
  3.59    0.000224           3        67           close
  3.45    0.000215           2       104           fstat
  2.72    0.000170           6        26           getdents64
  2.04    0.000127           7        17           brk
  1.75    0.000109           1        80         3 lseek
  1.71    0.000107           8        12           mprotect
  1.46    0.000091          91         1           seccomp
  0.93    0.000058           1        42        35 ioctl
  0.82    0.000051           4        12           tgkill
  0.75    0.000047           0       131           rt_sigaction
  0.74    0.000046           3        12           getpid
  0.58    0.000036           9         4           munmap
  0.53    0.000033           4         7         4 readlink
  0.40    0.000025           1        13           fcntl
  0.26    0.000016           3         5           epoll_ctl
  0.24    0.000015           3         4           write
  0.19    0.000012          12         1           chroot
  0.18    0.000011           1         6           rt_sigprocmask
  0.18    0.000011           5         2           pread64
  0.14    0.000009           9         1         1 access
  0.14    0.000009           9         1           pipe2
  0.10    0.000006           3         2           prlimit64
  0.08    0.000005           2         2           rt_sigreturn
  0.08    0.000005           2         2           sigaltstack
  0.08    0.000005           5         1           arch_prctl
  0.08    0.000005           5         1           set_robust_list
  0.08    0.000005           5         1           epoll_create1
  0.08    0.000005           5         1           rseq
  0.06    0.000004           4         1           setuid
  0.06    0.000004           4         1           setgid
  0.06    0.000004           4         1           set_tid_address
  0.06    0.000004           4         1           eventfd2
  0.05    0.000003           3         1           sched_getaffinity
  0.03    0.000002           1         2           chdir
  0.03    0.000002           2         1           prctl
  0.03    0.000002           1         2           gettid
  0.00    0.000000           0         2           getcwd
  0.00    0.000000           0         2           getrandom
  0.00    0.000000           0         1           clone3
------ ----------- ----------- --------- --------- ------------------
100.00    0.006240           6       980        80 total
```

# 运行 test.py 
1.复制需要运行的文件 json_print.py  在 /var/sandbox/sandbox-python 
2.在根目录执行 go run strace-syscall-map/main.go
