import json
#import argparse
import sys
#parser = argparse.ArgumentParser()
#parser.add_argument('-f',type=str,required=True)
#parser.add_argument('-o',type=str,required=False)



f = open(sys.argv[1],'r',encoding='utf-8')
o = open(sys.argv[2],'w')

data = [json.loads(line) for line in f]
for i in data:
    if i['failed'] == False:
        print(i['url'])
        if sys.argv[2] != None:
            o.write(i["url"]+'\n')