import csv

input_file = './data/worldcities.csv'
output_file = './data/worldcities_filtered.csv'

with open(input_file, mode='r', newline='', encoding='utf-8') as infile, \
     open(output_file, mode='w', newline='', encoding='utf-8') as outfile:
    reader = csv.DictReader(infile)
    writer = csv.writer(outfile)
    writer.writerow(['city_ascii', 'country'])
    for row in reader:
        writer.writerow([row['city_ascii'], row['country']])


# to run: python3 ./scripts/clean_csv.py
