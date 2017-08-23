package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
)

func getUserVisits(userID int, ctx *fasthttp.RequestCtx) (string, error) {
	fromDate, err := strconv.Atoi(string(ctx.FormValue("fromDate")))
	if len(ctx.FormValue("fromDate")) != 0 && err != nil {
		return "", errorBadRequest
	}
	toDate, err := strconv.Atoi(string(ctx.FormValue("toDate")))
	if len(ctx.FormValue("toDate")) != 0 && err != nil {
		return "", errorBadRequest
	}

	toDistance, err := strconv.Atoi(string(ctx.FormValue("toDistance")))
	if len(ctx.FormValue("toDistance")) != 0 && err != nil {
		return "", errorBadRequest
	}

	country := string(ctx.FormValue("country"))

	visitList := make(map[int]string)
	visitedAtKeys := []int{}

	for _, i := range allUsersVisit[userID] {
		v := allVisits[i]
		condition := fromDate != 0 && v.VisitedAt <= fromDate ||
			toDate != 0 && v.VisitedAt >= toDate ||
			len(country) != 0 && allLocations[v.Location].Country != country ||
			toDistance != 0 && allLocations[v.Location].Distance >= toDistance

		if !condition {
			visitList[v.VisitedAt] = fmt.Sprintf(
				`{"visited_at":%d,"mark":%d,"place":"%s"}`,
				v.VisitedAt,
				v.Mark,
				allLocations[v.Location].Place,
			)
			visitedAtKeys = append(visitedAtKeys, v.VisitedAt)
		}
	}

	sort.Ints(visitedAtKeys)

	sortedVisitList := make([]string, len(visitedAtKeys))

	for i, k := range visitedAtKeys {
		sortedVisitList[i] = visitList[k]
	}

	return strings.Join(sortedVisitList, ","), nil
}

func usersHandler(ctx *fasthttp.RequestCtx) {
	urlStr := string(ctx.Path()[7:])
	if ctx.IsGet() {
		userID, err := strconv.Atoi(urlStr)
		if err != nil {
			if !strings.Contains(urlStr, "/visits") {
				ctx.Response.SetStatusCode(404)
			} else {
				userID, _ := strconv.Atoi(urlStr[:len(urlStr)-7])
				if _, ok := allUsers[userID]; ok {
					userVisits, err := getUserVisits(userID, ctx)
					if err != nil {
						ctx.Response.SetStatusCode(400)
					} else {
						if len(userVisits) == 0 {
							fmt.Fprintf(ctx, `{"visits": []}`)
						} else {
							fmt.Fprintf(ctx, `{"visits":[%s]}`, userVisits)
						}
					}
				} else {
					ctx.Response.SetStatusCode(404)
				}
			}

		} else {
			if val, ok := allUsers[userID]; ok {
				fmt.Fprintf(
					ctx,
					`{"id":%d,"email":"%s","first_name":"%s","last_name":"%s","gender":"%s","birth_date":%d}`,
					val.Id, val.Email, val.FirstName, val.LastName, val.Gender, val.BirthDate,
				)
			} else {
				ctx.Response.SetStatusCode(404)
			}
		}
	} else if ctx.IsPost() {
		if urlStr == "new" {
			buf := ctx.PostBody()
			if len(buf) == 0 || strings.Contains(string(buf), ": null") {
				ctx.Response.SetStatusCode(400)
				return
			}
			var in user
			err := json.Unmarshal(buf, &in)
			if err != nil {
				ctx.Response.SetStatusCode(400)
			} else {
				allUsers[in.Id] = in
				ctx.SetBody(emptyJson)
			}
		} else {
			userID, _ := strconv.Atoi(urlStr)
			if userItem, ok := allUsers[userID]; ok {
				buf := ctx.PostBody()
				if len(buf) == 0 {
					ctx.Response.SetStatusCode(400)
					return
				}
				var in map[string]interface{}
				json.Unmarshal(buf, &in)
				if val, ok := in["email"]; ok {
					if val == nil {
						ctx.Response.SetStatusCode(400)
						return
					}
					userItem.Email = val.(string)
				}
				if val, ok := in["first_name"]; ok {
					if val == nil {
						ctx.Response.SetStatusCode(400)
						return
					}
					userItem.FirstName = val.(string)
				}
				if val, ok := in["last_name"]; ok {
					if val == nil {
						ctx.Response.SetStatusCode(400)
						return
					}
					userItem.LastName = val.(string)
				}
				if val, ok := in["gender"]; ok {
					if val == nil {
						ctx.Response.SetStatusCode(400)
						return
					}
					userItem.Gender = val.(string)
				}
				if val, ok := in["birth_date"]; ok {
					if val == nil {
						ctx.Response.SetStatusCode(400)
						return
					}
					userItem.BirthDate = int(val.(float64))
				}
				if userItem.Email == "" || userItem.FirstName == "" || userItem.LastName == "" || userItem.Gender == "" || userItem.BirthDate == 0 {
					ctx.Response.SetStatusCode(400)
				} else {
					allUsers[userID] = userItem
					ctx.SetBody(emptyJson)
				}
			} else {
				ctx.Response.SetStatusCode(404)
			}

		}
	}
}

func getLocationAvg(locationID int, ctx *fasthttp.RequestCtx) (float64, error) {
	fromDate, err := strconv.Atoi(string(ctx.FormValue("fromDate")))
	if len(ctx.FormValue("fromDate")) != 0 && err != nil {
		return .0, errorBadRequest
	}
	toDate, err := strconv.Atoi(string(ctx.FormValue("toDate")))
	if len(ctx.FormValue("toDate")) != 0 && err != nil {
		return .0, errorBadRequest
	}
	fromAge, err := strconv.Atoi(string(ctx.FormValue("fromAge")))
	if len(ctx.FormValue("fromAge")) != 0 && err != nil {
		return .0, errorBadRequest
	}
	toAge, err := strconv.Atoi(string(ctx.FormValue("toAge")))
	if len(ctx.FormValue("toAge")) != 0 && err != nil {
		return .0, errorBadRequest
	}
	gender := string(ctx.FormValue("gender"))
	if gender != "" && gender != "f" && gender != "m" {
		return .0, errorBadRequest
	}

	var sum int
	var count float64

	for _, i := range allLocationsVisit[locationID] {
		v := allVisits[i]
		condition := fromDate != 0 && v.VisitedAt <= fromDate ||
			toDate != 0 && v.VisitedAt >= toDate ||
			len(gender) != 0 && allUsers[v.User].Gender != gender ||
			fromAge != 0 && currentTime-allUsers[v.User].BirthDate < fromAge*oneYear ||
			toAge != 0 && currentTime-allUsers[v.User].BirthDate > toAge*oneYear

		if !condition {
			count++
			sum += v.Mark
		}
	}

	if sum == 0 {
		return .0, nil
	}

	return float64(sum) / count, nil
}

func locationsHandler(ctx *fasthttp.RequestCtx) {
	urlStr := string(ctx.Path()[11:])
	if ctx.IsGet() {
		locationID, err := strconv.Atoi(urlStr)
		if err != nil {
			if !strings.Contains(urlStr, "/avg") {
				ctx.Response.SetStatusCode(400)
			} else {
				locationID, _ := strconv.Atoi(urlStr[:len(urlStr)-4])
				if _, ok := allLocations[locationID]; ok {
					avg, err := getLocationAvg(locationID, ctx)
					if err != nil {
						ctx.Response.SetStatusCode(400)
					} else {
						fmt.Fprintf(ctx, "{\"avg\": %.5f }", avg)
					}
				} else {
					ctx.Response.SetStatusCode(404)
				}
			}

		} else {
			if val, ok := allLocations[locationID]; ok {
				fmt.Fprintf(
					ctx,
					`{"id":%d,"place":"%s","country":"%s","city":"%s","distance":%d}`,
					val.Id, val.Place, val.Country, val.City, val.Distance,
				)
			} else {
				ctx.Response.SetStatusCode(404)
			}
		}
	} else if ctx.IsPost() {
		if urlStr == "new" {
			buf := ctx.PostBody()
			if len(buf) == 0 || strings.Contains(string(buf), ": null") {
				ctx.Response.SetStatusCode(400)
				return
			}
			var in location
			err := json.Unmarshal(buf, &in)
			if err != nil {
				ctx.Response.SetStatusCode(400)
			} else {
				allLocations[in.Id] = in
				ctx.SetBody(emptyJson)
			}

		} else {
			locationID, _ := strconv.Atoi(urlStr)
			if locationItem, ok := allLocations[locationID]; ok {
				buf := ctx.PostBody()
				if len(buf) == 0 {
					ctx.Response.SetStatusCode(400)
					return
				}
				var in map[string]interface{}
				json.Unmarshal(buf, &in)
				if val, ok := in["place"]; ok {
					if val == nil {
						ctx.Response.SetStatusCode(400)
						return
					}
					locationItem.Place = val.(string)
				}
				if val, ok := in["country"]; ok {
					if val == nil {
						ctx.Response.SetStatusCode(400)
						return
					}
					locationItem.Country = val.(string)
				}
				if val, ok := in["city"]; ok {
					if val == nil {
						ctx.Response.SetStatusCode(400)
						return
					}
					locationItem.City = val.(string)
				}
				if val, ok := in["distance"]; ok {
					if val == nil {
						ctx.Response.SetStatusCode(400)
						return
					}
					locationItem.Distance = int(val.(float64))
				}
				if locationItem.Place == "" || locationItem.Country == "" || locationItem.City == "" {
					ctx.Response.SetStatusCode(400)
				} else {
					allLocations[locationID] = locationItem
					ctx.SetBody(emptyJson)
				}
			} else {
				ctx.Response.SetStatusCode(404)
			}

		}
	}
}

func remove(slice []int, s int) []int {
	return append(slice[:s], slice[s+1:]...)
}

func visitsHandler(ctx *fasthttp.RequestCtx) {
	urlStr := string(ctx.Path()[8:])
	if ctx.IsGet() {
		visitID, err := strconv.Atoi(urlStr)
		if err != nil {
			ctx.Response.SetStatusCode(400)
		} else {
			if val, ok := allVisits[visitID]; ok {
				fmt.Fprintf(
					ctx,
					`{"id":%d,"location":%d,"user":%d,"visited_at":%d,"mark":%d}`,
					val.Id, val.Location, val.User, val.VisitedAt, val.Mark,
				)
			} else {
				ctx.Response.SetStatusCode(404)
			}
		}
	} else if ctx.IsPost() {
		if urlStr == "new" {
			buf := ctx.PostBody()
			if len(buf) == 0 || strings.Contains(string(buf), ": null") {
				ctx.Response.SetStatusCode(400)
				return
			}
			var in visit
			err := json.Unmarshal(buf, &in)
			if err != nil {
				ctx.Response.SetStatusCode(400)
			} else {
				allVisits[in.Id] = in
				if val, ok := allUsersVisit[in.User]; ok {
					allUsersVisit[in.User] = append(val, in.Id)
				} else {
					allUsersVisit[in.User] = []int{in.Id}
				}
				if val, ok := allLocationsVisit[in.Location]; ok {
					allLocationsVisit[in.Location] = append(val, in.Id)
				} else {
					allLocationsVisit[in.Location] = []int{in.Id}
				}
				ctx.SetBody(emptyJson)
			}

		} else {
			visitID, _ := strconv.Atoi(urlStr)
			if visitItem, ok := allVisits[visitID]; ok {
				buf := ctx.PostBody()
				if len(buf) == 0 {
					ctx.Response.SetStatusCode(400)
					return
				}
				var in map[string]interface{}
				json.Unmarshal(buf, &in)
				if val, ok := in["location"]; ok {
					if val == nil {
						ctx.Response.SetStatusCode(400)
						return
					}
					visitItem.Location = int(val.(float64))
				}
				if val, ok := in["user"]; ok {
					if val == nil {
						ctx.Response.SetStatusCode(400)
						return
					}
					visitItem.User = int(val.(float64))
				}
				if val, ok := in["visited_at"]; ok {
					if val == nil {
						ctx.Response.SetStatusCode(400)
						return
					}
					visitItem.VisitedAt = int(val.(float64))
				}
				if val, ok := in["mark"]; ok {
					if val == nil {
						ctx.Response.SetStatusCode(400)
					} else {
						visitItem.Mark = int(val.(float64))
					}
				}
				if visitItem.Location == 0 || visitItem.User == 0 || visitItem.VisitedAt == 0 || visitItem.Mark == 0 {
					ctx.Response.SetStatusCode(400)
				} else {
					if allVisits[visitID].User != visitItem.User {
						for i, v := range allUsersVisit[allVisits[visitID].User] {
							if v == visitID {
								allUsersVisit[allVisits[visitID].User] = remove(allUsersVisit[allVisits[visitID].User], i)
								allUsersVisit[visitItem.User] = append(allUsersVisit[visitItem.User], visitID)
								break
							}
						}
					}
					if allVisits[visitID].Location != visitItem.Location {
						for i, v := range allLocationsVisit[allVisits[visitID].Location] {
							if v == visitID {
								allLocationsVisit[allVisits[visitID].Location] = remove(allLocationsVisit[allVisits[visitID].Location], i)
								allLocationsVisit[visitItem.Location] = append(allLocationsVisit[visitItem.Location], visitID)
								break
							}
						}
					}
					allVisits[visitID] = visitItem
					ctx.SetBody(emptyJson)
				}
			} else {
				ctx.Response.SetStatusCode(404)
			}
		}
	}
}