import requests
import sys
import argparse
import concurrent.futures
import validators
import threading

# Parsing arguments
parser = argparse.ArgumentParser()
parser.add_argument('-f','--file',type=str,help="The file that contains URLs to be examined")
parser.add_argument('-u','--url',type=str,help="A single URL to be examined")
parser.add_argument('-o','--output',type=str,help="A file to save the output")
parser.add_argument('-t','--threads',type=int,help='The number of concurrent threads to use',default=1)
args = parser.parse_args()

# Define the URLs file path
urls_file_path = sys.argv[1]

# Define the list of words to search for
words_to_search = open("/hunting/programs/linktree/recon/secretswordlist.txt",'r').readlines()

stop_event = threading.Event()

# Define a function to send HTTP requests and search for words in the response
def send_request(url):
    # Check if the stop event is set before performing the request
    if not stop_event.is_set():
        # Remove any whitespace characters from the URL
        url = url.strip()

        # Send a GET request to the URL
        response = requests.get(url,allow_redirects=False,timeout=10)
        response_post= requests.post(url,allow_redirects=False,timeout=10)
        # Check if the response was successful
        if response.status_code == 200:
            # Iterate over each word in the list
            for word in words_to_search:
                # Check if the word is in the response text
                if word in response.text:
                    print(f"{word} FOUND IN {url} METHOD GET")
        if response_post.status_code == 200:
            # Iterate over each word in the list
            for word in words_to_search:
                # Check if the word is in the response text
                if word in response.text:
                    print(f"{word} FOUND IN {url} METHOD POST")
        else:
            pass

def main():
    try:
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
    except KeyboardInterrupt:
        print("Execution stopped by user.")
        stop_event.set()
        executor.shutdown(wait=True)

if __name__ == '__main__':
    main()