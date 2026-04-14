#!/usr/bin/env python3
"""
Fetch NCAA school logos and upload them to S3.

Matches each college in colleges.json to a school in the NCAA /schools-index,
downloads the logo image, and uploads it to S3 under colleges/logos/<seo_slug>.<ext>.

Usage:
    python seed_college_logos.py [--colleges data/colleges.json] [--limit 10] [--dry-run]

Requirements:
    pip install boto3 requests python-dotenv thefuzz

Env vars (matches your Go S3 service):
    S3_BUCKET       - S3 bucket name (required)
    AWS_REGION      - AWS region (required)
"""

import argparse
import json
import os
import time
import unicodedata
import re
import sys
from pathlib import Path

import boto3
import requests
from thefuzz import process as fuzz_process

NCAA_API_BASE = "https://ncaa-api.henrygd.me"
NCAA_CDN_BASE = "https://i.turner.ncaa.com/sites/default/files/images/logos/schools/bgd"
RATE_LIMIT_DELAY = 0.25  # 4 req/s, safely under the 5 req/s cap
FUZZY_MATCH_THRESHOLD = 80  # minimum score to accept a name match
S3_KEY_PREFIX = "colleges/logos"
CONTENT_TYPE = "image/png"


def normalize(name: str) -> str:
    """Lowercase, strip accents, remove punctuation for fuzzy comparison."""
    name = unicodedata.normalize("NFKD", name)
    name = "".join(c for c in name if not unicodedata.combining(c))
    name = re.sub(r"[^\w\s]", "", name).lower().strip()
    return name


def load_colleges(path: str) -> list[dict]:
    with open(path) as f:
        return json.load(f)


def fetch_schools_index() -> list[dict]:
    """
    GET /schools-index — returns all NCAA schools.
    Each entry looks like:
      { "name": "Abilene Christian", "seo": "abilene-christian", ... }
    """
    print("Fetching NCAA schools index...")
    resp = requests.get(f"{NCAA_API_BASE}/schools-index", timeout=15)
    resp.raise_for_status()
    data = resp.json()
    # The endpoint returns { "schools": [ { "name": ..., "seo": ... }, ... ] }
    schools = data.get("schools", data) if isinstance(data, dict) else data
    print(f"  ✓ {len(schools)} schools returned")
    return schools


def build_lookup(schools: list[dict]) -> dict[str, dict]:
    """Build a normalized-name → school dict using both short and long names."""
    lookup = {}
    for s in schools:
        if s.get("long"):
            lookup[normalize(s["long"])] = s
        if s.get("name"):
            lookup[normalize(s["name"])] = s
    return lookup


def match_college(college_name: str, lookup: dict[str, dict]) -> dict | None:
    """Fuzzy-match a college name against the NCAA schools index."""
    normalized = normalize(college_name)
    result = fuzz_process.extractOne(normalized, lookup.keys())
    if result is None:
        return None
    matched_key, score = result[0], result[1]
    if score < FUZZY_MATCH_THRESHOLD:
        return None
    return lookup[matched_key]


def fetch_logo(seo_slug: str) -> tuple[bytes, str] | tuple[None, None]:
    """
    Try to download the logo for a school from the NCAA CDN.
    NCAA CDN pattern: https://i.turner.ncaa.com/sites/default/files/images/logos/schools/bgd/<seo>.svg
    Falls back to .png if .svg is not found.
    Returns (image_bytes, extension) or (None, None).
    """
    for ext in ("svg", "png"):
        url = f"{NCAA_CDN_BASE}/{seo_slug}.{ext}"
        try:
            resp = requests.get(url, timeout=10)
            if resp.status_code == 200 and len(resp.content) > 0:
                return resp.content, ext
        except requests.RequestException:
            continue
    return None, None


def s3_key(seo_slug: str, ext: str) -> str:
    """Mirrors the key format from your Go service: colleges/logos/<seo_slug>.<ext>"""
    return f"{S3_KEY_PREFIX}/{seo_slug}.{ext}"


def upload_to_s3(s3_client, bucket: str, key: str, data: bytes, ext: str) -> str:
    """Upload logo bytes to S3 and return the public S3 URI."""
    content_type = "image/svg+xml" if ext == "svg" else "image/png"
    s3_client.put_object(
        Bucket=bucket,
        Key=key,
        Body=data,
        ContentType=content_type,
        Metadata={"source": "ncaa-seed-script"},
    )
    region = s3_client.meta.region_name
    return f"https://{bucket}.s3.{region}.amazonaws.com/{key}"


# ── Main ─────────────────────────────────────────────────────────────────────


def main():
    parser = argparse.ArgumentParser(description="Fetch NCAA logos and upload to S3")
    parser.add_argument(
        "--colleges", default="data/colleges.json", help="Path to colleges.json"
    )
    parser.add_argument(
        "--limit", type=int, help="Only process first N colleges (for testing)"
    )
    parser.add_argument(
        "--dry-run", action="store_true", help="Skip S3 upload, just print matches"
    )
    args = parser.parse_args()

    bucket = os.environ.get("S3_BUCKET")
    region = os.environ.get("AWS_REGION")

    if not bucket or not region:
        print("ERROR: S3_BUCKET and AWS_REGION env vars are required")
        sys.exit(1)

    if not args.dry_run:
        # Uses the default AWS credential chain — same as your Go service
        s3_client = boto3.client("s3", region_name=region)
        print(f"✓ S3 client ready (bucket={bucket}, region={region})")
    else:
        s3_client = None
        print(f"[DRY RUN] Would upload to bucket={bucket}, region={region}")

    colleges = load_colleges(args.colleges)
    if args.limit:
        colleges = colleges[: args.limit]
    print(f"\n✓ Loaded {len(colleges)} colleges from {args.colleges}")

    schools = fetch_schools_index()
    lookup = build_lookup(schools)

    results = {"uploaded": [], "not_matched": [], "no_logo": [], "failed": []}

    for i, college in enumerate(colleges, 1):
        name = college["name"]
        print(f"\n[{i}/{len(colleges)}] {name}")

        # 1. Match to NCAA school
        school = match_college(name, lookup)
        if not school:
            print(f"  ✗ No NCAA match found (threshold={FUZZY_MATCH_THRESHOLD})")
            results["not_matched"].append(name)
            continue

        seo_slug = school.get("slug", "")
        if not seo_slug:
            print(f"  ✗ Matched '{school.get('name')}' but no seo slug")
            results["not_matched"].append(name)
            continue

        print(f"  → Matched: {school['name']} (seo={seo_slug})")

        # 2. Fetch logo from NCAA CDN
        time.sleep(RATE_LIMIT_DELAY)
        logo_bytes, ext = fetch_logo(seo_slug)
        if not logo_bytes:
            print(f"  ✗ No logo found on NCAA CDN")
            results["no_logo"].append(name)
            continue

        print(f"  ✓ Logo fetched ({len(logo_bytes)} bytes, .{ext})")

        if args.dry_run:
            key = s3_key(seo_slug, ext)
            print(f"  [DRY RUN] Would upload to s3://{bucket}/{key}")
            results["uploaded"].append({"name": name, "key": key})
            continue

        # 3. Upload to S3
        try:
            key = s3_key(seo_slug, ext)
            s3_url = upload_to_s3(s3_client, bucket, key, logo_bytes, ext)
            print(f"  ✓ Uploaded → {s3_url}")
            results["uploaded"].append({"name": name, "key": key, "url": s3_url})
            college["logo"] = key
        except Exception as e:
            print(f"  ✗ S3 upload failed: {e}")
            results["failed"].append({"name": name, "error": str(e)})

    print("\n" + "=" * 60)
    print("SUMMARY")
    print("=" * 60)
    print(f"  ✓ Uploaded:     {len(results['uploaded'])}")
    print(f"  ✗ No match:     {len(results['not_matched'])}")
    print(f"  ✗ No logo:      {len(results['no_logo'])}")
    print(f"  ✗ Upload error: {len(results['failed'])}")

    if results["not_matched"]:
        print(f"\nUnmatched colleges:")
        for n in results["not_matched"]:
            print(f"  - {n}")

    if results["no_logo"]:
        print(f"\nNo logo found:")
        for n in results["no_logo"]:
            print(f"  - {n}")

    with open(args.colleges, "w") as f:
        json.dump(colleges, f, indent=2)
    print(f"✓ colleges.json updated with logo keys")

    # Write results to JSON for reviewa
    output_path = Path(args.colleges).parent / "logo_results.json"
    with open(output_path, "w") as f:
        json.dump(results, f, indent=2)
    print(f"\n✓ Full results saved to {output_path}")


if __name__ == "__main__":
    main()
