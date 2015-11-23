package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
  	"gopkg.in/mgo.v2/bson"
  	"io/ioutil"
  	"os"
  	"bytes"
)

type PostRequest struct {
	StartingFromLocationID string `json:"starting_from_location_id"`
	LocationIds []string `json:"location_ids"`
}

type PutReqResponse struct {
	ID  bson.ObjectId `bson:"_id,omitempty"`
	Status string `json:"status"`
	StartingFromLocationID string `json:"starting_from_location_id"`
	NextDestinationLocationID string `json:"next_destination_location_id”"`
	BestRouteLocationIds []string `json:"best_route_location_ids"`
	TotalUberCosts int `json:"total_uber_costs"`
	TotalUberDuration int `json:"total_uber_duration"`
	TotalUberDistance float64 `json:"total_uber_distance"`
	UberWaitimeETA float64 `json:"uber_wait_time_eta"`
}
type PostResponse struct {
	ID  bson.ObjectId `bson:"_id,omitempty"`
	Status string `json:"status"`
	StartingFromLocationID string `json:"starting_from_location_id"`
	BestRouteLocationIds []string `json:"best_route_location_ids"`
	TotalUberCosts int `json:"total_uber_costs"`
	TotalUberDuration int `json:"total_uber_duration"`
	TotalUberDistance float64 `json:"total_uber_distance"`
}

type UberResponse struct {
	Prices []struct {
		LocalizedDisplayName string `json:"localized_display_name"`
		HighEstimate int `json:"high_estimate"`
		Minimum int `json:"minimum"`
		Duration int `json:"duration"`
		Estimate string `json:"estimate"`
		Distance float64 `json:"distance"`
		DisplayName string `json:"display_name"`
		ProductID string `json:"product_id"`
		LowEstimate int `json:"low_estimate"`
		SurgeMultiplier float64 `json:"surge_multiplier"`
		CurrencyCode string `json:"currency_code"`
	} `json:"prices"`
}

type ResponseAddOp struct {
	Coordinate struct {
		Lat float64 
		Lng float64 
	} `json:"coordinate"`
}
type ReuestIdPostRes struct {
	RequestID string `json:"request_id"`
	Status string `json:"status"`
	Vehicle interface{} `json:"vehicle"`
	Driver interface{} `json:"driver"`
	Location interface{} `json:"location"`
	Eta int `json:"eta"`
	SurgeMultiplier interface{} `json:"surge_multiplier"`
}
var putCounter=0
var finalDestReached=0

func main() {
		router := mux.NewRouter()
		router.HandleFunc("/trips", handlePostTrips).Methods("POST")
		router.HandleFunc("/trips/{trip_id}", handleTrips).Methods("GET", "DELETE", "PUT")
		http.ListenAndServe(":8080", router)
	}

func handleTrips(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(req)
		tripId := vars["trip_id"]
		sess, err := mgo.Dial("mongodb://kalpana:pass@ds035844.mongolab.com:35844/cmpe273")
		if err != nil {
		    fmt.Printf("Can't connect to mongo, go error 1 %v\n", err)
		    os.Exit(1)
		}
		defer sess.Close()		 
		sess.SetSafe(&mgo.Safe{})
		collection := sess.DB("cmpe273").C("TripPlanner")
			result := PostResponse{}
			err = collection.Find(bson.M{"_id": bson.ObjectIdHex(tripId)}).Select(bson.M{}).One(&result)
			if err != nil {
				res.WriteHeader(http.StatusNotFound)
			    fmt.Printf("Trip location not found: %v\n", err)
			}
	  switch req.Method {
		case "GET":		
			outputRes,_:= json.Marshal(result)
			res.WriteHeader(http.StatusOK)
			fmt.Fprint(res,string(outputRes))
		case "PUT":
			bestrouteArr:=result.BestRouteLocationIds
			responseStruct := &PutReqResponse{
				ID: bson.ObjectIdHex(tripId),
				Status: "finished",
				StartingFromLocationID: result.StartingFromLocationID,
				NextDestinationLocationID: result.StartingFromLocationID,
				BestRouteLocationIds: result.BestRouteLocationIds,
				TotalUberCosts : result.TotalUberCosts,
				TotalUberDuration : result.TotalUberDuration,
				TotalUberDistance : result.TotalUberDistance,
				UberWaitimeETA: 3,
			}
		 if finalDestReached!=2{			
			access_token:=" eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlsicHJvZmlsZSIsInJlcXVlc3RfcmVjZWlwdCIsInJlcXVlc3QiLCJoaXN0b3J5X2xpdGUiXSwic3ViIjoiNTgzYTI3ZmUtMWEyNy00YmI0LTllZWUtZGNmMTEyMDA1ZTYwIiwiaXNzIjoidWJlci11czEiLCJqdGkiOiI3ZDNjMjkyZi1iNzZjLTRmMzUtYjNjYS1kY2YxNTg2NjQxOGQiLCJleHAiOjE0NTA1NzIyOTksImlhdCI6MTQ0Nzk4MDI5OSwidWFjdCI6IkkzemZFR0o4cDZBeWdiU0NtVVlCcjdsNTBycWhWdiIsIm5iZiI6MTQ0Nzk4MDIwOSwiYXVkIjoiVjVZQnQwV2tCa1M4TGgyOWNTWlA0ZXpkQW9ZYUxZU1AifQ.UzZ8PtDkbsmik7kxOx-1ZQWfhLEdFtPglIE7j578DM5H2Y-Gk3_QIaoaotg2PMBX_d-t3L_kGeZ3tTN5nLiB5PZCsny1PbbeSe9L48AGfMqLnAi_LakJhOm0ufQf5J5Servs8QINZvny0MKMm6wANy2rHIeBCkx0a5y9bAEVdkBgrrfkVIq9SkYPjyBnQxwCLzQUXUnCFpv2LUFiYofbfCBiXsxRh0LNafnDtrQHMQ0Dj4p58xjkb1LsidSCp_RBqj5xN0HOameQjuTdIadqYHJE4-ui3nQWyFCRBpm68NqEFhdFbBYL7hN69p1JWrWDy0PItesfFbYyQjNXg8flwQ"
			var request_id interface{}
			var request_eta interface{}
			var data interface{}
			var nextDestiId string
			var statusUpdate="requesting"

			var startLat,startLng,endLat,endLng float64
			product_id:=""

			if len(bestrouteArr)==putCounter{
				finalDestReached=1
			}	

			if finalDestReached!=1{
				if putCounter==0{
					startLat,startLng=coordinateValues(result.StartingFromLocationID)
					endLat,endLng=coordinateValues(bestrouteArr[0])
					nextDestiId=bestrouteArr[0]
				} else{
					startLat,startLng=coordinateValues(bestrouteArr[putCounter-1])
					endLat,endLng=coordinateValues(bestrouteArr[putCounter])
					nextDestiId=bestrouteArr[putCounter]
				}
				} else {
					startLat,startLng=coordinateValues(bestrouteArr[len(bestrouteArr)-1])
					endLat,endLng=coordinateValues(result.StartingFromLocationID)
					finalDestReached=2
					statusUpdate="finished"
					nextDestiId=result.StartingFromLocationID
				}
				url:= fmt.Sprintf("https://sandbox-api.uber.com/v1/estimates/price?start_latitude=%f&start_longitude=%f&end_latitude=%f&end_longitude=%f&server_token=OnEamv71AJVwxWiXz_8_aN5LzsAZ-R6E_FcRI8ug", startLat,startLng,endLat,endLng)
				var uberRes UberResponse
	 			responseValue, err := http.Get(url)
	            defer responseValue.Body.Close()
	            reply, err := ioutil.ReadAll(responseValue.Body)
	            json.Unmarshal([]byte(reply), &uberRes)
	            if err != nil {
	                fmt.Printf("%s", err)
	                os.Exit(1)
	        	}
	        	product_id=uberRes.Prices[0].ProductID

				if len(bestrouteArr)>=putCounter{
					putCounter=putCounter+1
				}					
				var jsonStr = []byte(`{"start_longitude":`+fmt.Sprint("", startLng) +`,"start_latitude":`+fmt.Sprint("", startLat) +`,"product_id":"`+product_id+`"}`)				
				url= fmt.Sprintf("https://sandbox-api.uber.com/v1/requests?start_latitude=%f&start_longitude=%f&end_latitude=%f&end_longitude=%f&product_id=%s", startLat,startLng,endLat,endLng,product_id)

				req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
				req.Header.Set("Authorization", "Bearer"+access_token)
				req.Header.Set("Content-Type", "application/json")

				client:=http.Client{}
				resp, err := client.Do(req)
				if err != nil {
	    			panic(err)
				}				
				err=json.NewDecoder(resp.Body).Decode(&data)
				if err!=nil {
					fmt.Println(err)
				}
				request_id=data.(map[string]interface{})["request_id"]
				request_eta=data.(map[string]interface{})["eta"]
				var jsonStatus=[]byte(`{"status":"accepted"}`)
				if _, ok := request_id.(string); ok {
				    url="https://sandbox-api.uber.com/v1/sandbox/requests/"+request_id.(string)
				    req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStatus))
				    if err!=nil {
						fmt.Println(err)
					}
				    req.Header.Set("Authorization", "Bearer"+access_token)
					req.Header.Set("Content-Type", "application/json")
					response, err := client.Do(req)
					if response.StatusCode==204{
						responseStruct= &PutReqResponse{
						ID: bson.ObjectIdHex(tripId),
						Status: statusUpdate,
						StartingFromLocationID: result.StartingFromLocationID,
						NextDestinationLocationID: nextDestiId,
						BestRouteLocationIds: bestrouteArr,
						TotalUberCosts : result.TotalUberCosts,
						TotalUberDuration : result.TotalUberDuration,
						TotalUberDistance : result.TotalUberDistance,
						UberWaitimeETA: request_eta.(float64),
						}
						outgoingJSON, err := json.Marshal(responseStruct)
						if err != nil {
							log.Println(err)
							http.Error(res, err.Error(), http.StatusInternalServerError)
							return
						}
						res.WriteHeader(http.StatusOK)
						fmt.Fprint(res, string(outgoingJSON))
					} else {
								fmt.Printf("Not 204 response from Uber API\n")
			    				os.Exit(1) }					
				} else {
				    fmt.Printf("Cannot convert request_id to string\n")
			    	os.Exit(1)
				}
		} else{
			 outgoingJSON, err := json.Marshal(responseStruct)
				if err != nil {
					log.Println(err)
					http.Error(res, err.Error(), http.StatusInternalServerError)
					return
				}
			res.WriteHeader(http.StatusOK)
			fmt.Fprint(res, string(outgoingJSON))
		}
	}
}

func handlePostTrips(res http.ResponseWriter, req *http.Request) {
		bestRouteFinalStringArr := make([]string, 0)
		res.Header().Set("Content-Type", "application/json")
		tripInput:= new(PostRequest)
		decoder := json.NewDecoder(req.Body)
		error := decoder.Decode(&tripInput)
		if error != nil {
			log.Println(error.Error())
			http.Error(res, error.Error(), http.StatusInternalServerError)
			return
		}
		inputString:=tripInput.StartingFromLocationID

		finTotalCost:=0
		finTotalDuration:=0
		finTotalDistance:=0.0

		modifiedArray := make([]string, 0)
		bestRouteLocation,totalCost,totalDuration,totalDistance:=bestRouteLocationId(inputString,tripInput.LocationIds)	
		finTotalCost+=totalCost
		finTotalDuration+=totalDuration
		finTotalDistance+=totalDistance


		bestRouteFinalStringArr=append(bestRouteFinalStringArr, bestRouteLocation)

	outerLoop:
		for len(tripInput.LocationIds)>=len(bestRouteFinalStringArr){
			if len(bestRouteFinalStringArr)==1{
				modifiedArray=removeDuplicate(bestRouteFinalStringArr[0],tripInput.LocationIds)
				bestRouteLocation,totalCost,totalDuration,totalDistance=bestRouteLocationId(bestRouteFinalStringArr[0],modifiedArray)
				finTotalCost+=totalCost
				finTotalDuration+=totalDuration
				finTotalDistance+=totalDistance
				bestRouteFinalStringArr=append(bestRouteFinalStringArr, bestRouteLocation)
				if len(tripInput.LocationIds)==len(bestRouteFinalStringArr){
					break outerLoop
				}
			}
			if len(bestRouteFinalStringArr)==2{
				modifiedArray=removeDuplicate(bestRouteFinalStringArr[1],modifiedArray)
				bestRouteLocation,totalCost,totalDuration,totalDistance=bestRouteLocationId(bestRouteFinalStringArr[1],modifiedArray)	
				totalCost,totalDuration,totalDistance=fetchLowestValues(bestRouteFinalStringArr[1],modifiedArray[0])
				finTotalCost+=totalCost
				finTotalDuration+=totalDuration
				finTotalDistance+=totalDistance
				bestRouteFinalStringArr=append(bestRouteFinalStringArr, bestRouteLocation)
				if len(tripInput.LocationIds)==len(bestRouteFinalStringArr){
					break outerLoop
				}
			}
			if len(bestRouteFinalStringArr)==3{
				modifiedArray=removeDuplicate(bestRouteFinalStringArr[2],modifiedArray)
				bestRouteLocation,totalCost,totalDuration,totalDistance=bestRouteLocationId(bestRouteFinalStringArr[2],modifiedArray)
				finTotalCost+=totalCost
				finTotalDuration+=totalDuration
				finTotalDistance+=totalDistance
				bestRouteFinalStringArr=append(bestRouteFinalStringArr, bestRouteLocation)
				if len(tripInput.LocationIds)==len(bestRouteFinalStringArr){
					break outerLoop
				}
			}
			if len(bestRouteFinalStringArr)==4{
				modifiedArray=removeDuplicate(bestRouteFinalStringArr[3],modifiedArray)
				bestRouteLocation,totalCost,totalDuration,totalDistance=bestRouteLocationId(bestRouteFinalStringArr[3],modifiedArray)	
				finTotalCost+=totalCost
				finTotalDuration+=totalDuration
				finTotalDistance+=totalDistance
				bestRouteFinalStringArr=append(bestRouteFinalStringArr, bestRouteLocation)
				if len(tripInput.LocationIds)==len(bestRouteFinalStringArr){
					break outerLoop
				}
			}
			if len(bestRouteFinalStringArr)==5{
				modifiedArray=removeDuplicate(bestRouteFinalStringArr[4],modifiedArray)
				bestRouteLocation,totalCost,totalDuration,totalDistance=bestRouteLocationId(bestRouteFinalStringArr[4],modifiedArray)
				finTotalCost+=totalCost
				finTotalDuration+=totalDuration
				finTotalDistance+=totalDistance
				bestRouteFinalStringArr=append(bestRouteFinalStringArr, bestRouteLocation)
				if len(tripInput.LocationIds)==len(bestRouteFinalStringArr){
					break outerLoop
				}
			}
			if len(bestRouteFinalStringArr)==6{
				modifiedArray=removeDuplicate(bestRouteFinalStringArr[5],modifiedArray)
				bestRouteLocation,totalCost,totalDuration,totalDistance=bestRouteLocationId(bestRouteFinalStringArr[5],modifiedArray)
				finTotalCost+=totalCost
				finTotalDuration+=totalDuration
				finTotalDistance+=totalDistance
				bestRouteFinalStringArr=append(bestRouteFinalStringArr, bestRouteLocation)
				if len(tripInput.LocationIds)==len(bestRouteFinalStringArr){
					break outerLoop
				}
			}
			if len(bestRouteFinalStringArr)==7{
				modifiedArray=removeDuplicate(bestRouteFinalStringArr[6],modifiedArray)
				bestRouteLocation,totalCost,totalDuration,totalDistance=bestRouteLocationId(bestRouteFinalStringArr[6],modifiedArray)
				finTotalCost+=totalCost
				finTotalDuration+=totalDuration
				finTotalDistance+=totalDistance
				bestRouteFinalStringArr=append(bestRouteFinalStringArr, bestRouteLocation)
				if len(tripInput.LocationIds)==len(bestRouteFinalStringArr){
					break outerLoop
				}
			}
			if len(bestRouteFinalStringArr)==8{
				modifiedArray=removeDuplicate(bestRouteFinalStringArr[7],modifiedArray)
				bestRouteLocation,totalCost,totalDuration,totalDistance=bestRouteLocationId(bestRouteFinalStringArr[7],modifiedArray)
				finTotalCost+=totalCost
				finTotalDuration+=totalDuration
				finTotalDistance+=totalDistance
				bestRouteFinalStringArr=append(bestRouteFinalStringArr, bestRouteLocation)
				if len(tripInput.LocationIds)==len(bestRouteFinalStringArr){
					break outerLoop
				}
			}
		}

		rountTripCost,rountTripDuration,rountTripDistance:=fetchLowestValues(bestRouteFinalStringArr[len(bestRouteFinalStringArr)-1],inputString)
		finTotalCost+=rountTripCost
		finTotalDuration+=rountTripDuration
		finTotalDistance+=rountTripDistance

		responseStruct := &PostResponse{
			ID: bson.NewObjectId(),
			Status: "Planning",
			StartingFromLocationID: inputString,
			BestRouteLocationIds: bestRouteFinalStringArr,
			TotalUberCosts : finTotalCost,
			TotalUberDuration : finTotalDuration,
			TotalUberDistance : finTotalDistance,
		}
			sess, err := mgo.Dial("mongodb://kalpana:pass@ds035844.mongolab.com:35844/cmpe273")
			if err != nil {
			    fmt.Printf("Can't connect to mongo, go error 2 %v\n", err)
			    os.Exit(1)
			}
			defer sess.Close()				 
			sess.SetSafe(&mgo.Safe{})
			collection := sess.DB("cmpe273").C("TripPlanner")
			err = collection.Insert(responseStruct)
			if err != nil {
			    fmt.Printf("Can't insert document: %v\n", err)
			    os.Exit(1)
			}
			outgoingJSON, err := json.Marshal(responseStruct)
			if err != nil {
			log.Println(error.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
			}
			res.WriteHeader(http.StatusCreated)
			fmt.Fprint(res, string(outgoingJSON))
}

func coordinateValues(locationid string) (lat,lng float64) {
		sess, err := mgo.Dial("mongodb://kalpana:pass@ds035844.mongolab.com:35844/cmpe273")
		if err != nil {
		    fmt.Printf("Can't connect to mongo, go error 3 %v\n", err)
		    os.Exit(1)
		}
		defer sess.Close()				 
		sess.SetSafe(&mgo.Safe{})
		collection := sess.DB("cmpe273").C("AddressBook")

		result:= ResponseAddOp{}
		err = collection.Find(bson.M{"_id": bson.ObjectIdHex(locationid)}).Select(bson.M{}).One(&result)
		if err != nil {
		    fmt.Printf("Address not found: %v\n", err)
		}
		lat=result.Coordinate.Lat
		lng=result.Coordinate.Lng
		return 
}

func fetchLowestValues(startLocation,nextLocation string) (costEstimate int,duration int,distance float64) {	
			startLat, startLng := coordinateValues(startLocation)
			endLat, endLng := coordinateValues(nextLocation)

			url := fmt.Sprintf("https://sandbox-api.uber.com/v1/estimates/price?start_latitude=%f&start_longitude=%f&end_latitude=%f&end_longitude=%f&server_token=OnEamv71AJVwxWiXz_8_aN5LzsAZ-R6E_FcRI8ug", startLat,startLng,endLat,endLng)
			var uberRes UberResponse
 			responseValue, err := http.Get(url)
            defer responseValue.Body.Close()
            reply, err := ioutil.ReadAll(responseValue.Body)
            json.Unmarshal([]byte(reply), &uberRes)
            if err != nil {
                fmt.Printf("%s", err)
                os.Exit(1)
        	}
			costEstimate=uberRes.Prices[0].LowEstimate
        	duration=uberRes.Prices[0].Duration
        	distance=uberRes.Prices[0].Distance

			for i := 0; i < len(uberRes.Prices); i++ {
				if uberRes.Prices[i].LowEstimate < costEstimate {
						if uberRes.Prices[i].LowEstimate!=0{
                         costEstimate = uberRes.Prices[i].LowEstimate
                         }
                 }
                 if uberRes.Prices[i].Duration < duration {
						if uberRes.Prices[i].Duration!=0{
                         duration = uberRes.Prices[i].Duration
                         }
                 }
                 if uberRes.Prices[i].Distance < distance {
						if uberRes.Prices[i].Distance!=0{
                         distance = uberRes.Prices[i].Distance
                         }
                 }
			}
        	return 
}

func bestRouteLocationId(startLocation string,locationArray []string) (bestRoute string,cost,duration int,distance float64){
		minCost:=0
		minDuration:=0
		for i := 0; i < len(locationArray); i++ {
			cost,duration,distance=fetchLowestValues(startLocation,locationArray[i])
			if cost < minCost || minCost==0 {
					if cost!=0{
	                 	minCost = cost
					if duration < minDuration || minDuration==0 {
						if duration!=0{
		                 minDuration = duration
		                }
            		}
	            	bestRoute=locationArray[i]
	            }
            }			
    	}
    	return
}

func removeDuplicate(duplicateStr string,locationArray []string)(modifiedLocationArray []string) {
	for i := 0; i < len(locationArray); i++ {
		if duplicateStr==locationArray[i]{
		} else {
			modifiedLocationArray=append(modifiedLocationArray,locationArray[i])
		}
	}
	return modifiedLocationArray
}