avgTimeSeries
=============

This script will average the data column (second column, column 1 zero based) of MCell 
reaction data output files. The script reads any number of two column reaction data 
files and average the data column across all of them, effectively averaging across 
the time series. If a directory is provided instead of a reaction data file name, 
the script will consider all reaction data files contained within but not recurse 
beyond that. 

(C) Markus Dittrich, 2015
