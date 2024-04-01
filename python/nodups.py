from email import parser
import sys
import argparse

# numbers of files
parser = argparse.ArgumentParser()
parser.add_argument('-n',type=int,help="number of files",nargs=1,required=True)
parser.add_argument('-d',type=str,help='input files with duplicates',nargs='?',action='append',required=True)
parser.add_argument('-o',type=str,help="file to save the output",nargs=1,default='nodups.txt')
args = parser.parse_args()

o = open(args.o[0],'w')

nodups= []
for i in range(args.n[0]):
    with open(args.d[i]) as file:
        lines = file.readlines()
        for line in lines:
            if line in nodups:
                continue
            else:
                nodups.append(line)
    file.close()

for i in nodups:
    if i == '\n':
        continue
    o.write(i)

o.close()