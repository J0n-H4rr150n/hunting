#!/bin/bash

output_file="certificate_validity.json"
input_file="${1:-targets.txt}"

if [ ! -f "$input_file" ]; then
  echo "Error: Input file '$input_file' not found."
  echo "Usage: $0 <input_file.txt>"
  exit 1
fi

> "$output_file"

json_output_array="["

process_count=0

while IFS= read -r target; do
  if [[ $process_count -gt 0 ]]; then
    json_output_array+=", "
  fi
  process_count=$((process_count + 1))

  echo "Processing target: $target"

  if [[ "$target" =~ ":" ]]; then
    target_for_openssl="[$target]"
  else
    target_for_openssl="$target"
  fi

  certificate_output=$(openssl s_client -connect "${target_for_openssl}:443" -showcerts </dev/null 2>/dev/null)

  if echo "$certificate_output" | grep -q "BEGIN CERTIFICATE"; then
    certificate_raw=$(echo "$certificate_output")

    not_before_raw=$(echo "$certificate_raw" | openssl x509 -noout -startdate | sed 's/^notBefore=//')
    not_before=$(date -d "$not_before_raw" +%Y-%m-%d)

    not_after_raw=$(echo "$certificate_raw" | openssl x509 -noout -enddate | sed 's/^notAfter=//')
    not_after=$(date -d "$not_after_raw" +%Y-%m-%d)

    echo "Certificate Validity for $target:"
    echo "  Not Before: $not_before"
    echo "  Not After:  $not_after"

    json_output_array+="{ \"target\": \"$target\", \"NotBefore\": \"$not_before\", \"NotAfter\": \"$not_after\" }"

  else
    echo "No certificate retrieved from $target."
    json_output_array+="{ \"target\": \"$target\", \"status\": \"No Certificate Found\" }"
  fi
  echo "---"
done < "$input_file"

json_output_array+="]"

if command -v jq &> /dev/null; then
  echo "$json_output_array" | jq '.' > "$output_file"
  echo "Certificate validity dates extracted and saved to $output_file in JSON format (formatted with jq)"
else
  echo "$json_output_array" > "$output_file"
  echo "Certificate validity dates extracted and saved to $output_file in JSON format (unformatted - jq not found)"
  echo "Consider installing 'jq' for formatted JSON output."
fi
