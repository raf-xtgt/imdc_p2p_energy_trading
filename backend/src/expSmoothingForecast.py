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

load_dotenv()
USERID = ""
# handle arguments passed to the script
# total arguments
n = len(sys.argv)
print("Total arguments passed:", n)
 
# Arguments passed
print("\nName of Python script:", sys.argv[0])
 
print("\nArguments passed:", sys.argv[1], "\n")


def energyForecast():
    # Load the data
    fields = ['PeriodStart', 'Ghi Curr Day'] # we only want to use these rows
    filename = "Irradiance_39.xlsx"
    df = pd.read_excel(filename, usecols=fields)
    #print(df.info) # print the dataframe

    # select a specific segment of the data
    data = df['Ghi Curr Day'][1:50]
    # training the model
    model = ExponentialSmoothing(data, trend="add", damped=True, seasonal="add", seasonal_periods=2)
    model_fit = model.fit() # points according to the model
    #print("Forecast plots:", model_fit.fittedvalues)
    
    # get forecast/prediction for the next point
    forecast = model_fit.predict()
    #print("Actual values:", data)
    #print("Forecast values:", forecast)

    # Plot the forecast 
    actual_x = df['PeriodStart'][1:50].values
    ypoints = df['Ghi Curr Day'][1:50].values
    actual_y = ypoints


    prediction_x = df['PeriodStart'][1:51].values # include the next time series
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