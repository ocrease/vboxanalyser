# vboxanalyser  [![CircleCI](https://circleci.com/gh/ocrease/vboxanalyser.svg?style=svg)](https://circleci.com/gh/ocrease/vboxanalyser)
A tool to perform bulk analysis of Racelogic VBOX files

## Usage
Download the zip file in the latest release and unzip into a folder. 

Run vboxanalyser.exe -dir "Path to vbox files" -t "threshold". Eg:

`vboxanalyser.exe -dir C:\Racing -t 8400`

This will analyse all .vbo files under the directory and will list any files where the RPM reached more than 8400.