#include <iostream>
#include <vector>
#include <fstream>
#include <sstream>
#include <string>
#include <chrono>
#include <iomanip>
#include <algorithm>
#include <numeric>

struct Candle {
    std::string timestamp;
    double close;
};

void calculateSMA(const std::vector<double>& closes, int period, std::vector<double>& sma) {
    sma.resize(closes.size());
    if (closes.size() < period) return;

    double sum = std::accumulate(closes.begin(), closes.begin() + period, 0.0);
    sma[period - 1] = sum / period;

    for (size_t i = period; i < closes.size(); ++i) {
        sum += closes[i] - closes[i - period];
        sma[i] = sum / period;
    }
}

void movingAverageCrossover(const std::vector<double>& closes, int shortPeriod, int longPeriod, 
                            std::vector<size_t>& buySignals, std::vector<size_t>& sellSignals) {
    std::vector<double> shortSMA, longSMA;
    calculateSMA(closes, shortPeriod, shortSMA);
    calculateSMA(closes, longPeriod, longSMA);

    size_t start = std::max(shortPeriod, longPeriod);
    for (size_t i = start; i < closes.size(); ++i) {
        if (shortSMA[i] > longSMA[i] && shortSMA[i - 1] <= longSMA[i - 1]) {
            buySignals.push_back(i);
        } else if (shortSMA[i] < longSMA[i] && shortSMA[i - 1] >= longSMA[i - 1]) {
            sellSignals.push_back(i);
        }
    }
}

std::vector<Candle> readCSV(const std::string& filename) {
    std::vector<Candle> candles;
    std::ifstream file(filename);
    std::string line;

    candles.reserve(200000);  

    std::getline(file, line);

    while (std::getline(file, line)) {
        size_t pos1 = line.find(',');
        size_t pos2 = line.find_last_of(',');
        if (pos1 != std::string::npos && pos2 != std::string::npos) {
            candles.push_back({line.substr(0, pos1), std::stod(line.substr(pos2 + 1))});
        }
    }

    return candles;
}

int main() {
    auto start = std::chrono::high_resolution_clock::now();

    std::vector<Candle> candles = readCSV("data.csv");
    std::vector<double> closes;
    closes.reserve(candles.size());
    for (const auto& candle : candles) {
        closes.push_back(candle.close);
    }

    int shortPeriod = 5;
    int longPeriod = 10;
    std::vector<size_t> buySignals, sellSignals;
    movingAverageCrossover(closes, shortPeriod, longPeriod, buySignals, sellSignals);

    std::cout << "Buy Signals:\n";
    for (size_t signal : buySignals) {
        std::cout << "Buy at " << candles[signal].timestamp << ", Price: " << std::fixed << std::setprecision(2) << candles[signal].close << "\n";
    }

    std::cout << "\nSell Signals:\n";
    for (size_t signal : sellSignals) {
        std::cout << "Sell at " << candles[signal].timestamp << ", Price: " << std::fixed << std::setprecision(2) << candles[signal].close << "\n";
    }

    auto end = std::chrono::high_resolution_clock::now();
    std::chrono::duration<double> diff = end - start;
    std::cout << "Execution Time: " << diff.count() << " seconds\n";

    return 0;
}
