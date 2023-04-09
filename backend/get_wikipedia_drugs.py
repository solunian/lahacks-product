import requests
import pandas as pd
from tqdm import tqdm

def add_suffix_if_not_present(string, suffix):
    """
    Add a suffix to a string if the string does not already end with that suffix.

    Args:
        string (str): The string to add a suffix to.
        suffix (str): The suffix to add.

    Returns:
        str: The original string with the suffix added if it was not already present.
    """
    if not string.endswith(suffix):
        string += suffix
    return string

def get_all_medication_ids(): 
    # ?medication wdt:P31 wd:Q12140 . # medication
    # ?medication wdt:P31 wd:Q112193867 . # virus

    query = """
    SELECT DISTINCT ?medication
    WHERE {
        ?medication wdt:P31 wd:Q12140 .
    }
    """
    url = f"https://query.wikidata.org/sparql?query={query}&format=json"
    response = requests.get(url)
    data = response.json()
    medication_ids = [item['medication']['value'].split('/')[-1] for item in data['results']['bindings']]
    return medication_ids

def get_first_paragraphs(wikidata_ids):
    # initialize list to store results
    results = []
    # chunk the Wikidata IDs into groups of 50
    chunked_ids = [wikidata_ids[i:i+50] for i in range(0, len(wikidata_ids), 50)]
    # iterate over each chunk of Wikidata IDs
    for chunk in tqdm(chunked_ids):
        try:
            # join list of Wikidata IDs into a comma-separated string
            wikidata_ids_string = "|".join(chunk)
            # make API request to Wikidata to get English Wikipedia article titles and URLs
            url = f"https://www.wikidata.org/w/api.php?action=wbgetentities&ids={wikidata_ids_string}&format=json&props=sitelinks&sitefilter=enwiki"
            
            response = requests.get(url)
            data = response.json()

            # extract the article titles from the Wikidata API response
            titles = []
            for wikidata_id in chunk: 
                try:
                    titles.append(data['entities'][wikidata_id]['sitelinks']['enwiki']['title'])
                except KeyError:
                    pass
            
            if titles == []:
                print(url, response.status_code, response.text)
                continue
            
            # make API request to Wikipedia to get article content for all titles in the current chunk
            titles_string = "|".join(titles)
            url = f"https://en.wikipedia.org/w/api.php?action=query&titles={titles_string}&prop=extracts&format=json&exintro=1&explaintext=1"
            response = requests.get(url)
            data = response.json()
            # iterate over each Wikidata ID in the current chunk
            for article in data['query']['pages'].values():
                title = article['title']
                try:
                    first_paragraph = add_suffix_if_not_present(". ".join(article['extract'].split("\n")[0].split('. ')[0:3]).strip(), ".")
                except KeyError:
                    continue
                if len(first_paragraph) < 100:
                    continue
                results.append({
                    "id": article['pageid'],
                    "title": title,
                    "first_paragraph": first_paragraph
                })
        except KeyError as e:
            print(e)
            print(url, response.text)
            # if English Wikipedia article not found, skip and continue to next chunk
            continue
    return results


medication_ids = get_all_medication_ids()
medications = get_first_paragraphs(medication_ids)

df = pd.DataFrame(medications)

df.to_csv("wikipedia_drugs.csv")