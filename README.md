# RecurExif
A simple CLI tool to find all files containing metadata.

## Installation
1. Verify you have Go 1.25 or higher
2. `go install github.com/varppi/RecurExif/cmd/RecurExif`
3. `export PATH=~/go/bin`
4. `RecurExif`

## Help
```
--extensions string   Only parses files with specific extensions. Example: pdf,png,jpg
--json                Outputs everything in JSON
--no-progress         Disables progress bar
--quiet               Doesn't show any extra info
--search string       Finds all files that contain metadata containing the search keyword
```

## Example output
```
______________________________  ______________________  __________________
___  __ \__  ____/_  ____/_  / / /__  __ \__  ____/_  |/ /___  _/__  ____/
__  /_/ /_  __/  _  /    _  / / /__  /_/ /_  __/  __    / __  / __  /_    
_  _, _/_  /___  / /___  / /_/ / _  _, _/_  /___  _    | __/ /  _  __/    
/_/ |_| /_____/  \____/  \____/  /_/ |_| /_____/  /_/|_| /___/  /_/       


starting...

picture.jpg:
EXIF: {"DateTime":"2025:02:01 9:32:37","DateTimeDigitized":"2025:02:01 9:32:37\u0000","DateTimeOriginal":"2025:02:01 9:32:37\u0000","ExifTag":[134],"ExifVersion":"0220","ExposureTime":[{"Numerator":30126919,"Denominator":1000000000}],"FNumber":[{"Numerator":16400,"Denominator":10000}],"Flash":[0],"FocalLength":[{"Numerator":44900,"Denominator":10000}],"ISOSpeedRatings":[425],"ImageLength":[3025],"ImageWidth":[4301],"LightSource":[0],"Make":"Oneplus","Model":"10","Orientation":[1],"PixelXDimension":[4051],"PixelYDimension":[3054],"SubSecTime":"255","SubSecTimeDigitized":"127216","SubSecTimeOriginal":"125726","WhiteBalance":[0],"no_tag_name_10":"+06:00","no_tag_name_19":"+06:00","no_tag_name_20":"+06:00"}

all files processed!
```
