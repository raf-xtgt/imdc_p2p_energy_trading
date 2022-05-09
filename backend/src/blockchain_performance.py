import matplotlib.pyplot as plt


def plotGraph():
    transactions = [1, 5, 25, 50, 100, 500, 1000, 2500, 5000] # number of transactions in the block
    #mining_time = [0, 4.484193152, 10.49449102,  35.7907982, 2.898173851, 26.90469269, 75.00000099, 12.85793697, 7.332277854 ] # time taken to mine each block with the corresponding number of transactions
    mining_time = [0, 4.484, 10.494,  35.791, 2.898, 26.905, 75.000, 12.858, 7.332] # time taken to mine each block with the corresponding number of transactions
    #Plot the graph
    plt.xticks(rotation=30)
    plt.grid(True)
    plt.title("Change in Mining Time vs Total Transactions per Block", fontsize=14)
    plt.plot(transactions, mining_time, label="Mining Time per Block", linestyle="-", color='#D9381E', marker="o")
    plt.xlabel("Number of Transactions per Block", fontsize=14)
    plt.ylabel("Mining Time (seconds)", fontsize=14)
    plt.legend()
    plt.show()

print(plotGraph())




def trnSizeVsMiningTime():
    transaction = [1, 5, 25, 50, 100, 500, 1000, 2500, 5000] # number of transactions in the block
    transaction_size = [767, 1100, 19200, 38400, 77900, 389400, 778800, 1900000, 3900000] # transaction size in bytes
    
    #mining_time = [0, 4.484193152, 10.49449102,  35.7907982, 2.898173851, 26.90469269, 75.00000099, 12.85793697, 7.332277854 ] # time taken to mine each block with the corresponding number of transactions
    mining_time = [0, 4.484, 10.494,  35.791, 2.898, 26.905, 75.000, 12.858, 7.332] # time taken to mine each block with the corresponding number of transactions
    #Plot the graph
    plt.xticks(rotation=30)
    plt.grid(True)
    plt.title("Change in Mining Time vs Block Size", fontsize=14)
    plt.plot(transaction_size, mining_time, label="Mining Time per Block", linestyle="-", color='#D9381E', marker="o")
    plt.xlabel("Transaction Size per Block (megabytes)", fontsize=14)
    plt.ylabel("Mining Time (seconds)", fontsize=14)
    plt.legend()
    plt.show()

trnSizeVsMiningTime()


"""
Index:  1 No. of Transactions: 2
Index:  2 No. of Transactions: 1
Index:  3 No. of Transactions: 2
Index:  4 No. of Transactions: 1
Index:  5 No. of Transactions: 1
Index:  6 No. of Transactions: 1
Index:  7 No. of Transactions: 1
Index:  8 No. of Transactions: 1
Index:  9 No. of Transactions: 1
Index:  10 No. of Transactions: 1
Index:  11 No. of Transactions: 118
Index:  12 No. of Transactions: 124
Index:  13 No. of Transactions: 107
Index:  14 No. of Transactions: 157
Index:  15 No. of Transactions: 498
Index:  16 No. of Transactions: 1070
Index:  17 No. of Transactions: 867
Index:  18 No. of Transactions: 969


"""