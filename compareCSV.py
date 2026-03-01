#!/usr/bin/env python3

import csv
from collections import Counter


def load_csv_as_counter(path):
    """
    Loads CSV into a Counter of row tuples.
    Handles duplicates correctly.
    Returns headers and Counter.
    """

    with open(path, newline="", encoding="utf-8") as f:
        reader = csv.DictReader(f)
        headers = reader.fieldnames

        counter = Counter(
            tuple(row[h] for h in headers)
            for row in reader
        )

    return headers, counter


def find_differences(headers, c1, c2):
    """
    Finds missing and extra rows.
    Returns list of diff entries.
    """

    diffs = []

    all_rows = set(c1.keys()) | set(c2.keys())

    for row in all_rows:
        count1 = c1.get(row, 0)
        count2 = c2.get(row, 0)

        if count1 > count2:
            for _ in range(count1 - count2):
                diffs.append(("Missing in CSV2", row, None))

        elif count2 > count1:
            for _ in range(count2 - count1):
                diffs.append(("Extra in CSV2", None, row))

    return diffs


def print_diff_table(headers, diffs, max_rows=50):
    """
    Pretty prints differences.
    """

    if not diffs:
        print("✓ Files are identical (order independent)")
        return

    print(f"\nFound {len(diffs)} differences\n")

    for i, (dtype, row1, row2) in enumerate(diffs[:max_rows], 1):

        print("=" * 60)
        print(f"Difference #{i}: {dtype}")

        if row1:
            print("\nCSV1:")
            for h, v in zip(headers, row1):
                print(f"{h:20} {v}")

        if row2:
            print("\nCSV2:")
            for h, v in zip(headers, row2):
                print(f"{h:20} {v}")

        if row1 and row2:
            print("\nColumn differences:")
            for h, v1, v2 in zip(headers, row1, row2):
                if v1 != v2:
                    print(f"{h:20} {v1}  !=  {v2}")

        print()

    if len(diffs) > max_rows:
        print(f"... and {len(diffs) - max_rows} more differences")


def compare_csv(csv1, csv2):

    headers1, counter1 = load_csv_as_counter(csv1)
    headers2, counter2 = load_csv_as_counter(csv2)

    if headers1 != headers2:
        print("ERROR: Headers differ")
        print(headers1)
        print(headers2)
        return

    diffs = find_differences(headers1, counter1, counter2)

    print_diff_table(headers1, diffs)


if __name__ == "__main__":

    import argparse

    parser = argparse.ArgumentParser(description="Compare two CSV files ignoring row order")
    parser.add_argument("csv1")
    parser.add_argument("csv2")

    args = parser.parse_args()

    compare_csv(args.csv1, args.csv2)
