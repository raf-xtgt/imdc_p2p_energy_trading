import numpy as np # linear algebra
# to forecast the energy data
import pandas as pd
import matplotlib.pyplot as plt
from statsmodels.tsa.api import ExponentialSmoothing, SimpleExpSmoothing, Holt
import pymongo
# to access the .env file
import os
from dotenv import load_dotenv
# to pass arguments to the script
import sys
# to determin points to take based on time
from datetime import datetime
# to increase the buy forecast amount by a little bit
import random
# import all functions from matchmaking.py
from sellerSelection import *  
from doubleAuction import *

load_dotenv()
USERID = ""
# handle arguments passed to the script
# total arguments
n = len(sys.argv)
print("Total arguments passed:", n)
 
# Arguments passed
print("\nName of Python script:", sys.argv[0])
 
print("\nArguments passed:", sys.argv[1], "\n") # userId or global
print("\nProcess:", sys.argv[2], "\n") # type of process(prod, consm, optm)

process = sys.argv[2] # what process needs to be run

# prod --> sellEnergyForecast(energy production forecast)
# consm --> buyEnergyForecast(energy consumption forecast)
# optm --> matchmaking algo


timeStamps = ["12:00AM", "12:30AM", "01:00AM","01:30AM","02:00AM","02:30AM", "03:00AM", "03:30AM",
"04:00AM", "04:30AM", "05:00AM", "05:30AM", "06:00AM", "06:30AM","07:00AM", "07:30AM",
"08:00AM", "08:30AM", "09:00AM", "09:30AM", "10:00AM", "10:30AM", "11:00AM", "11:30AM", 
"12:00PM","12:30PM", "01:00PM","01:30PM","02:00PM","02:30PM", "03:00PM", "03:30PM",
"04:00PM", "04:30PM", "05:00PM", "05:30PM", "06:00PM", "06:30PM","07:00PM", "07:30PM",
"08:00PM", "08:30PM", "09:00PM", "09:30PM", "10:00PM", "10:30PM", "11:00PM", "11:30PM",]

def energyForecast():
    upper_limit = findTime()
    print("Upper Limit:", upper_limit)
    # Load the data
    fields = ['PeriodStart', 'Ghi Curr Day'] # we only want to use these rows
    filename = "Irradiance_data_for_one_day.xlsx"
    df = pd.read_excel(filename, usecols=fields)
    #print(df.info) # print the dataframe

    # select a specific segment of the data
    data = df['Ghi Curr Day'][1:upper_limit+1]
    # training the model
    model = ExponentialSmoothing(data, trend="add", damped_trend=True, seasonal="add", seasonal_periods=2)
    model_fit = model.fit() # points according to the model
    #print("Forecast plots:", model_fit.fittedvalues)
    
    # get forecast/prediction for the next point
    forecast = model_fit.predict()
    #print("Actual values:", data)
    #print("Forecast values:", forecast)

    # Plot the forecast 
    actual_x = timeStamps[1:upper_limit+1]
    ypoints = df['Ghi Curr Day'][1:upper_limit+1].values
    actual_y = ypoints


    prediction_x = timeStamps[1:upper_limit+2] # include the next time series
    prediction_y = np.concatenate((model_fit.fittedvalues, forecast))
    
    x_axis_pred = []
    y_axis_pred = []
    x_axis_actual = []
    y_axis_actual = []

    addIncrease = None
    if process == 'consm':
        addIncrease = False
    elif process == 'prod':
        addIncrease = True

    # prediction values
    for i in range(len(prediction_x)):
        x_axis_pred.append(str(prediction_x[i]))
        # convert from W/m^2 to kWh
        energy = unitConversion(prediction_y[i], addIncrease)
        y_axis_pred.append(float(energy))
    
    # we do not want to bias the actual data!!
    for j in range(len(actual_x)):
        x_axis_actual.append(str(actual_x[j]))
        energy = unitConversion(actual_y[j], False)
        y_axis_actual.append(float(energy))
    

    # print("Predicted x-axis values:",x_axis_pred)
    # print("Predicted y-axis values:",y_axis_pred)
    # print()
    # print("Actual x-axis values:",x_axis_actual)
    # print("Actual y-axis values:",y_axis_actual)

    return [x_axis_actual, y_axis_actual, x_axis_pred, y_axis_pred]
    #Plot the graph
    # plt.xticks(rotation=30)
    # plt.plot(prediction_x, prediction_y, label='Energy Forecast', linestyle="-", color='red', marker="o")
    # plt.plot(actual_x, actual_y, label='Energy Data', linestyle="-", color='blue', marker="o") 
    # plt.legend()
    # plt.show()

# function to get the number of points in the data that needs to be taken based on time
def findTime():
    now = datetime.now()
    hrs = int(now.strftime('%H'))
    mint = int(now.strftime('%M'))
    upper_end = 0
    if not hrs == 0:
        if mint < 30:
            upper_end = hrs*2
        elif mint>=30 and mint<=45:
            upper_end = (hrs*2) + 1
        else:
            upper_end = (hrs*2) + 2
    else:
        upper_end = 4
    return upper_end

def getDateString():
    now = datetime.now()
    return str(now.strftime('%d-%m-%Y')) 

# function to convert Watts per square metre to kilowatt hours
def unitConversion(value, addRandomIncrease):
    increase = random.uniform(0, 0.5) # small number to increase overall demand
    avg_roof_top_size = 50 # 50 square metres worth of panels
    power_hrs = 0.5 # since we get for 30 min intervals
    energy_val = 0

    # for energy consumption, we cannot let the energy consumption fall to zero, so we add random increase
    if addRandomIncrease:
        energy_val = ((value * avg_roof_top_size * power_hrs)/100) + (  ((value * avg_roof_top_size * power_hrs)/100) *increase )
    else:
        energy_val = ((value * avg_roof_top_size * power_hrs)/100)
    

    # for consumption, the energy cannot fall to zero, so we add random increase
    if process == 'consm':    
        if energy_val <=50:
            generate_energy = random.uniform(80, 90)
            return generate_energy
        else:
            return energy_val

    # solar energy production can fall to zero, as there is no production at night
    elif process == 'prod':
        if energy_val <0:
            return 0
        else:
            return energy_val


def connectToDb():
    # load environment variables
    db_username = os.getenv("DB_USERNAME")
    db_password = os.getenv("DB_PASSWORD")
    db_cluster = os.getenv("DB_CLUSTER_ADDR")

    # connect to the database    
    client = pymongo.MongoClient('mongodb+srv://'+db_username+':'+ db_password +'@imdc-p2p-energy.y0a68.mongodb.net/'+db_cluster+'?retryWrites=true&w=majority')
    print("Connected to database")
    cluster=client["IMDC-p2p-energy"]
    collection = 0
    if process =='consm':
        collection = cluster.buyOrderForecast
    elif process =='prod':
        collection = cluster.energy_forecast
    # when matchmaking needs to be run
    else:
        print("Running matchmaking")
        initMatchmaking(client)
        print("Running double auction")
        initDoubleAuction(client)
        return

    # Get sample data
    forecast = energyForecast()
    actual_x = forecast[0]
    actual_y = forecast[1]
    pred_x = forecast[2]
    pred_y = forecast[3]
    dateString = getDateString()
    # data object
    data = {
            'actual_x' : actual_x,
            'actual_y' : actual_y,
            'pred_x' : pred_x,
            'pred_y': pred_y,
            'date': dateString,
            'userId':  sys.argv[1], 
            'current_pred': pred_y[len(pred_y)-1]
        }
    # if the document of this dateString exists for userId
    if collection.count_documents({"userId": sys.argv[1], "date":dateString}, limit=1)==1:
        result = collection.update_one({"userId": sys.argv[1], "date":dateString}, {"$set":data})
        print('Updated forecast data for', dateString)

    # insert new one if no document exists 
    else:
        # Insert business object directly into MongoDB via insert
        result=collection.insert_one(data)
        #Print to the console the ObjectID of the new document
        print('Added forecast data, id:', result.inserted_id)



connectToDb()


"""
Irradiance (watts per square metre) to kWh
(2 W/m^2 * 84 m^2 * 0.5hr) /1000

Solar Irradinace: watts per square metre
Malaysia has an average house size of about 1,264 sq ft 


Most residential solar panels on today's market are rated to produce 
between 250 and 400 watts each per hour.
take average then 325 per hr 

Domestic solar panel systems typically have a capacity of between 1 kW and 4 kW.


- When selling say this can produce 325 w per hr (predict it)
- When buying say based on capacitance 

"""