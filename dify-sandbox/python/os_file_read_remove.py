import os
file_path = "python.so"

if os.path.exists(file_path):
    os.remove(file_path)
    print("File deleted")
else:
    print("File does not exist")