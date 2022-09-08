/*
	Name: 		Byron Dowling
	Course:		5443 Spatial Databases
	Semester:	Fall, 2022	
	Assignment:	Program 1 - Simple API using Go Gin and PostgreSQL

	References used:

		- https://pkg.go.dev/database/sql (SQL)
		- https://pkg.go.dev/net/http (HTTP Responses)
		- https://pkg.go.dev/github.com/gin-gonic/gin#section-readme (Gin API)
		- https://go.dev/doc/tutorial/web-service-gin (Web Service Gin)
		- https://pkg.go.dev/strconv (String Convert)
		- https://github.com/umahmood/haversine (Haversine Formula)
		- https://dev.to/ramu_mangalarapu/building-rest-apis-in-golang-go-gin-with-persistence-database-postgres-4616


	Assignment Instructions:

		- Install Postgres DB and PostGIS
		- Install pgAdmin4
		- Created DB called Project1 using Public schema
		- Load a data file from the MSU CS server
		- Have the following API routes
			- findAll
			- findOne
			- findClosest


	Frameworks used:
		- Go
		- Gin
		- PostgreSQL		
*/

package main

// Libraries Needed
import(
	"fmt"
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/umahmood/haversine"
	"strconv"
	_ "github.com/lib/pq"
)

// PostgreSQL Database login info
const(
	host = "localhost"
	port = 5432
	user = "postgres"
	password = "PorkandBeans"
	dbname = "Project1"
)


/*

	██╗   ██╗███████╗███╗   ██╗██╗   ██╗
	██║   ██║██╔════╝████╗  ██║██║   ██║
	██║   ██║█████╗  ██╔██╗ ██║██║   ██║
	╚██╗ ██╔╝██╔══╝  ██║╚██╗██║╚██╗ ██╔╝
	 ╚████╔╝ ███████╗██║ ╚████║ ╚████╔╝ 
	  ╚═══╝  ╚══════╝╚═╝  ╚═══╝  ╚═══╝  
                                    

	Struct that holds a SQL Database Pointer

	Used in lieu of a global variable for our Database so that when we reference our 
	PostgreSQL DB in our API routes, our DB does not go out of scope.

	Found this workaround in the references below:

	https://pkg.go.dev/database/sql
	https://dev.to/ramu_mangalarapu/building-rest-apis-in-golang-go-gin-with-persistence-database-postgres-4616
*/

type vEnv struct {
	DB *sql.DB
}


/*

	██╗   ██╗███████╗ ██████╗ 
	██║   ██║██╔════╝██╔═══██╗
	██║   ██║█████╗  ██║   ██║
	██║   ██║██╔══╝  ██║   ██║
	╚██████╔╝██║     ╚██████╔╝
	╚═════╝ ╚═╝      ╚═════╝ 
							
	Structs that holdsall the info of a UFO sighting that we will be querying from
	our database and return from our Gin API.

	Ex:) 
		
		{
			"DateTime": "12/21/16 19:15",
			"Country": "USA",
			"City": "Waynesboro",
			"State": "VA",
			"Shape": "Sphere",
			"Lat": 38.065228,
			"Long": -78.90588
		}

		{
			"Message": "The closest ufo sighting is 0.690783 miles away",
			"Value": 
			{
				"DateTime": "3/15/16 06:14",
				"Country": "USA",
				"City": "Wichita Falls",
				"State": "TX",
				"Shape": "Unknown",
				"Lat": 33.913708,
				"Long": -98.493385
			}
		}
*/

type UFO struct {
	DateTime string
	Country string
	City string
	State string
	Shape string
	Lat float32
	Long float32
}

type UFO_Distance struct {
	DateTime string
	Country string
	City string
	State string
	Shape string
	Lat float32
	Long float32
	distanceMiles float64
	distanceKM float64
}



// Main Function
func main() {

	venv := new(vEnv)	// New SQL DB pointer object to keep DB within scope

	// PostgreSQL Connection string
	postgresCX := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Connecting to PostgreSQL
	var err error
	venv.DB, err = sql.Open("postgres", postgresCX)
	checkError(err)


	// Sanity Check that DB Connection was successful
	err = venv.DB.Ping()
	checkError(err)
	fmt.Println("Database Connected!")
	fmt.Println()


	/*
		 █████╗ ██████╗ ██╗    ██╗███╗   ██╗███████╗ ██████╗                             
		██╔══██╗██╔══██╗██║    ██║████╗  ██║██╔════╝██╔═══██╗                            
		███████║██████╔╝██║    ██║██╔██╗ ██║█████╗  ██║   ██║                            
		██╔══██║██╔═══╝ ██║    ██║██║╚██╗██║██╔══╝  ██║   ██║                            
		██║  ██║██║     ██║    ██║██║ ╚████║██║     ╚██████╔╝                            
		╚═╝  ╚═╝╚═╝     ╚═╝    ╚═╝╚═╝  ╚═══╝╚═╝      ╚═════╝                             
																						
		 █████╗ ███╗   ██╗██████╗     ██████╗  ██████╗ ██╗   ██╗████████╗███████╗███████╗
		██╔══██╗████╗  ██║██╔══██╗    ██╔══██╗██╔═══██╗██║   ██║╚══██╔══╝██╔════╝██╔════╝
		███████║██╔██╗ ██║██║  ██║    ██████╔╝██║   ██║██║   ██║   ██║   █████╗  ███████╗
		██╔══██║██║╚██╗██║██║  ██║    ██╔══██╗██║   ██║██║   ██║   ██║   ██╔══╝  ╚════██║
		██║  ██║██║ ╚████║██████╔╝    ██║  ██║╚██████╔╝╚██████╔╝   ██║   ███████╗███████║
		╚═╝  ╚═╝╚═╝  ╚═══╝╚═════╝     ╚═╝  ╚═╝ ╚═════╝  ╚═════╝    ╚═╝   ╚══════╝╚══════╝
																							
																						

		app.GET("/Home/All", venv.FindAll)
			- Route will return all ufo sighting records

		app.GET("/Home/Closest/:lat/:long", venv.NearestNeighbor)
			- Route accepts a Lat/Long and returns nearest sighting

		app.GET("/Home/FindSightingByCity/:id", venv.FindSightingByCity)
			- Returns all matches by City

		app.GET("/Home/FindSightingByState/:id", venv.FindSightingByState)
			- Returns all matches by State

		app.GET("/Home/FindSightingByCountry/:id", venv.FindSightingByCountry)
			- Returns all matches by Country
	*/

	// Declaring an instance of Gin
	app := gin.Default()

	// Home page info
	app.GET("/Home", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Message": `Hello and thanks for stopping by, this API is brought to you by Squarespace.`,
			"DISCLAIMER": "Please no Solicitors",
		})
	})

	// Additional routes
	app.GET("/Home/All", venv.FindAll)
	app.GET("/Home/Closest/:lat/:long", venv.NearestNeighbor)
	app.GET("/Home/FindSightingByCity/:id", venv.FindSightingByCity)
	app.GET("/Home/FindSightingByState/:id", venv.FindSightingByState)
	app.GET("/Home/FindSightingByCountry/:id", venv.FindSightingByCountry)
	
	app.Run() 		// Launch and run our Gin API


	venv.DB.Close()	// Close our Database
}


/*
    Public API Function void: NearestNeighbor()

    Description:
        Function that determines the closest ufo sighting to the coordinates
		that are passed into the URL.

        This is performed by converting the string input into float64 values
		and plugging them into the haversine distance formula. This is checked
		against all results in the database and the one with the smallest distance
		is returned as a JSON string on the API page.

    Params:
        - c *gin.Context
		- Lat
		- Long

    Returns:
        - c.JSON
*/

func (env vEnv)NearestNeighbor(c *gin.Context) {

	// Values that will be read from PostgreSQL
	var date string
	var country string
	var city string
	var state string
	var shape string
	var lat float32
	var long float32

	// Variables needed for GPS
	var DM float64
	var DKM float64


	// Need to convert from float32 to float64 since float64 is required for Haversine
	lattitude, e := strconv.ParseFloat(c.Param("lat"), 64)
	checkError(e)

	// Need to convert from float32 to float64 since float64 is required for Haversine
	longitude, er := strconv.ParseFloat(c.Param("long"), 64)
	checkError(er)

	// Coordinate Pair that is passed in
	cp1 := haversine.Coord{Lat: lattitude, Lon: longitude}

	result, err := env.DB.Query(`SELECT * FROM public.ufo_sightings`)
	checkError(err)

	// Arbitratily large value that will reduce as we find the nearest sighting
	var closestDistance float64 = 10000.123456
	var closestSighting UFO_Distance

	// Looping through matches
	for result.Next() {
		result.Scan(&date, &country, &city, &state, &shape, &lat, &long)

		// Casting to float64
		L1 := float64(lat)
		L2 := float64(long)
		
		// Haversine coord pair of the current entry
		cp2 := haversine.Coord{Lat: L1, Lon: L2}

		mi, km := haversine.Distance(cp1, cp2)
		DM = mi
		DKM = km

		if mi < closestDistance {

			closestDistance = mi
			temp := UFO_Distance{DateTime: date, Country: country, City: city, State: state,
				Shape: shape, Lat: lat, Long: long, distanceMiles: DM, distanceKM: DKM }
			
			closestSighting = temp
		}

	}

	/*
		Coordinate Pairs to Test:

			Wichita Falls TX
				33.913708/-98.493385

			Las Cruces NM
				32.3140354/-106.779807

			Area 51/Edwards AFB CA/NV
				34.900814/-117.9439533

	*/

	output := fmt.Sprintf("The closest ufo sighting is %f miles away", closestDistance)
	
	c.JSON(200, gin.H{
		"Message": output,
		"Value": closestSighting,
	})
}


/*
    Public API Function void: FindSightingByCountry()

    Description:
        Function that returns a list of ufo sightings that match the country
		passed in. This was predominantly added as there are sighting in both
		the United States and Canada in this list.

    Params:
        - c *gin.Context
		- Country

    Returns:
        - c.JSON
*/

func (env vEnv)FindSightingByCountry(c *gin.Context) {

	// List of ufo structs
	var ufos []UFO

	// Values that will be read from PostgreSQL
	var date string
	var country string
	var city string
	var state string
	var shape string
	var lat float32
	var long float32

	id := c.Param("id")

	query := (`SELECT * FROM public.ufo_sightings WHERE country=$1`)
	result, e := env.DB.Query(query, id)
	checkError(e)

	for result.Next() {
		result.Scan(&date, &country, &city, &state, &shape, &lat, &long)

		ufos = append(ufos, UFO{DateTime: date, Country: country, City: city, State: state,
			Shape: shape, Lat: lat, Long: long})
	}

	c.JSON(http.StatusOK, ufos)
}


/*
    Public API Function void: FindSightingByState()

    Description:
        Function that returns a list of ufo sightings that match the state
		passed in. This allows to see all the sightings that have occured
		in your state and inspect the details of each.

    Params:
        - c *gin.Context
		- State

    Returns:
        - c.JSON
*/

func (env vEnv)FindSightingByState(c *gin.Context) {

	// List of ufo structs
	var ufos []UFO

	// Values that will be read from PostgreSQL
	var date string
	var country string
	var city string
	var state string
	var shape string
	var lat float32
	var long float32

	id := c.Param("id")

	query := (`SELECT * FROM public.ufo_sightings WHERE state=$1`)
	result, e := env.DB.Query(query, id)
	checkError(e)

	for result.Next() {
		result.Scan(&date, &country, &city, &state, &shape, &lat, &long)

		ufos = append(ufos, UFO{DateTime: date, Country: country, City: city, State: state,
			Shape: shape, Lat: lat, Long: long})
	}

	c.JSON(http.StatusOK, ufos)
}


/*
    Public API Function void: FindSightingByCity()

    Description:
        Function that returns a list of ufo sightings that match the city
		passed in. This allows to see all the sightings that have occured
		in your city and inspect the details of each.

    Params:
        - c *gin.Context
		- City

    Returns:
        - c.JSON
*/

func (env vEnv)FindSightingByCity(c *gin.Context) {

	// List of ufo structs
	var ufos []UFO

	// Values that will be read from PostgreSQL
	var date string
	var country string
	var city string
	var state string
	var shape string
	var lat float32
	var long float32

	id := c.Param("id")

	query := (`SELECT * FROM public.ufo_sightings WHERE city=$1`)
	result, e := env.DB.Query(query, id)
	checkError(e)

	for result.Next() {
		result.Scan(&date, &country, &city, &state, &shape, &lat, &long)

		ufos = append(ufos, UFO{DateTime: date, Country: country, City: city, State: state,
			Shape: shape, Lat: lat, Long: long})
	}

	c.JSON(http.StatusOK, ufos)
}


/*
    Public API Function void: FindAll()

    Description:
        Default function that will return all results of the database and
		display them as a JSON array.

    Params:
        - c *gin.Context

    Returns:
        - c.JSON
*/

func (env vEnv)FindAll(c *gin.Context) {

	rows, err := env.DB.Query(`SELECT * FROM public.ufo_sightings`)
	checkError(err)

	// List of ufo structs
	var ufos []UFO

	// Values that will be read from PostgreSQL
	var date string
	var country string
	var city string
	var state string
	var shape string
	var lat float32
	var long float32

	for rows.Next() {
		rows.Scan(&date, &country, &city, &state, &shape, &lat, &long)

		ufos = append(ufos, UFO{DateTime: date, Country: country, City: city, State: state,
			Shape: shape, Lat: lat, Long: long})
	}


	c.JSON(http.StatusOK, ufos)
}


/*
	 ██████╗██╗  ██╗███████╗ ██████╗██╗  ██╗ 
	██╔════╝██║  ██║██╔════╝██╔════╝██║ ██╔╝ 
	██║     ███████║█████╗  ██║     █████╔╝  
	██║     ██╔══██║██╔══╝  ██║     ██╔═██╗  
	╚██████╗██║  ██║███████╗╚██████╗██║  ██╗ 
	 ╚═════╝╚═╝  ╚═╝╚══════╝ ╚═════╝╚═╝  ╚═╝ 
											
	███████╗██████╗ ██████╗  ██████╗ ██████╗ 
	██╔════╝██╔══██╗██╔══██╗██╔═══██╗██╔══██╗
	█████╗  ██████╔╝██████╔╝██║   ██║██████╔╝
	██╔══╝  ██╔══██╗██╔══██╗██║   ██║██╔══██╗
	███████╗██║  ██║██║  ██║╚██████╔╝██║  ██║
	╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝
										
	Utility function that errors are sent to check if an error has
	occured in an operation. Many built in functions in Go require you
	to declare two variable, the data variable and an error variable.

	Instead of checking all of them inline, the errors are sent here to
	be checked and handled.

*/
func checkError(err error) {
	if err != nil {
        panic(err)
    }
}
