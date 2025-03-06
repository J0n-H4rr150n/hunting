# target selection  

bounty-targets-data
[https://github.com/arkadiyt/bounty-targets-data/tree/main](https://github.com/arkadiyt/bounty-targets-data/tree/main)

## HackerOne Platform

json data:  
[https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/refs/heads/main/data/hackerone_data.json](https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/refs/heads/main/data/hackerone_data.json)

get in_scope URL targets:  
```
cat hackerone_data.json | jq -r '.[] | .targets.in_scope[] | select(.eligible_for_bounty == true and .asset_type == "URL") | .asset_identifier' | sort -u | anew hackerone_targets.txt
```
