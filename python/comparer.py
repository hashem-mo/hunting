import sys

def compare_files(file1, file2):
    with open(file1, 'r') as f1:
        urls1 = set(f1.read().splitlines())
    with open(file2, 'r') as f2:
        urls2 = set(f2.read().splitlines())
    dead = urls1 - urls2
    print("DOWN HOSTS")
    for url in dead:
        print(url)
    alive = urls2 - urls1
    print("NEW UP HOSTS")
    for i in alive:
        print(i)

if __name__ == '__main__':
    file1 = sys.argv[1]
    file2 = sys.argv[2]
    compare_files(file1, file2)