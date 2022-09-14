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
	// "github.com/umahmood/haversine"
	// "strconv"
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
							
	Structs that holds all the info of a UFO sighting that we will be querying from
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
	Geo string
	dist float64
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

        This is performed by converting the PostGreSQL ST_Distance query

    Params:
        - c *gin.Context
		- Lat
		- Long

    Returns:
        - c.JSON
*/

func (env vEnv)NearestNeighbor(c *gin.Context) {

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
	var geo string
	var distance float64

	/*
		Coordinate Pairs to Test:

			Wichita Falls TX
				33.913708/-98.493385

			Las Cruces NM
				32.3140354/-106.779807

			Area 51/Edwards AFB CA/NV
				34.900814/-117.9439533

	*/


	sqlStatemnt := fmt.Sprintf(`SELECT *, ST_Distance('SRID=4326;POINT(%s %s)'::geometry, location)
	AS dist FROM public.ufo_sightings ORDER BY dist LIMIT 10`, c.Param("long"), c.Param("lat"))

	fmt.Println(sqlStatemnt)

	result, err := env.DB.Query(sqlStatemnt)
	checkError(err)

	for result.Next() {
		result.Scan(&date, &country, &city, &state, &shape, &lat, &long, &geo, &distance)

		ufos = append(ufos, UFO{DateTime: date, Country: country, City: city, State: state,
			Shape: shape, Lat: lat, Long: long, Geo: geo, dist: distance})
	}

	c.JSON(http.StatusOK, ufos)
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
	var geo string

	id := c.Param("id")

	query := (`SELECT * FROM public.ufo_sightings WHERE country=$1`)
	result, e := env.DB.Query(query, id)
	checkError(e)

	for result.Next() {
		result.Scan(&date, &country, &city, &state, &shape, &lat, &long, &geo)

		ufos = append(ufos, UFO{DateTime: date, Country: country, City: city, State: state,
			Shape: shape, Lat: lat, Long: long, Geo: geo})
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
	var geo string

	id := c.Param("id")

	query := (`SELECT * FROM public.ufo_sightings WHERE state=$1`)
	result, e := env.DB.Query(query, id)
	checkError(e)

	for result.Next() {
		result.Scan(&date, &country, &city, &state, &shape, &lat, &long, &geo)

		ufos = append(ufos, UFO{DateTime: date, Country: country, City: city, State: state,
			Shape: shape, Lat: lat, Long: long, Geo: geo})
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
	var geo string

	id := c.Param("id")

	query := (`SELECT * FROM public.ufo_sightings WHERE city=$1`)
	result, e := env.DB.Query(query, id)
	checkError(e)

	for result.Next() {
		result.Scan(&date, &country, &city, &state, &shape, &lat, &long, &geo)

		ufos = append(ufos, UFO{DateTime: date, Country: country, City: city, State: state,
			Shape: shape, Lat: lat, Long: long, Geo: geo})
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
	var geo string

	for rows.Next() {
		rows.Scan(&date, &country, &city, &state, &shape, &lat, &long, &geo)

		ufos = append(ufos, UFO{DateTime: date, Country: country, City: city, State: state,
			Shape: shape, Lat: lat, Long: long, Geo: geo})
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
