#!/bin/bash

# Check if the input file is provided
if [ -z "$1" ]; then
  echo "Usage: $0 <ip_list_file>"
  exit 1
fi

input_file="$1"

# Check if the input file exists
if [ ! -f "$input_file" ]; then
  echo "Error: Input file '$input_file' not found."
  exit 1
fi

# Read each IP from the input file
while IFS= read -r ip; do
  # Sanitize the IP for use in a filename (replace dots with underscores)
  sanitized_ip=$(echo "$ip" | tr '.' '_')

  # Create the output filename
  output_file="${sanitized_ip}.txt"

  # Run nmap and redirect the output to the file
  nmap "$ip" > "$output_file" 2>&1 #redirect standard output and standard error

  # Check if nmap was successful (optional)
  if [ $? -eq 0 ]; then
    echo "Nmap scan for $ip saved to $output_file"
  else
    echo "Nmap scan for $ip failed. Errors saved to $output_file"
  fi
done < "$input_file"

echo "All scans completed."

exit 0
