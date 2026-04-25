import os
import shutil

def clear_current_directory():
    current_dir = os.getcwd()
    print(current_dir)
    
    for filename in os.listdir(current_dir):
        file_path = os.path.join(current_dir, filename)
        try:
            if os.path.isfile(file_path) or os.path.islink(file_path):
                os.unlink(file_path)
            elif os.path.isdir(file_path):
                shutil.rmtree(file_path)
        except Exception as e:
            print(f'delete file error {file_path}: {e}')

if __name__ == "__main__":
    clear_current_directory()