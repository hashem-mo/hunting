import sys
# for removing the "https" or "http" from httpx tool for the tools that needs subdomains nor urls
# takes  first arg a file contains urls and second arg for output


f = open(sys.argv[1],'r').readlines()
out = open(sys.argv[2], 'w')
for i in f:
    if i[5] == 's':
        out.write(i[9:])
    else:
        out.write(i[8:])

out.close()
