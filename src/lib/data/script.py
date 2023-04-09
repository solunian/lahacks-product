import json

meds = open('meds_short.json', "r")
out_meds = open("meds_out.json", "w")

meds_dict: dict = json.load(meds)
new_list = []

for key, in_dict in zip(meds_dict.keys(), meds_dict.values()):
    inside_dict = {"name": key}
    
    for in_key, in_value in zip(in_dict.keys(), in_dict.values()):
        inside_dict[in_key] = in_value[0:75]
        
    new_list.append(inside_dict)

json.dump(new_list, out_meds)