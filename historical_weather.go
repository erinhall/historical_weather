package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)
//to hold data from csv
type WeatherDataCSV struct {
	Name   string
	Date   string
	Precip string
	snow   string
	Tmax  string
	Tmin   string
}
//hold data for daysofprecip json output
type AvgPrecipJSON struct{
  City string
  Days_of_precip float64
}
//hold data for maxtempdelta json
type TempChangeJSON struct{
  City string
  Date time.Time
  Temp_change float64
}

func main() {
 args := os.Args[1:]
  ProcessArgs(args)
}

//on limited time, this was my solution for allowing a cli call that was similar to the instructions
func ProcessArgs(Args []string){
var Isdaysofprecip bool = false
var Ismaxtempdelta bool = false
var functionArg string
var cityArg string
for i := 1; i < len(Args); i++ {
    //trim the brackets from arg because of desired formatting in instructions
    Args[i]=strings.Trim(Args[i], "[],")
}
 argLen := len(Args)
  //require function name and city as args
  if argLen < 2{
    fmt.Println("Please try again with function name and city")
    os.Exit(0)
  } else {
    cityArg = os.Args[2]
    functionArg = os.Args[1]
    if functionArg == "days_of_precip"{
      Isdaysofprecip=true
    } else if functionArg =="max_temp_delta"{
      Ismaxtempdelta=true
    } else{
      fmt.Println("Please try again with function name of days_of_precip or max_temp_delta")
       os.Exit(0)
    }
  }
  switch argLen{
  case 2:    //only city given as arg for either function
    if Isdaysofprecip{
      days_of_precip(cityArg)
    } else if Ismaxtempdelta {
      max_temp_delta(cityArg)
    }
  case 3: //city and year given for maxtempdelta only
     if Isdaysofprecip{
       fmt.Println("days_of_precip can only accect a city arg")
       os.Exit(0)
     }
    yearArg, err := strconv.Atoi(os.Args[3])
    max_temp_delta(cityArg, yearArg)
          if err != nil {
        fmt.Println(err)
        }
  case 4: //city, year, month given for maxtempdelta only
    if Isdaysofprecip{
      fmt.Println("days_of_precip can only            accect a city arg")
       os.Exit(0)
     }
    monthArg, err := strconv.Atoi(os.Args[4])
    yearArg, err := strconv.Atoi(os.Args[3])
    max_temp_delta(cityArg, yearArg, monthArg)
         if err != nil {
        fmt.Println(err)
        }
} 
}

//read csv into array
func ReadWeatherData(filename string) ([][]string, error) {
	fileContent, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer fileContent.Close()
	lines, err := csv.NewReader(fileContent).ReadAll()
	if err != nil {
		return [][]string{}, err
	}
	return lines, nil
}

//can be called from cli
func days_of_precip(city string) {
	csvWeatherData, err := ReadWeatherData("csv/WeatherData.csv")
	if err != nil {
		panic(err)
	}
	var AvgRainStruct AvgPrecipJSON
  AvgRainStruct.Days_of_precip=0.0
  //check user input
  city = strings.ToLower(city)
  switch city{
    case "bos":
     AvgRainStruct.City = "BOSTON"     
    case "mia":
     AvgRainStruct.City = "MIAMI"
    case "jnu":
     AvgRainStruct.City = "JUNEAU"
    default: 
      fmt.Println("Accepted values are BOS, MIA, or JNU")
      os.Exit(0)
  }

  //loop through csv data in array and count days with precip for desired city. divide by 10 to get average for the 10 year period
	for _, line := range csvWeatherData {
		data := WeatherDataCSV{
			Name:   line[1],
			Date:   line[5],
			Precip: line[9],
		}      
    if (strings.HasPrefix(data.Name, AvgRainStruct.City )) {
      if data.Precip != "0" || (len(data.snow) != 0 && data.snow != "0") {
        AvgRainStruct.Days_of_precip += 1     
        }
    }
	}
	  AvgRainStruct.Days_of_precip =   AvgRainStruct.Days_of_precip / 10
//turn the struct into json and print to screen
   j,err :=json.Marshal(AvgRainStruct)
     if err != nil {
        fmt.Println(err)
    }
  fmt.Println(string(j))
}
//can be called from cli.  date...int allows for zero or more ints to be added as parameters
func max_temp_delta (city string, date...int ){
  var month int 
  var year int
  switch len(date){
    case 1: //get and validate year if given
    if InBetween(date[0], 2010, 2019){
     year =date[0]
      } else {
     fmt.Println("Accepted years are 2010-2019")
      os.Exit(0)
      }
    case 2: //get and validate month if given
     if InBetween(date[0], 2010,2019) && InBetween(date[1],1,12) {
    year =date[0]
    month =date[1]  
      } else {
     fmt.Println("Accepted years are 2010-2019 and accepted months are 1-12")
      os.Exit(0)
      }
    default: //no month or year
    year=0
    month=0
  }

  csvWeatherData, err := ReadWeatherData("csv/WeatherData.csv")
	if err != nil {
		panic(err)
	}
	//initializa struct for JSON
  var TempChangeStruct TempChangeJSON
  //validate user input for city
  city = strings.ToLower(city)
  switch city{
    case "bos":
     TempChangeStruct.City = "BOSTON"     
    case "mia":
     TempChangeStruct.City = "MIAMI"
    case "jnu":
     TempChangeStruct.City = "JUNEAU"
    default: 
      fmt.Println("Accepted values are BOS, MIA, or JNU")
      os.Exit(0)
  }
  	for _, line := range csvWeatherData {
		data := WeatherDataCSV{
			Name:   line[1],
			Date:   line[5],
			Tmax:   line[13],
      Tmin:   line[14],
		}     
      //parsing going on to convert strings to int and int to date as needed
       monthStr := strconv.Itoa(month)
       yearStr := strconv.Itoa(year)
      
    if (strings.HasPrefix(data.Name, TempChangeStruct.City )) {
      if year !=0{        
        if month !=0{   
          //start month and year
          if(strings.HasSuffix(data.Date, yearStr))&&(strings.HasPrefix(data.Date,monthStr)) {
           
        strconv.ParseFloat(data.Tmax, 64)
        maxTemp, err := strconv.ParseFloat(data.Tmax, 8)
        minTemp, err := strconv.ParseFloat(data.Tmin, 8)
        if err != nil {
		      fmt.Println(err)
	      }
      //calculating temp change for month,year
       tempchangeValue := maxTemp-minTemp
        if tempchangeValue > TempChangeStruct.Temp_change{
          TempChangeStruct.Temp_change = tempchangeValue
                    TempChangeStruct.Date, err = time.Parse("1-2-2006", strings.Replace(data.Date, "/", "-", -1) )
        }
          }
        } else { //start only year
          if(strings.HasSuffix(data.Date, yearStr)){
             strconv.ParseFloat(data.Tmax, 64)
        maxTemp, err := strconv.ParseFloat(data.Tmax, 8)
        minTemp, err := strconv.ParseFloat(data.Tmin, 8)
        if err != nil {
		      fmt.Println(err)
	      }
            //calculating temp change for year
       tempchangeValue := maxTemp-minTemp
        if tempchangeValue > TempChangeStruct.Temp_change{
          TempChangeStruct.Temp_change = tempchangeValue
                    TempChangeStruct.Date, err = time.Parse("1-2-2006", strings.Replace(data.Date, "/", "-", -1) )
        }
          }
        }  
      }else{ //start no month or year only city
         strconv.ParseFloat(data.Tmax, 64)
        maxTemp, err := strconv.ParseFloat(data.Tmax, 8)
        minTemp, err := strconv.ParseFloat(data.Tmin, 8)
        if err != nil {
		      fmt.Println(err)
	      }
      //calculating temp change for all dates
       tempchangeValue := maxTemp-minTemp
        if tempchangeValue > TempChangeStruct.Temp_change{
          TempChangeStruct.Temp_change = tempchangeValue
          TempChangeStruct.Date, err = time.Parse("1-2-2006", strings.Replace(data.Date, "/", "-", -1) )
	if err != nil {
		panic(err)
	}
        }
      }
    }
	}
    j,err :=json.Marshal(TempChangeStruct)
     if err != nil {
        fmt.Println(err)
        return
    }
  fmt.Println(string(j))
	}

//function to make sure year, month in range
 func InBetween(i, min, max int) bool {
         if (i >= min) && (i <= max) {
                 return true
         } else {
                 return false
         }
 }