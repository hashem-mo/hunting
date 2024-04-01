import sys,argparse,hashlib

parser = argparse.ArgumentParser()
parser.add_argument('-t',type=str,help='the real text NOT HASHED')
parser.add_argument('-hash',type=str,help="the hash that is wanted to be indentified")
 
args = parser.parse_args()
text = args.t.encode('utf-8')
hash = args.hash 

algs = ['sha3_256', 'sha3_512', 'shake_128', 'shake_256', 'md5', 'sha512', 'sha384', 'sha3_384', 'sha3_224', 'blake2b', 'blake2s', 'sha256', 'sha1', 'sha224']

for alg in algs:
    hashtest = getattr(hashlib,alg)
    hashobject = hashtest()
    hashobject.update(text)
    hex = hashobject.hexdigest()
    if hash == hex:
        print(f"Hash indentified!\nAlgorithm: {alg}")
        sys.exit(0)
    else:
        continue
print("sorry! your hash doesn't match any known algorithms")
