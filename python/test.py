# import sys,os
# sys.path.insert(1, '/../../../../../../../hunting/tools/sublist3r')

# import sublist3r



# def sublister(input,inputType,silent):
    
#     if silent == False:
#         print('Starting Sublist3r....')

#     if inputType == 'domain':
#         subdomains = sublist3r.main(input, 10, 'sublister.txt', ports= None, silent=False, verbose= False, enable_bruteforce= False, engines=None)
#     else:
#         f = open(input,'r')
#         list = f.readlines()
#         f.close()

#         for link in list:
#             subdomains = sublist3r.main(link, 10, 'sublister.txt', ports= None, silent=True, verbose= False, enable_bruteforce= False, engines=None)
    
#     if silent == False:
#         print('Sublist3r Done! \nFILE : sublister.txt')

#     return


# sublister("google.com",'domain',False)



















#     i = 0
#     while True:
#         if t1.is_alive() == True:
#             pass
#         else:
#             i = i + 1
#             if silent == False:
#                 print('Subfinder Done! \nFILE : subfinder.txt')
#         if t2.is_alive() == True:
#             pass
#         else:
#             i = i + 1
#             if silent == False:
#                 print('Amass Done! \nFILE : amass.txt')
#         if t3.is_alive() == True:
#             pass
#         else:
#             i = i + 1
#             if silent == False:
#                 print('Sublist3r Done! \nFILE : sublister.txt')
#         if i == 3:
#             break





import os
c = os.popen('ls')

s = c.read()
