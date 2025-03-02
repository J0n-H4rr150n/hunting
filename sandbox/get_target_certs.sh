#!/bin/bash

output_file="cert_domains.txt"

# Clear the output file at the beginning (optional)
> "$output_file"

while IFS= read -r target; do
  # Handle IPv6 addresses correctly by enclosing in brackets for openssl
  if [[ "$target" =~ ":" ]]; then
    target_for_openssl="[$target]"
  else
    target_for_openssl="$target"
  fi

  # Use openssl s_client to connect and get the certificate
  certificate=$(openssl s_client -connect "${target_for_openssl}:443" -showcerts </dev/null 2>/dev/null | openssl x509 -noout -text)

  if [[ -n "$certificate" ]]; then # Check if a certificate was retrieved
    # Extract Subject Alternative Names (SANs)
    sans=$(echo "$certificate" | grep "Subject Alternative Name" -A 1 | grep -oE 'DNS:[^,]+')
    if [[ -n "$sans" ]]; then
      echo "$sans" | sed 's/DNS://g' | sed 's/,/\n/g' | sed 's/^  *//g' | while IFS= read -r domain; do
        echo "$domain" >> "$output_file" # Output SAN domains to file
      done
    else
      # Extract Common Name (CN) - if SANs are not present or for fallback
      cn=$(echo "$certificate" | grep "Subject:.*CN=" | sed 's/.*CN=//' | sed 's/,.*//')
      if [[ -n "$cn" ]]; then
        echo "$cn" >> "$output_file" # Output CN domain to file
      fi
    fi
  fi
done < targets.txt

echo "Domains extracted and saved to $output_file"
