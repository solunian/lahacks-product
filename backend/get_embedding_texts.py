import pandas as pd
import urllib.parse

drugs_drugs = pd.read_csv("drugs.csv") # drugName,url,description
wikipedia_drugs = pd.read_csv("wikipedia_drugs.csv") # ,id,title,first_paragraph
wikipedia_virus = pd.read_csv("wikipedia_virus.csv") # ,id,title,first_paragraph

wikipedia_drugs['url'] = "https://en.wikipedia.org/wiki/" +\
        wikipedia_drugs['title'].map(lambda t:
            urllib.parse.quote_plus(t.replace(" ", "_")))
wikipedia_virus['url'] = "https://en.wikipedia.org/wiki/" +\
        wikipedia_virus['title'].map(lambda t:
            urllib.parse.quote_plus(t.replace(" ", "_")))


# Rename columns in drugs_drugs
drugs_drugs = drugs_drugs.rename(columns={'drugName': 'name'})

# Rename columns in wikipedia_drugs
wikipedia_drugs = wikipedia_drugs.rename(columns={'title': 'name', 'first_paragraph': 'description'})

# Rename columns in wikipedia_virus
wikipedia_virus = wikipedia_virus.rename(columns={'title': 'name', 'first_paragraph': 'description'})


combined_df = pd.concat([drugs_drugs, wikipedia_drugs, wikipedia_virus], axis=0, ignore_index=True)

export_columns = ['name', 'url', 'description']
combined_df = combined_df[export_columns]

combined_df.index.name = 'id'
combined_df.reset_index()

combined_df.to_csv('for_embeddings.csv', index_label="id")
combined_df.to_json('for_embeddings.json', orient="index")
