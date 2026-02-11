#!/usr/bin/env python3
"""
Generate college seed data from NCAA College Division Database

Usage:
    python seed_colleges.py [--limit 20] [--output data/colleges.json] [--division 1]

Example:
    python seed_colleges.py --limit 20 --division 1  # Test with 20 Division 1 schools
    python seed_colleges.py --division 1              # All Division 1 schools
"""

import pandas as pd
import json
import argparse
import kagglehub


def transform_colleges(df, output_path="data/colleges.json", limit=None, division=None):
    """
    Transform NCAA college DataFrame to colleges JSON format.
    """
    print(f"\nProcessing NCAA college data...")
    print(f"Columns available: {df.columns.tolist()}")
    print(f"First few records:\n{df.head()}")

    # Filter by division if specified
    if division:
        # Check which column contains division info
        division_col = None
        for col in df.columns:
            if "division" in col.lower() or "div" in col.lower():
                division_col = col
                break

        if division_col:
            original_count = len(df)
            # Try different formats: "Division I", "D1", "I", "1", etc.
            df = df[
                (
                    df[division_col]
                    .astype(str)
                    .str.contains(f"Division {division}", case=False, na=False)
                )
                | (
                    df[division_col]
                    .astype(str)
                    .str.contains(f"D{division}", case=False, na=False)
                )
                | (df[division_col].astype(str) == str(division))
                | (df[division_col].astype(str) == "I" if division == 1 else False)
            ]
            print(
                f"Filtered from {original_count} to {len(df)} Division {division} schools"
            )
        else:
            print(f"⚠️  Warning: Could not find division column")

    print(f"Total records after filtering: {len(df)}")

    # Apply limit if specified
    if limit:
        df = df.head(limit)
        print(f"Limited to first {limit} records")

    colleges = []

    for _, row in df.iterrows():
        # Try to find the right columns - adapt based on what's in the dataset
        name = None
        state = None
        city = None
        website = None
        conference = None

        # Look for name column (case-insensitive)
        for col in [
            "Name",
            "name",
            "school",
            "School",
            "institution",
            "college",
            "INSTNM",
        ]:
            if col in df.columns:
                name = str(row.get(col, "")).strip()
                if name and name != "nan":
                    break

        # Look for Location column (format: "City, State")
        if "Location" in df.columns:
            location = str(row.get("Location", "")).strip()
            if location and location != "nan" and "," in location:
                parts = location.split(",")
                city = parts[0].strip()
                state = parts[1].strip() if len(parts) > 1 else ""

        # Fall back to separate state/city columns if Location doesn't exist
        if not state:
            for col in ["state", "State", "st", "state_abbr", "STABBR"]:
                if col in df.columns:
                    state = str(row.get(col, "")).strip()
                    if state and state != "nan":
                        break

        if not city:
            for col in ["city", "City", "CITY"]:
                if col in df.columns:
                    city = str(row.get(col, "")).strip()
                    if city and city != "nan":
                        break

        # Look for website column (case-insensitive)
        for col in ["URL", "url", "website", "Website", "web", "INSTURL"]:
            if col in df.columns:
                website = str(row.get(col, "")).strip()
                if website and website != "nan":
                    if not website.startswith("http"):
                        website = "https://" + website
                    break

        # Look for conference column (case-insensitive)
        for col in ["Conference", "conference", "conf"]:
            if col in df.columns:
                conference = str(row.get(col, "")).strip()
                if conference and conference != "nan":
                    break

        # Skip if missing essential data
        if not name or not state:
            continue

        college = {
            "name": name,
            "state": state,
            "city": city or "",
            "website": website or "",
            "conference": conference or "",
            "division_rank": division if division else 1,
        }

        colleges.append(college)

    print(f"Transformed {len(colleges)} colleges")

    # Write to JSON
    with open(output_path, "w") as f:
        json.dump(colleges, f, indent=2)

    print(f"✓ Saved to {output_path}")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Generate college seed data from NCAA College Division Database"
    )
    parser.add_argument(
        "--output", default="data/colleges.json", help="Output JSON file path"
    )
    parser.add_argument(
        "--limit",
        type=int,
        help="Limit number of records to process (e.g., 20 for testing)",
    )
    parser.add_argument(
        "--division",
        type=int,
        choices=[1, 2, 3],
        help="Filter by NCAA division (1, 2, or 3)",
    )

    args = parser.parse_args()

    # Load dataset from Kaggle
    print("Loading NCAA College Division dataset from Kaggle...")
    print("(This may take a moment on first run)")

    # Download the dataset first
    dataset_path = kagglehub.dataset_download("flynn28/college-division-database")
    print(f"✓ Dataset downloaded to: {dataset_path}")

    # Find CSV files in the dataset
    from pathlib import Path

    csv_files = list(Path(dataset_path).rglob("*.csv"))
    if not csv_files:
        raise FileNotFoundError(f"No CSV files found in {dataset_path}")

    print(f"\nFound {len(csv_files)} CSV file(s):")
    for i, f in enumerate(csv_files, 1):
        print(f"  {i}. {f.name}")

    # Look for Colleges.csv specifically (NCAA schools)
    colleges_file = None
    for f in csv_files:
        if f.name.lower() == "colleges.csv":
            colleges_file = f
            break

    if not colleges_file:
        # Fall back to first CSV if Colleges.csv not found
        colleges_file = csv_files[0]

    print(f"\nUsing file: {colleges_file.name}")
    df = pd.read_csv(str(colleges_file))

    print(f"✓ Dataset loaded: {len(df)} total records")

    # Transform the data
    transform_colleges(df, args.output, args.limit, args.division)
