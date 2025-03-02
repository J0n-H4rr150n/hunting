import psycopg2
import psycopg2.extras
import os
import json
from dotenv import load_dotenv
from mitmproxy import http
from mitmproxy import addonmanager

# Load environment variables from .env file
load_dotenv()

class PostgresFlowSaver:
    """
    Mitmproxy addon to save HTTP flows (requests and responses) to a PostgreSQL database
    with a JSONB column.
    """
    def __init__(self):
        self.db_host = os.environ.get("DB_HOST")
        self.db_port = os.environ.get("DB_PORT", "5432")
        self.db_name = os.environ.get("DB_NAME")
        self.db_user = os.environ.get("DB_USER")
        self.db_password = os.environ.get("DB_PASSWORD")

        if not all([self.db_host, self.db_name, self.db_user, self.db_password]):
            #print("Error: Database connection details are missing in .env file.")
            #print("Please ensure DB_HOST, DB_PORT, DB_NAME, DB_USER, and DB_PASSWORD are set in your .env file.")
            self.enabled = False # Disable addon if DB config is missing
        else:
            self.enabled = True

    def load(self, loader: addonmanager.Loader):
        if not self.enabled:
            #print("PostgresFlowSaver addon is disabled due to missing database configuration.")
            return
        #print("PostgresFlowSaver addon loaded and ready to save flows to PostgreSQL.")

    def response(self, flow: http.HTTPFlow):
        if not self.enabled:
            return # Do nothing if addon is disabled

        request_content_decoded = None
        response_content_decoded = None

        if flow.request.content:
            # Remove null bytes and then decode, replacing errors
            request_content_decoded = flow.request.content.replace(b'\x00', b'').decode('utf-8', 'replace')

        if flow.response.content:
            # Remove null bytes and then decode, replacing errors
            response_content_decoded = flow.response.content.replace(b'\x00', b'').decode('utf-8', 'replace')


        flow_data = {
            "request": {
                "method": flow.request.method,
                "url": flow.request.url,
                "headers": dict(flow.request.headers.items()), # Convert headers to dict for JSON
                "content": request_content_decoded, # Use decoded content
            },
            "response": {
                "status_code": flow.response.status_code,
                "reason": flow.response.reason,
                "headers": dict(flow.response.headers.items()), # Convert headers to dict for JSON
                "content": response_content_decoded, # Use decoded content
                "timestamp_start": str(flow.response.timestamp_start), # Serialize timestamps to string for JSON
                "timestamp_end": str(flow.response.timestamp_end),     # Serialize timestamps to string for JSON
            },
            # Safely get client IP, handling potential None values, missing ip_address, or empty ip_address sequence
            "client_ip": flow.client_conn.ip_address[0] if flow.client_conn and hasattr(flow.client_conn, 'ip_address') and isinstance(flow.client_conn.ip_address, (tuple, list)) and len(flow.client_conn.ip_address) > 0 else None,
            # Safely get server IP, handling potential None values, missing ip_address, or empty ip_address sequence
            "server_ip": flow.server_conn.ip_address[0] if flow.server_conn and hasattr(flow.server_conn, 'ip_address') and isinstance(flow.server_conn.ip_address, (tuple, list)) and len(flow.server_conn.ip_address) > 0 else None,
            "timestamp_flow_start": str(flow.timestamp_start), # Serialize flow start time to string for JSON
        }

        self.insert_flow_data_to_postgres(flow_data)

    def insert_flow_data_to_postgres(self, flow_data):
        """Inserts mitmproxy flow data (dictionary) into PostgreSQL JSONB column in cli_tool_data table."""
        conn = None
        try:
            conn = psycopg2.connect(
                host=self.db_host,
                port=self.db_port,
                database=self.db_name,
                user=self.db_user,
                password=self.db_password,
                sslmode='require'  # Enforce SSL for DigitalOcean Managed DB
            )
            cursor = conn.cursor(cursor_factory=psycopg2.extras.DictCursor)

            sql = "INSERT INTO hunting.cli_tool_data (tool_name, raw_data) VALUES (%s, %s)" # Updated table name and added tool_name
            cursor.execute(sql, ('mitmproxy', json.dumps(flow_data),)) # Added 'mitmproxy' as tool_name

            conn.commit()
            #print("Flow data saved to PostgreSQL (hunting.cli_tool_data table).") # Updated message

        except psycopg2.Error as e:
            #print(f"Database error while saving flow data: {e}")
            if conn:
                conn.rollback()
        except json.JSONEncodeError as e:
            #print(f"JSON Encoding error while saving flow data: {e}")
            pass
        finally:
            if conn:
                cursor.close()
                conn.close()


addons = [
    PostgresFlowSaver()
