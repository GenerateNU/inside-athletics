import json
import os
import re
from typing import Dict, List, Set, Tuple
import requests
from bs4 import BeautifulSoup
from pathlib import Path

class ScrapeAthletes:

    info_pattern: str
    sports: List[str]

    def __init__(self, info_pattern: str):
        self.info_pattern = info_pattern
        self.sports = self.read_sports_from_file()
        
    def read_sports_from_file(self) -> List[str]:
        """
        Read all the possible sports from the sports file, create possible url representations for given sports
        """
        sports_path = Path(__file__).parent / "sports.txt"
        print(sports_path)
        sports_list = []
        with open(sports_path, "r") as file:
            for line in file:

                # seperate line into sport information 
                line = line.lower()
                split = line.split(' ')
                gender_string = split.pop(-1).rstrip()

                # parse line into sport name and gender info
                genders: List[str] = gender_string[1: len(gender_string) - 1].split('/')
                sport_str = split[0]

                
                sport_str = "-".join(split)
                sports_list.append(sport_str) # add sport by itself
                sports_list = sports_list + [g + "s-" + sport_str for g in genders]
            return sports_list
        
    def clean_url(self, url: str) -> str:
        """
            Given a url, put it in https:// format 
        """
        if url.startswith(("http://", "https://")):
            return url
        elif url.startswith("www."):
            return url.replace("www.", "https://")
        else:
            return "https://" + url  
    
    def get_all_athletes(self, sample_size: int | None = None):
        college_info_path = Path(__file__).parent / "college_info.json" 
        missed_colleges = []
        college_athletes = {}
        total_schools_looped = 0
        with open(college_info_path, "r") as file:
            college_info = json.load(file)
            for college in college_info:
                sport_url = self.clean_url(college["athletics_website"])
                name = college["name"]
                athletes, found = self.get_athletes(sport_url)
                if not found:
                    missed_colleges.append(name)
                else:
                    college_athletes[name] = athletes
                
                total_schools_looped += 1
                if sample_size and total_schools_looped >= sample_size:
                    break

            with open(Path(__file__).parent / "athletes.json", "w", encoding="utf-8") as athlete_file:
                json.dump(college_athletes, athlete_file, indent=2)
            with open(Path(__file__).parent / "missed_schools.txt", "w") as missed_college_file:
                for school in missed_colleges:
                    missed_college_file.write(school + "\n")
            print(f"Number of schools with found athletes {len(college_athletes)} number of missed colleges {len(missed_colleges)}")
    
    def get_athletes(self, base_url: str) -> Tuple[Dict[str, List[str]], bool]:
        """
        Get the names of all athletes for a given base_url college 

        params:
            base_url: str: Represents the sport website for a given college 
        
        Returns: Dictionary of found sports for the given college sports url and the athlete names for each found sport along with if any sports were found for the given college
        """
        sports_url = f"{base_url}/sports/"
        # create links to respective sports rosters
        sports_rosters = {}
        for sport in self.sports:
            roster_url = f"{sports_url}{sport}/roster"
            try:
                roster: List[str] = self.filter_for_sublinks(self.create_soup_from_request(roster_url), self.info_pattern)
            except Exception as e:
                print(f"{sport} ERROR {roster_url} {e}")
                continue
            roster = self.get_athlete_names(roster)
            if len(roster) > 0:
                sports_rosters[sport] = roster
                print(f"{sport} FOUND {roster_url}")
            else:
                print(f"{sport} NOT FOUND {roster_url}")
        return sports_rosters, len(sports_rosters) > 0
    
    def get_college_info(self, save: bool) -> Dict[str, str]:
        """
        Get college information from the NCAA API, specifically used to get the athletics website for a given school, division, name, and school website

        Return Dictionary storing all the necessary information for a given college. 
        """
        # get cached list of athletes that luis got
        cached_list = None 
        try:
            with open(Path(__file__).parent.parent / "seed" / "data" / "colleges.json") as file:
                cached_list = json.load(file)
        except:
            print("Couldn't open cached list of colleges")
            pass

        # Use NCAA free API to scrape all colleges with their sport website information
        sport_websites: List[Dict[str, str]] = []
        ncaa_url = "https://web3.ncaa.org/directory/api/directory/memberList?type=12"
        response = requests.get(ncaa_url)
        if response.status_code == 200:
            college_information = response.json()
            for college_data in college_information:
                name: str = college_data["nameOfficial"]
                school_url: str = college_data["webSiteUrl"]
                sport_website: str = college_data["athleticWebUrl"]
                division: str = college_data["division"]
                city, state = self.get_address_information(name, cached_list=cached_list)
                if not (name and school_url and sport_website and division and city and state):
                    continue # all information must be present for us to add into database of schools

                if not sport_website: # filter out null sports websites 
                    continue

                sport_websites.append({
                    "name": name,
                    "school_url": school_url,
                    "athletics_website": sport_website,
                    "division": division,
                    "city": city,
                    "state": state
                })

            if save:
                save_path = Path(__file__).parent / "college_info.json"
                with open(save_path, "w", encoding="utf-8") as file:
                    json.dump(sport_websites, file, indent=2)

            return sport_websites
        else:
            raise Exception("Unable to get college information from NCAA website")
        
    def get_address_information(self, college_name: str, cached_list = None) -> Tuple[str | None, str | None]:
        """
        Given a college name return the city and state the college is in using the Google places API
        """
        for c in cached_list:
            if college_name == c["name"]:
                return c["city"], c["state"]

        url = "https://maps.googleapis.com/maps/api/geocode/json"
        params = {
            "address": college_name,
            "key": os.getenv("GOOGLE_MAPS_API_KEY")
        }
        response = requests.get(url, params=params)
        data = response.json()
        
        if data["results"]:
            formatted_addr = data["results"][0]["formatted_address"]
            spl: List[str] = formatted_addr.split(",")
            if len(spl) < 3:
                return None, None
            print(f"{college_name} {spl}")
            city = spl[1].strip()
            state_str = spl[2]
            white_space = state_str.index(' ')
            state = state_str[white_space + 1 : white_space + 3]
            return city, state
        return None, None

    def get_athlete_names(self, roster_links: List[str]) -> List[str]:
        """
        Parse athlete names from roster links and return them as a set
        """
        athlete_names = set()
        for link in roster_links:
            match = re.search(r"roster/[a-z]+-[a-z]+", link)
            if match:
                found_pattern = match.group()
                name = found_pattern.split("/")[1]
                athlete_names.add(name)
        return list(athlete_names)
    
    def create_soup_from_request(self, url: str) -> BeautifulSoup:
        """
        Get the HTML for a given link, return it parsed into BeautifulSoup object
        """
        page = requests.get(url)
        if page.status_code != 200:
            raise Exception(f"Unable to reach page, returned code: {page.status_code}")
        return BeautifulSoup(page.content, "html.parser")
    
    def filter_for_sublinks(self, soup: BeautifulSoup, pattern: str, filter: str | None = None):
        """
        Filter for sublinks based on the given pattern, checks all the href links for the given pattern
        """
        roster_links = soup.find_all("a", href=re.compile(pattern))
        l = []
        for link in roster_links:
            extracted_link = link.get("href")
            l.append(re.search(filter, extracted_link).group(0) if filter else extracted_link)
        return l
    

c = ScrapeAthletes(".*/roster/[a-z]+-[a-z]+/[0-9]+")
athlete_dict = c.get_all_athletes(sample_size=50)