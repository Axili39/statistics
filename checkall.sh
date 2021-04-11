for TICKER in $(cat TICKERS.txt)
do
	if [ ! -f ${TICKER}.csv ]
	then
		echo "downloading ${TICKER}"
    		sleep 10
    		curl "https://query1.finance.yahoo.com/v7/finance/download/${TICKER}?period1=946857600&period2=1618012800&interval=1d&events=history&includeAdjustedClose=true" -o ${TICKER}.csv -s  -A "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:59.0) Gecko/20100101 Firefox/59.0"
	fi
    	./statistics -t ${TICKER}
done
