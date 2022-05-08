import json

lines  = []
with open('/home/rafaquat/buyer_ids.txt') as f:
    lines = f.readlines()

#print(lines) 
#print("\n")

id_strs = []
for id in lines:
    id = id.replace("\n", "")
    id_strs.append(id)

#print(id_strs)


x = "run double auction"
y = json.dumps(x)
print(y)



"""
# dummy test data
order_data =  '{ "buyerid":"6235eada63fc9762da8b7e54", "energyamount":99, "fiatamount":9.109, "requestclosed":True, "reqtime":"test time", "reqid":"", "auctioned":False }'

# make json
data = json.loads(order_data)


# a Python object (dict):
x = {
  "name": "John",
  "age": 30,
  "city": "New York"
}

# convert into JSON:
y = json.dumps(x)

"""