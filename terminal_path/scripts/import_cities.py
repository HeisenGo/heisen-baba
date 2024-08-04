import csv
import sqlite3
import os
import sys

# Add the project root to the Python path
project_root = os.path.abspath(os.path.join(os.path.dirname(__file__), '..'))
sys.path.insert(0, project_root)

# Ensure the db_helper directory exists
os.makedirs(os.path.join(project_root, 'db_helper'), exist_ok=True)

# Connect to the SQLite database (it will be created if it doesn't exist)
db_path = os.path.join(project_root, 'db_helper', 'worldcities.db')
conn = sqlite3.connect(db_path)
cursor = conn.cursor()

# Create the table if it doesn't exist
cursor.execute('''CREATE TABLE IF NOT EXISTS worldcities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    city TEXT NOT NULL,
    country TEXT NOT NULL
)
''')

# Create indexes for faster searching
cursor.execute('CREATE INDEX IF NOT EXISTS idx_city ON worldcities(city)')
cursor.execute('CREATE INDEX IF NOT EXISTS idx_country ON worldcities(country)')

# Open and read the CSV file
csv_path = os.path.join(project_root, 'data', 'worldcities_filtered.csv')
with open(csv_path, 'r', encoding='utf-8') as csv_file:
    csv_reader = csv.reader(csv_file)
    next(csv_reader)  # Skip the header row

    # Prepare the INSERT statement
    cursor.execute('BEGIN TRANSACTION')
    for row in csv_reader:
        cursor.execute('INSERT INTO worldcities (city, country) VALUES (?, ?)', (row[0], row[1]))

# Commit the changes and close the connection
conn.commit()
conn.close()

print("Data import completed successfully!")

# Verify the import
conn = sqlite3.connect(db_path)
cursor = conn.cursor()

# Count the total number of rows
cursor.execute('SELECT COUNT(*) FROM worldcities')
total_rows = cursor.fetchone()[0]
print(f"Total number of imported rows: {total_rows}")

# Display the first 5 rows
print("\nFirst 5 rows of imported data:")
cursor.execute('SELECT * FROM worldcities LIMIT 5')
for row in cursor.fetchall():
    print(row)

conn.close()