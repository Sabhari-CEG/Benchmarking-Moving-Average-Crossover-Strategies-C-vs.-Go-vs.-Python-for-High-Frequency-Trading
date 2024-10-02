import time
import pandas as pd
import numpy as np

def calculate_sma_manual(closes, period):
    sma = np.empty(len(closes))
    sma[:] = np.nan
    cumsum = np.cumsum(closes)
    sma[period-1:] = (cumsum[period-1:] - np.concatenate(([0], cumsum[:-period]))) / period
    return sma

def moving_average_crossover_manual(closes, short_period, long_period):
    short_sma = calculate_sma_manual(closes, short_period)
    long_sma = calculate_sma_manual(closes, long_period)

    buy_signals = []
    sell_signals = []

    start = max(short_period, long_period)
    short_above = short_sma[start-1:] > long_sma[start-1:]
    crossovers = short_above[1:] != short_above[:-1]
    buy_indices = np.where(crossovers & short_above[1:])[0] + start
    sell_indices = np.where(crossovers & ~short_above[1:])[0] + start

    return buy_indices, sell_indices

start_time = time.time()

# Load only the necessary columns
data = pd.read_csv('data.csv', usecols=['timestamp', 'close'])

# Convert 'close' to numpy array for faster operations
closes = data['close'].to_numpy()

short_period = 5
long_period = 10
buy_signals, sell_signals = moving_average_crossover_manual(closes, short_period, long_period)

end_time = time.time()

print("Buy Signals:")
for signal in buy_signals:
    print(f"Buy at {data['timestamp'][signal]}, Price: {data['close'][signal]}")

print("\nSell Signals:")
for signal in sell_signals:
    print(f"Sell at {data['timestamp'][signal]}, Price: {data['close'][signal]}")

print(f"Execution Time: {end_time - start_time:.6f} seconds")
