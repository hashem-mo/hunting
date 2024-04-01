import sys

f = open(sys.argv[1],'r').readlines()
o = open('tmp.txt','w')
i = list(dict.fromkeys(f))
for l in i :
    o.write(l)