import pandas as pd
import json

RANDOM_STATE=sum([ord(c) for c in "e621.net"])

df = pd.read_csv("webmdQAs.csv")

alpaca_df = pd.read_parquet("alpaca.parquet")

df = pd.concat([
    df.sample(n=18_000, random_state=RANDOM_STATE),
    alpaca_df.sample(n=21_000, random_state=RANDOM_STATE)
])

df.to_csv("aaaa.csv")

with open('output.txt', 'w') as f:
    # iterate over each row of the dataframe
    for index, row in df.sample(frac=1).iterrows():
        if isinstance(row["text"], str):
            f.write(f"{row['text']}\n####\n\n")
            continue
        tags = ', '.join(json.loads(row['tags']))
        if len(tags) == 0:
            tags = '[NO TAGS]'
        f.write(f"### Human: {row['question']}\n\n### Assistant: {row['answer']}\n\n####\n\n")
