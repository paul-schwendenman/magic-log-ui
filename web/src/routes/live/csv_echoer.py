import csv
import time
import random

# Path to your CSV file
csv_path = 'your_logs.csv'
column_name = 'your_column'  # Change this to the column you want

# Open and read the CSV
with open(csv_path, newline='') as csvfile:
    reader = csv.DictReader(csvfile)
    for row in reader:
        print(row[column_name])
        time.sleep(random.uniform(0.2, 0.5))
