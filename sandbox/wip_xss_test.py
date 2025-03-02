#!/usr/bin/env python3

import requests
import sys
import urllib.parse
import io  # For parsing request file
import warnings # Import the warnings module
import urllib3 # Import urllib3 to reference the warning class

# --- Configuration ---
FUZZ_KEYWORD = "FUZZ"           # Keyword in the request file to replace with payloads
MITMPROXY_ADDRESS = "http://localhost:8080"  # Address of your mitmproxy
DEFAULT_SCHEME = "https://"     # Default scheme to use if not explicitly in URL

# Suppress the InsecureRequestWarning
warnings.filterwarnings('ignore', category=urllib3.exceptions.InsecureRequestWarning)


def parse_request_file(request_file_path):
    """Parses the request file to extract method, url_path, headers, and body,
       and constructs the full URL."""
    method = None
    url_path = None # Now storing path, not full URL initially
    headers = {}
    body = ""
    is_body = False
    host_header_value = None # To store the Host header value

    with open(request_file_path, 'r') as f:
        request_lines = f.readlines()

    if not request_lines:
        raise ValueError("Request file is empty.")

    # First line: Request method and URL Path
    request_line_parts = request_lines[0].strip().split()
    if len(request_line_parts) >= 2:
        method = request_line_parts[0].upper()
        url_path = request_line_parts[1] # Extracting path only
    else:
        raise ValueError("Invalid request line in request file.")

    # Headers and body
    for line in request_lines[1:]:
        line = line.strip()
        if line == "":
            is_body = True  # Empty line indicates start of body
            continue
        if not is_body:
            header_parts = line.split(":", 1)
            if len(header_parts) == 2:
                header_name = header_parts[0].strip()
                header_value = header_parts[1].strip()
                headers[header_name] = header_value
                if header_name.lower() == "host": # Capture Host header value
                    host_header_value = header_value
        else:
            body += line + "\n" # Reconstruct body, preserving newlines

    if not host_header_value:
        raise ValueError("Host header is missing in the request file.")

    # Construct the full URL using scheme, Host header, and URL path
    full_url = DEFAULT_SCHEME + host_header_value + url_path

    return method, full_url, headers, body.rstrip("\n") # Return full URL


def main():
    # --- Check for correct number of arguments ---
    if len(sys.argv) != 3:
        print("Usage: python xss_injector.py <request_file> <payload_file>")
        print("  <request_file>: Path to the request template file (e.g., basic.req)")
        print("  <payload_file>: Path to the XSS payload list file (e.g., xss_basic.txt)")
        sys.exit(1)

    request_file = sys.argv[1]
    payload_file = sys.argv[2]

    # --- Check if files exist and parse request ---
    try:
        method, base_url, headers, request_body = parse_request_file(request_file) # parse_request_file now returns full URL
    except FileNotFoundError:
        print(f"Error: Request file '{request_file}' not found.")
        sys.exit(1)
    except ValueError as e:
        print(f"Error parsing request file '{request_file}': {e}")
        sys.exit(1)


    try:
        with open(payload_file, 'r') as f:
            payloads = [line.strip() for line in f] # Read payloads, removing whitespace
    except FileNotFoundError:
        print(f"Error: Payload file '{payload_file}' not found.")
        sys.exit(1)

    # --- Loop through payloads and send requests through mitmproxy ---
    for payload in payloads:
        print(f"--- Testing payload: {payload} ---")

        # URL-encode the payload
        url_encoded_payload = urllib.parse.quote_plus(payload)

        # Replace FUZZ keyword with the URL-encoded payload
        modified_url = base_url.replace(FUZZ_KEYWORD, url_encoded_payload) # Use full URL now
        modified_body = request_body.replace(FUZZ_KEYWORD, url_encoded_payload)


        proxies = {
            "http": MITMPROXY_ADDRESS,
            "https": MITMPROXY_ADDRESS, # If you need to test HTTPS as well
        }

        print("REQUEST:") # Request Header
        print("") # Line break after header
        print(f"Method: {method}")
        print(f"URL: {modified_url}") # Print full constructed URL
        print(f"Headers: {headers}")
        if modified_body:
            print(f"Body: {modified_body}")


        try:
            response = requests.request(
                method=method,
                url=modified_url, # Use full constructed URL
                headers=headers,
                data=modified_body, # Use data for body in requests
                proxies=proxies,
                verify=False # Disable SSL verification for proxy - be cautious in real scenarios
            )

            print("") # Line break before response header
            print("RESPONSE:") # Response Header
            print("") # Line break after header
            print(f"Response Status Code: {response.status_code}")
            # Optionally print response content for debugging:
            # print("Response Content:")
            # print(response.text)


        except requests.exceptions.RequestException as e:
            print(f"Request Error: {e}")

        print("-----------------------------\n") # Separator line after each response


    print("--- XSS testing script finished ---")


if __name__ == "__main__":
    main()
