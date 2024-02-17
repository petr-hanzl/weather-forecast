# weather-forecast

Weather forecast is an small API used for getting temperature and precipitation for given city name.

## How to use
Make an API call `localhost:3000//api/{city}` where city is name of the city.
Your request should look like this:
```
{
  "temperature": true, // if you want to get temperature
  "precipitation: true, // if you want to get precipitations
  "days": 5 // if you want to specifies number of days - default is 1
}
```
