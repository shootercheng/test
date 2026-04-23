import ctypes
import os
import sys
import traceback


# setup sys.excepthook
def excepthook(type, value, tb):
    sys.stderr.write("".join(traceback.format_exception(type, value, tb)))
    sys.stderr.flush()
    sys.exit(-1)


sys.excepthook = excepthook

lib = ctypes.CDLL("/home/scd/code/go/dify/dify-sandbox/internal/core/runner/python/python.so")
lib.DifySeccomp.argtypes = [ctypes.c_uint32, ctypes.c_uint32, ctypes.c_bool]
lib.DifySeccomp.restype = None

# get running path
running_path = sys.argv[1]
if not running_path:
    exit(-1)

os.chdir(running_path)

lib.DifySeccomp(10099, 1001, 1)

with open('json_print.py', 'r', encoding='utf-8') as f:
    code = f.read()
    exec(compile(code, 'json_print.py', 'exec'))
