import math
import weaviate
import numpy as np
import cohere
import pandas as pd
import json
import os
from tqdm import tqdm

client = weaviate.Client("http://localhost:8080")

schema = {
    "classes": [{
            "class": "MedObject",
            "vectorizer": "none", # explicitly tell Weaviate not to vectorize anything, we are providing the vectors ourselves through our BERT model
            "properties": [{
                "name": "doc_id",
                "dataType": ["text"],
            },{
                "name": "name",
                "dataType": ["text"],
            },{
                "name": "url",
                "dataType": ["text"],
            },{
                "name": "description",
                "dataType": ["text"],
            }]
    }]
}

# cleanup from previous runs
client.schema.delete_all()

client.schema.create(schema)


client.batch.configure(
    # `batch_size` takes an `int` value to enable auto-batching
    # (`None` is used for manual batching)
    batch_size=100,
    # dynamically update the `batch_size` based on import speed
    dynamic=False,
    # `timeout_retries` takes an `int` value to retry on time outs
    timeout_retries=3,
    # checks for batch-item creation errors
    # this is the default in weaviate-client >= 3.6.0
    callback=weaviate.util.check_batch_result,
    consistency_level=weaviate.data.replication.ConsistencyLevel.ALL,  # default QUORUM
)


co = cohere.Client(os.getenv("COHERE_API_KEY"))

# Load text data into dataframe
data = pd.read_csv('for_embeddings.csv')

data = data.reset_index()

embeddings_with_ids={}

chunk_size = 95
chunks = math.floor(len(data) / chunk_size + 1)


for chunk, num in tqdm(
    ((data[x * chunk_size:(x + 1) * chunk_size], x) for x in range(chunks)),
        total=chunks):
    embeddings = co.embed(texts=list(chunk['description']))
    print(len(chunk), len(embeddings))
    for i, e in enumerate(embeddings):
        entity = chunk.loc[i+(chunk_size*num)]        
        embeddings_with_ids[str(entity['id'])] = {
            "name": entity['name'],
            "url": entity['url'],
            "description": entity['description'],
            "embedding": e
        }

# Save embeddings with IDs to file
with open('embeddings_with_ids.json', 'w') as f:
    json.dump(embeddings_with_ids, f)

with client.batch as batch:
    # Create Weaviate objects for each text and insert them into the database
    for i, v in enumerate(embeddings_with_ids.items()):
        id, data = v
        object = {
            "doc_id": id,
            "name": data['name'],
            "url": data['url'],
            "description": data['description'],
        }
        batch.add_data_object(object, "MedObject", vector=data['embedding'])
