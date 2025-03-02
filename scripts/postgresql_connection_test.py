#!/usr/bin/env python3

import psycopg2
import os
from dotenv import load_dotenv

load_dotenv()

def connection_test():
    """
    Connects to a PostgreSQL database using credentials from environment variables
    loaded from a .env file and executes a simple query.

    Reads database connection details from the following environment variables:
        DB_HOST: Hostname or IP address of PostgreSQL database.
        DB_PORT: Port number of your PostgreSQL database (usually 5432).
        DB_NAME: Name of the database to connect to.
        DB_USER: Username for database access.
        DB_PASSWORD: Password for the database user.

    Returns:
        str: PostgreSQL version string if connection is successful, None otherwise.
    """
    db_host = os.environ.get("DB_HOST")
    db_port = os.environ.get("DB_PORT", "5432")  # Default port if not in .env
    db_name = os.environ.get("DB_NAME")
    db_user = os.environ.get("DB_USER")
    db_password = os.environ.get("DB_PASSWORD")

    if not all([db_host, db_name, db_user, db_password]):
        print("Error: Database connection details are missing in .env file.")
        print("Please ensure DB_HOST, DB_PORT, DB_NAME, DB_USER, and DB_PASSWORD are set in your .env file.")
        return None

    conn = None
    try:
        # Connect to the database
        conn = psycopg2.connect(
            host=db_host,
            port=db_port,
            database=db_name,
            user=db_user,
            password=db_password,
            sslmode='require'  # Enforce SSL
        )

        # Create a cursor object to interact with the database
        cur = conn.cursor()

        # Execute a simple query to check the connection
        cur.execute("SELECT VERSION();")

        # Fetch the version result
        db_version = cur.fetchone()[0]

        print("Successfully connected to PostgreSQL database using .env variables!")
        print(f"PostgreSQL database version: {db_version}")

        return db_version

    except (Exception, psycopg2.Error) as error:
        print("Error while connecting to PostgreSQL:", error)
        return None

    finally:
        # Close the database connection
        if conn:
            cur.close()
            conn.close()
            print("Database connection closed.")

if __name__ == "__main__":
    connection_test()
