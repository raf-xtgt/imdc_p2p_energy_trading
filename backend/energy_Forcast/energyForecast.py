import numpy as np # linear algebra
import pandas as pd
import matplotlib.pyplot as plt
from statsmodels.tsa.api import ExponentialSmoothing, SimpleExpSmoothing, Holt



# Load the data
fields = ['PeriodStart', 'Ghi Curr Day'] # we only want to use these rows
filename = "Irradiance_39.xlsx"
df = pd.read_excel(filename, usecols=fields)
print(df.info) # print the dataframe

# select a specific segment of the data
data = df['Ghi Curr Day'][1:50]

# training the model
model = ExponentialSmoothing(data, trend="add", damped=True, seasonal="add", seasonal_periods=2)
model_fit = model.fit()
#print("Forecast plots:", model_fit.fittedvalues)
# get forecast/prediction
forecast = model_fit.predict()
#print("Actual values:", data)
print("Forecast values:", forecast)


# Plot the forecast 
actual_x = df['PeriodStart'][1:50].values
ypoints = df['Ghi Curr Day'][1:50].values
actual_y = ypoints


prediction_x = df['PeriodStart'][1:51].values # include the next time series
prediction_y = np.concatenate((model_fit.fittedvalues, forecast))
plt.xticks(rotation=30)
plt.plot(prediction_x, prediction_y, label='Energy Forecast', linestyle="-", color='red', marker="o")
plt.plot(actual_x, actual_y, label='Energy Data', linestyle="-", color='blue', marker="o") 
plt.legend()
plt.show()

