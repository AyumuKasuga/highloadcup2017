package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func updateUsers(userList []user) {
	for _, item := range userList {
		allUsers[item.Id] = item
	}
}

func updateLocations(locationList []location) {
	for _, item := range locationList {
		allLocations[item.Id] = item
	}
}

func updateVisits(visitList []visit) {
	for _, item := range visitList {
		allVisits[item.Id] = item
		if val, ok := allUsersVisit[item.User]; ok {
			allUsersVisit[item.User] = append(val, item.Id)
		} else {
			allUsersVisit[item.User] = []int{item.Id}
		}
		if val, ok := allLocationsVisit[item.Location]; ok {
			allLocationsVisit[item.Location] = append(val, item.Id)
		} else {
			allLocationsVisit[item.Location] = []int{item.Id}
		}

	}
}

func loadFromFile() {
	allUsers = make(map[int]user)
	allLocations = make(map[int]location)
	allVisits = make(map[int]visit)
	allUsersVisit = make(map[int][]int)
	allLocationsVisit = make(map[int][]int)
	r, err := zip.OpenReader("/tmp/data/data.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer rc.Close()
		if strings.Contains(f.Name, "locations") {
			var locationList locations
			buf, _ := ioutil.ReadAll(rc)
			json.Unmarshal(buf, &locationList)
			updateLocations(locationList.Locations)
		} else if strings.Contains(f.Name, "users") {
			var userList users
			buf, _ := ioutil.ReadAll(rc)
			json.Unmarshal(buf, &userList)
			updateUsers(userList.Users)
		} else if strings.Contains(f.Name, "visits") {
			var visitList visits
			buf, _ := ioutil.ReadAll(rc)
			json.Unmarshal(buf, &visitList)
			updateVisits(visitList.Visits)
		}

	}
	fmt.Println("ready!")
}
