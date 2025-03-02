#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 <proxy_list_file.txt>"
  exit 1
fi

proxy_file="$1"
check_url="https://api.ipify.org?format=json" # Or any IP checking API
timeout=5 # Timeout in seconds
alive_proxies_file="alive_proxies.txt"

# Remove the old alive_proxies.txt file
rm -f "$alive_proxies_file"

# Trap the SIGINT signal (Ctrl+C)
trap 'echo "Exiting..."; exit 0' SIGINT

while IFS= read -r proxy; do
  protocol=$(echo "$proxy" | cut -d'/' -f1 | tr -d ':')
  ip_port=$(echo "$proxy" | cut -d'/' -f3)
  ip=$(echo "$ip_port" | cut -d':' -f1)
  port=$(echo "$ip_port" | cut -d':' -f2)

  echo "Testing: $proxy"

  # Test Liveness (TCP Connection)
  if timeout $timeout nc -zv "$ip" "$port" >/dev/null 2>&1; then
    echo "  - Alive: Yes"

    # Add the proxy to alive_proxies.txt
    echo "$proxy" >> "$alive_proxies_file"

    # Test Anonymity (IP Check)
    if [[ "$protocol" == "http" || "$protocol" == "https" ]]; then
      proxy_env="--proxy $protocol://$ip:$port"
    elif [[ "$protocol" == "socks5" ]]; then
      proxy_env="--socks5 $ip:$port"
    else
        echo "unknown protocol"
        continue
    fi

    external_ip=$(curl -s $proxy_env $check_url | jq -r '.ip')

    if [ -n "$external_ip" ]; then
      local_ip=$(curl -s $check_url | jq -r '.ip')

      if [ "$external_ip" != "$local_ip" ]; then
        echo "  - Anonymous: Yes (External IP: $external_ip)"
      else
        echo "  - Anonymous: No (Your local IP is still visible)"
      fi
    else
      echo "  - Anonymity check failed."
    fi
  else
    echo "  - Alive: No"
  fi
done < "$proxy_file"
