#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 <alive_proxies_file.txt>"
  exit 1
fi

alive_proxies_file="$1"
output_file="list_proxychains4.conf"

# Start with the header
echo "dynamic_chain" > "$output_file"
echo "[ProxyList]" >> "$output_file"

# Read proxies from the alive_proxies.txt file
while IFS= read -r proxy; do
  # Extract protocol, IP, and port (adjust if your format is different)
  protocol=$(echo "$proxy" | cut -d ':' -f 1)
  ip_port=$(echo "$proxy" | cut -d ':' -f 2-)
  ip=$(echo "$ip_port" | cut -d ':' -f 1)
  port=$(echo "$ip_port" | cut -d ':' -f 2)

  # Remove "//" from the IP if present
  ip=$(echo "$ip" | sed 's/\/\///g')

  # Add the proxy to the output file with proper formatting
  echo "$protocol $ip $port" >> "$output_file"
done < "$alive_proxies_file"

echo "Proxy list updated in $output_file"
