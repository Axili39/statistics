TICKER=ORA.PA
curl "https://query1.finance.yahoo.com/v7/finance/download/${TICKER}?period1=946857600&period2=1618012800&interval=1d&events=history&includeAdjustedClose=true" -o ${TICKER}.csv
