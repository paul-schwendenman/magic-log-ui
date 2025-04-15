#! /usr/bin/env python
import csv
import time
import random
import argparse
import sys


def main():
    parser = argparse.ArgumentParser(description="Stream a CSV column with semi-random timing.")
    parser.add_argument('file', help='Path to the CSV file')
    parser.add_argument('--column', required=True, help='Name of the column to output')
    parser.add_argument('--min', type=float, default=0.2, help='Minimum sleep time in seconds')
    parser.add_argument('--max', type=float, default=0.5, help='Maximum sleep time in seconds')

    args = parser.parse_args()

    if args.min > args.max:
        parser.error("--min must be less than or equal to --max")

    try:
        with open(args.file, newline='') as csvfile:
            reader = csv.DictReader(csvfile)

            if args.column not in reader.fieldnames:
                print(f"Error: Column '{args.column}' not found in CSV headers: {reader.fieldnames}", file=sys.stderr)
                sys.exit(1)

            for row in reader:
                value = row.get(args.column)
                if value is not None:
                    print(value)
                else:
                    print(f"Warning: Missing value for column '{args.column}' in row: {row}", file=sys.stderr)
                time.sleep(random.uniform(args.min, args.max))
    except FileNotFoundError:
        print(f"Error: File '{args.file}' not found.", file=sys.stderr)
        sys.exit(1)

if __name__ == '__main__':
    main()
