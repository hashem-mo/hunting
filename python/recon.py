from xml import dom
import requests,sys,os,threading,argparse
from time import sleep

# def sublister(input,inputType,silent):
    
#     i = 0

#     if silent == False:
#         print('Starting Sublist3r....')
#     l = []
#     allSubs = []
#     if inputType == 'domain':
#         command = os.popen(f'python3 /../../../../../../../hunting/tools/sublist3r/sublist3r.py -d {input} -t 5')
#         output = list(command.read())
#     else:
#         f = open(input,'r')
#         l = f.readlines()
#         f.close()
#         for domain in l:
#             command = os.popen(f'python3 /../../../../../../../hunting/tools/sublist3r/sublist3r.py -d {domain} -t 5')
#             i = i + 1 
#             output = list(command.read())
#             for sub in output:
#                 allSubs.append(sub)


            
#         # for o in range(i):
#         #     output = open(f'sublister{o}.txt','r').readlines()
#         #     for sub in output:
#         #         list.append(sub)
#         # f = open('all_sublister','w')

#         for i in allSubs:
#             f.write(i)
#     return


def amass(input,inputType,silent):
    if silent == False:
        print('Starting Amass....')
    
    if inputType == 'domain' :
        os.system(f'amass enum -d {input} -passive -norecursive -noalts -silent > amass.txt')
        # output = list(command.read())
        # file = open('amass.txt','w')
        # for i in output:
        #     file.write(i)
    else:
        amass = open("amass.txt",'w')
        files_names = []
        f = open(input,'r').readlines()
        for domain in f:
            os.system(f'amass enum -d {input} -passive -norecursive -noalts -silent > {domain}.txt')
            files_names.append(domain)
        print(files_names)
        for file in files_names:
            subdomains = open(file.rstrip(),'r')
            for i in subdomains.readlines():
                amass.write(i)
            subdomains.close()

        for name in files_names:
            os.system(f"rm -rf {name}")



        # output = list(command.read())
        # file = open('amass.txt','w')
        # for i in output:
        #     file.write(i)

    return


def subfinder(input,inputType,silent):

    if silent == False:
        print('Starting Subfinder....')    
    print(input)
    if inputType == 'domain' :
        os.system(f'subfinder -d {input} -all -silent > subfinder.txt')
        # output = list(command.readlines())
        # file = open('subfinder.txt','w')
        # for i in output:
        #     file.write(i)
    else:
        os.system(f'subfinder -dL {input} -all -silent > subfinder.txt')   
        # output = list(command.read())
        # file = open('subfinder.txt','w')
        # for i in output:
        #     file.write(i)

    return    



def subdomain_enumeration(input,inputType,silent):

    global t1,t2,t3

    t1 = threading.Thread(target=subfinder,args=(input,inputType,silent),daemon=True)
    t2 = threading.Thread(target=amass,args=(input,inputType,silent),daemon=True)
    # t3 = threading.Thread(target=sublister,args=(input,inputType,silent),daemon=True)

    t1.start()
    t2.start()
    # t3.start()    


    i = 0
    while True:
        if not t1.is_alive():
            i = i + 1
            if silent == False:
                print('Subfinder Done! \nFILE : subfinder.txt')

        if not t2.is_alive():
            i = i + 1
            if silent == False:
                print('Amass Done! \nFILE : amass.txt')

        # if not t3.is_alive():
        #     i = i + 1
        #     if silent == False:
        #         print('Sublist3r Done! \nFILE : sublister.txt')

        if i == 1:
            break
    t1.join()
    t2.join()
    
    sleep(3)

    subfinder_f = open('subfinder.txt','r')
    amass_f = open('amass.txt','r')
    # sublister_f = open('sublister.txt','r')

    subfinder_list = subfinder_f.readlines()
    amass_list = amass_f.readlines()
    # sublister_list = sublister_f.readlines()
    
    subfinder_f.close()
    amass_f.close()
    # sublister_f.close()

    all_subdomains = []
    all_subdomains.extend(subfinder_list)
    all_subdomains.extend(amass_list)
    # all_subdomains.extend(sublister_list)

    all_subdomains = list(dict.fromkeys(all_subdomains))

    subdomains = open('subdomains.txt','w')

    for sub in all_subdomains:
        subdomains.write(sub)
    
    return all_subdomains


def httpProbing():
    pass


def main():
    
    parser = argparse.ArgumentParser()
    parser.add_argument('-d','--domain',type=str,help='Single domain to do recon on')
    parser.add_argument('-l','--file',type=str,help='File contains list of domains/subdomains to use')
    parser.add_argument('-o','--output',type=str,help='File to save the output')
    parser.add_argument('-s', '--silent',action='store_true')
    args = parser.parse_args()
    
    if args.domain == None and args.file == None:
        sys.exit("YOU MUST PROVIDE AN INPUT")

    if args.file == None:
        subdomains = subdomain_enumeration(args.domain,'domain',args.silent)        
    elif args.domain == None:
        if os.path.exists(args.file) == True:
            subdomains = subdomain_enumeration(args.file,'file',args.silent)
        else:
            sys.exit('FILE DOES NOT EXIST!')
    else:
        sys.exit('ERROR!')

    #urls = httpProbing(subdomains)


if __name__=='__main__':
    main()