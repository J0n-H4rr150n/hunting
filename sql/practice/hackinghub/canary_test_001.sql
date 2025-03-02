select 
	created_at, 
	id, 
	tool_name, 
	raw_data -> 'request' ->> 'method' as method, 
	raw_data -> 'response' ->> 'status_code' as status, 
	raw_data -> 'request' ->> 'url' as url, 
	raw_data as data 
from hunting.cli_tool_data 
where raw_data -> 'request' ->> 'url' like '%ctfio%' 
and raw_data -> 'request' ->> 'url' like '%19628397001%' 
order by id desc 
