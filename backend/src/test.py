from datetime import datetime
from calendar import monthrange
now = datetime.now()
print(now.strftime('%Y/%m/%d %H:%M:%S')) #24-hour format

print("Time in 24hr format:", now.strftime('%H:%M')) #24-hour format


def findTime():
    now = datetime.now()
    hrs = int(now.strftime('%H'))
    mint = int(now.strftime('%M'))
    month = int(now.strftime('%m'))
    year = int(now.strftime('%Y'))
    upper_end = 0
    num_days = monthrange(year, month)[1] # number of days in current month
    print("Number of days in this month", num_days)
    if mint < 30:
        upper_end = hrs*2
    elif mint>=30 and mint<=45:
        upper_end = (hrs*2) + 1
    else:
        upper_end = (hrs*2) + 2
    return upper_end



print(findTime())


#print("12 hr format")
#print(now.strftime('%Y/%m/%d %I:%M:%S')) #12-hour format
"""
24hrs in 1 day
so 24 points in 1 day
the data we got is of 30 min intervals 
so we got 48 points (2 points per hour)

at 16:08, we take 16*2 = 32 points (prediction for time 16:30)
at 16:30 to 16:45, we take = (16*2)+1 = 33 points (prediction for time 17:00)
at 16:45 to 16:59, we take = (16*2)+2 = 34 points (prediction for time 17:30)


24hrs in 1 day
so 24 points in 1 day
assuming n days per month: n*24 = points. 


"""