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

load_dotenv()
USERID = ""
# handle arguments passed to the script
# total arguments
n = len(sys.argv)
print("Total arguments passed:", n)
 
# Arguments passed
print("\nName of Python script:", sys.argv[0])
 
print("\nArguments passed:", sys.argv[1], "\n")
print("\nArguments passed:", sys.argv[2], "\n")

timeStamps = ["12:00AM", "12:30AM", "01:00AM","01:30AM","02:00AM","02:30AM", "03:00AM", "03:30AM",
"04:00AM", "04:30AM", "05:00AM", "05:30AM", "06:00AM", "06:30AM","07:00AM", "07:30AM",
"08:00AM", "08:30AM", "09:00AM", "09:30AM", "10:00AM", "10:30AM", "11:00AM", "11:30AM", 
"12:00PM","12:30PM", "01:00PM","01:30PM","02:00PM","02:30PM", "03:00PM", "03:30PM",
"04:00PM", "04:30PM", "05:00PM", "05:30PM", "06:00PM", "06:30PM","07:00PM", "07:30PM",
"08:00PM", "08:30PM", "09:00PM", "09:30PM", "10:00PM", "10:30PM", "11:00PM", "11:30PM",]

def energyForecast():
    upper_limit = findTime()
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

    # prediction values
    for i in range(len(prediction_x)):
        x_axis_pred.append(prediction_x[i])
        y_axis_pred.append(prediction_y[i])
    

    for j in range(len(actual_x)):
        x_axis_actual.append(actual_x[j])
        y_axis_actual.append(actual_y[j])
    

    print("Predicted x-axis values:",x_axis_pred)
    print("Predicted y-axis values:",y_axis_pred)
    print()
    print("Predicted x-axis values:",x_axis_actual)
    print("Predicted y-axis values:",y_axis_actual)


    plt.xticks(rotation=30)
    plt.plot(prediction_x, prediction_y, label='Energy Forecast', linestyle="-", color='red', marker="o")
    plt.plot(actual_x, actual_y, label='Energy Data', linestyle="-", color='blue', marker="o") 
    plt.legend()
    plt.show()

# function to get the number of points in the data that needs to be taken based on time
def findTime():
    now = datetime.now()
    hrs = int(now.strftime('%H'))
    mint = int(now.strftime('%M'))
    upper_end = 0
    if mint < 30:
        upper_end = hrs*2
    elif mint>=30 and mint<=45:
        upper_end = (hrs*2) + 1
    else:
        upper_end = (hrs*2) + 2
    return upper_end

def getDateString():
    now = datetime.now()
    return str(now.strftime('%d/%m/%Y')) 

def connectToDb():
    # load environment variables
    db_username = os.getenv("DB_USERNAME")
    db_password = os.getenv("DB_PASSWORD")
    db_cluster = os.getenv("DB_CLUSTER_ADDR")

    # connect to the database    
    client = pymongo.MongoClient('mongodb+srv://'+db_username+':'+ db_password +'@imdc-p2p-energy.y0a68.mongodb.net/'+db_cluster+'?retryWrites=true&w=majority')
    print("Connected to database")


#connectToDb()
energyForecast()


"""
1. Connect the python script to db.
2. Find a way to call the python script from golang and pass in the user id as argument
3. Create a forecast dataset for each user. [Data and user id]
"""