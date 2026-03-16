#!/usr/bin/env python3
"""
Generate sports seed data from US Collegiate Sports Dataset

Usage:
    python seed_sports.py [--limit 100]

Example:
    python seed_sports.py --limit 100  # Test with 100 records
    python seed_sports.py              # Process all records
"""

import pandas as pd
import json
import argparse
import kagglehub


def generate_sports_data(df, output="data/sports.json", limit=None):
    """
    Generate sports.json from the US Collegiate Sports dataset.
    """
    print(f"\nProcessing US Collegiate Sports data...")
    print(f"Total records: {len(df)}")

    if limit:
        df = df.head(limit)
        print(f"Limited to first {limit} records")

    # Extract unique sports
    print("\n--- Extracting Sports ---")

    # Get unique sports
    unique_sports = df["sports"].dropna().unique()

    sports = []
    for sport_name in sorted(unique_sports):
        sport_name = str(sport_name).strip()
        if sport_name and sport_name != "nan":
            # Count occurrences as a rough popularity metric
            count = len(df[df["sports"] == sport_name])

            sports.append({"name": sport_name, "popularity": count})

    # Sort by popularity (most popular first)
    sports = sorted(sports, key=lambda x: x["popularity"], reverse=True)

    print(f"Extracted {len(sports)} unique sports")
    print(f"\nTop 10 sports by popularity:")
    for i, sport in enumerate(sports[:10], 1):
        print(f"  {i}. {sport['name']} ({sport['popularity']} programs)")

    # Write sports JSON
    with open(output, "w") as f:
        json.dump(sports, f, indent=2)

    print(f"\n✓ Saved to {output}")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Generate sports seed data from US Collegiate Sports Dataset"
    )
    parser.add_argument(
        "--limit",
        type=int,
        help="Limit number of records to process (for testing)",
    )
    parser.add_argument(
        "--output", default="data/sports.json", help="Output path for sports JSON"
    )

    args = parser.parse_args()

    # Load dataset from Kaggle
    print("Downloading US Collegiate Sports dataset from Kaggle...")
    print("(This may take a moment on first run)")

    dataset_path = kagglehub.dataset_download("umerhaddii/us-collegiate-sports-dataset")
    print(f"✓ Dataset downloaded to: {dataset_path}")

    # Find CSV file
    from pathlib import Path

    csv_files = list(Path(dataset_path).rglob("*.csv"))

    if not csv_files:
        raise FileNotFoundError(f"No CSV files found in {dataset_path}")

    csv_file = csv_files[0]
    print(f"Using file: {csv_file.name}")

    # Load the CSV
    df = pd.read_csv(str(csv_file))

    print(f"✓ Dataset loaded")

    # Generate sports data
    generate_sports_data(df, args.output, args.limit)

    print("\n✅ Sports data generation complete!")
