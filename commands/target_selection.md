# target selection  

bounty-targets-data
[https://github.com/arkadiyt/bounty-targets-data/tree/main](https://github.com/arkadiyt/bounty-targets-data/tree/main)

## HackerOne Platform

json data:  
[https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/refs/heads/main/data/hackerone_data.json](https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/refs/heads/main/data/hackerone_data.json)

## In Scope URLs

get in_scope URL targets:  
`cat hackerone_data.json | jq -r '.[] | .targets.in_scope[] | select(.eligible_for_bounty == true and .asset_type == "URL") | .asset_identifier' | sort -u | anew hackerone_targets.txt`  

Use curl to get_in_scope URL targets from the raw file hosted on githubusercontent.com  
`curl https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/refs/heads/main/data/hackerone_data.json | jq -r '.[] | .targets.in_scope[] | select(.eligible_for_bounty == true and .asset_type == "URL") | .asset_identifier' | sort -u | anew hackerone_targets.txt`  

## In Scope wildcard URLs  

get in_scope URLS where targets start with "*":  
`cat hackerone_data.json | jq -r '.[] | .targets.in_scope[] | select(.eligible_for_bounty == true and .asset_type == "URL" and (.asset_identifier | startswith("*"))) | .asset_identifier' | sort -u | anew hackerone_targets_wildcard.txt`  
