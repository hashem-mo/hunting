

import urllib.parse
import re
import os,requests
import urllib3
urllib3.disable_warnings()
from concurrent.futures import ProcessPoolExecutor
import argparse


parser = argparse.ArgumentParser()
parser.add_argument('-i',type=str,required=True,help="Input file contains URLs")
parser.add_argument('-q',type=str,required=True,help="A string to test reflection")
parser.add_argument("-f",type=str,help="A folder to save the responses in")
parser.add_argument('-t',type=int,default=3,help="Numbers of threads to run")
parser.add_argument("-c",type=bool,default=False,help="This option is used to when the params query include Non-alphanumeric characters to check if they are html to reduce false positives")
parser.add_argument("-o",type=str,help="file name to save the vulnerable endpoints")
parser.add_argument("-p",type=str,help="http proxy to use")
args = parser.parse_args()

if args.f != None:
    os.mkdir(args.f)
    folder = True
if args.o != None: 
    outputfile = open(args.o,'w')
    file = True

os.system('color')
colors = {'magenta': '\x1b[95m', 'blue': '\x1b[94m', 'cyan': '\x1b[96m', 'green': '\x1b[92m', 'yellow': '\x1b[93m', 'red': '\x1b[91m', 'ENDC': '\x1b[0m', 'BOLD': '\x1b[1m', 'UNDERLINE': '\x1b[4m'}


proxies = {'https':args.p,
'http':args.p}
sent = []
headers = {"User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36","Accept-Encoding":"gzip, deflate, br","Accept-Language":"en-US,en;q=0.9"}
def requester(url):

    url = urllib.parse.unquote(url)
    url = url.strip()

    if '=' not in url:
        return

    info = urllib.parse.urlparse(url)
    
    global sent
    if info.netloc in sent:
        return
    sent.append(info.netloc)

    u = "https://" + info.netloc + info.path
    if url[:5] == "http:":
        u = "http://" + info.netloc + info.path 
    param = dict(item.split('=') for item in info.query.split('&'))
    
    try:
        if args.p != None:
            r = requests.get(u,params=param,allow_redirects=False,timeout=10,proxies=proxies,verify=False)   
        else:
            r = requests.get(u,params=param,allow_redirects=False,timeout=10,verify=False)           
        if r.status_code == 200:
            if re.findall(args.q,r.text) != []:

                print('[' + colors['red'] + "Critical" + colors['ENDC'] +']' + colors['BOLD'] + " Found reflection in:\t" + colors['ENDC'] + url)


                if folder == True:
                    f = open(args.f + "/" + url +'.txt','w')
                    for line in r.text:
                        f.write(line)

                if file == True:
                    outputfile.write(url)
                
    except:
        return




def main():

        
    urls = open(args.i,'r').readlines()
    with ProcessPoolExecutor(max_workers=args.t) as executor:
        futures = [executor.submit(requester,url) for url in urls]
    







if __name__ == "__main__":
    main()