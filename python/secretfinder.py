import requests,sys,argparse,re
import concurrent.futures
import validators,os

# Parsing arguments
parser = argparse.ArgumentParser()
parser.add_argument('-f','--file',type=str,help="The file that contains URLs to be examined")
parser.add_argument('-u','--url',type=str,help="A single URL to be examined")
parser.add_argument('-o','--output',type=str,help="A file to save the output")
parser.add_argument('-t','--threads',type=int,help='The number of concurrent threads to use',default=1)
args = parser.parse_args()

os.system('color')
colors = {'magenta': '\x1b[95m', 'blue': '\x1b[94m', 'cyan': '\x1b[96m', 'green': '\x1b[92m', 'yellow': '\x1b[93m', 'red': '\x1b[91m', 'ENDC': '\x1b[0m', 'BOLD': '\x1b[1m', 'UNDERLINE': '\x1b[4m'}


# Define the list of words to search for
words_to_search = open("/hunting/programs/linktree/recon/secretswordlist.txt",'r').readlines()


# Define a function to send HTTP requests and search for words in the response
def send_request(url):
    
        # Remove any whitespace characters from the URL
        url = url.strip()
        # Send a GET request to the URL
        response = requests.get(url,allow_redirects=False,timeout=10)
        response_post= requests.post(url,allow_redirects=False,timeout=10)
        # Check if the response was successful
        if response.status_code == 200:
            # Iterate over each word in the list
            for word in words_to_search:
                match = re.findall(word, response.text)
                match_post = re.findall(word, response_post.text)
                # Check if the word is in the response text
                if match != []:
                    print('[' + colors['red'] + "Critical" + colors['ENDC'] +']' + f"{word}" + colors['BOLD'] + "FOUND IN" + colors["ENDC"] + f"{url}\t\tMETHOD:"+ colors['yellow'] + "GET" + colors['ENDC'])
        if response_post.status_code == 200:
            # Iterate over each word in the list
            for word in words_to_search:
                # Check if the word is in the response text
                if match_post != []:
                    print('[' + colors['red'] + "Critical" + colors['ENDC'] +']' + f"{word}" + colors['BOLD'] + "FOUND IN" + colors["ENDC"] + f"{url}\t\tMETHOD:" + colors['yellow'] + "POST" + colors['ENDC'])
        else:
            pass



def main():
    if args.file:
        file = open(args.file,'r')
        urllist = file.readlines()
        file.close()
    elif args.url:
        url = args.url
        if not validators.url(url):
            print("Please provide a valid URL\n\n Exitting....")
            exit()
        send_request(url)
        exit()

    else:
        print("You must provide an input whether a URL or a file contains URLs a one by line! \n\n Exitting....")
        exit()
    # Create a thread pool executor with the specified number of threads
    
    with concurrent.futures.ThreadPoolExecutor(max_workers=args.threads) as executor:
        # Submit a task for each URL in the file
        futures = [executor.submit(send_request, url) for url in urllist]
        # Wait for all tasks to complete
        concurrent.futures.wait(futures)
    





if __name__ == '__main__':
    main()