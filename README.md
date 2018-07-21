# vboxanalyser  [![CircleCI](https://circleci.com/gh/ocrease/vboxanalyser.svg?style=svg)](https://circleci.com/gh/ocrease/vboxanalyser)
A tool to perform bulk analysis of Racelogic VBOX files

## Usage
Download the zip file in the latest release and unzip into a folder. 

Run vboxanalyser.exe and specify the options below.

## Options 
There are three command line parameters that can be passed to the program.

`-dir` Specify the directory to scan. Defaults to the current directory

`-c` Specify the data channel to scan. Choose from rpm, speedKph, speedMph. Defaults to rpm

`-t` Specify the threshold. Only files with a value higher than the threshold will be printed. Defaults to 8300 (for RPM) 

Example:

`vboxanalyser.exe -dir C:\Racing -c speedKph -t 190`



