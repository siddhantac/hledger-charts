# hledger-charts

## Todos
1. Fix: investment CSV output gives the total in the last row which gets included in the pie chart (does not happen for expenses because the first column in last row is empty string. In investment it is `Total`)
1. add flags for:
    a. output-dir
    b. hledger journal file
    c. fileserver dir (reuse output dir?)
    d. date-range
2. create and use hledger-cli tool (from puffin)
