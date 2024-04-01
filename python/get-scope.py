

import argparse
# to get the scope of the targer
parser = argparse.ArgumentParser()
parser.add_argument('-of',type=str,required=True, help='file contains out of scope words or domains')
parser.add_argument('-o',type=str,required=True,help='file to save the output')
parser.add_argument('-i',type=str,help='input file', required=True, help='input file')
args = parser.parse_args()


i = open(args.i, "r").readlines()
of = open(args.of,'r').readlines()
o = open(args.o,"w")

for line in i:
    for outofscopeword in of:
        s = line.find(outofscopeword)
        if s == -1:
            o.write(line)




o.close()




