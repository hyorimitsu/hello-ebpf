#!/usr/bin/python

from bcc import BPF

text="""
#include <uapi/linux/ptrace.h>

int execve_hook(struct pt_regs *ctx) {
	int *p = (int*)0x1234;
	*p = 0xabcd;
	return 0;
}
"""

bpf = BPF(text=text)
bpf.attach_kprobe(event=bpf.get_syscall_fnname("execve"),
                  fn_name="execve_hook")

while True:
   bpf.trace_print()