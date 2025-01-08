package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"transport-api/internal/models"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type BusImpl struct {
	DB *pgx.Conn
}

func registerBusRoutes(router *mux.Router, dbpool *pgx.Conn) {
	impl := &BusImpl{
		DB: dbpool,
	}

	router.HandleFunc("/buses/stops/latitude/{latitude}/longitude/{longitude}", impl.GetNearestStops).Methods("GET")

	router.HandleFunc("/buses/company/{company_slug}/route/{route_id}/direction/{direction}", impl.GetRouteByCompanyRouteDirection).Methods("GET")
}

func (impl *BusImpl) GetNearestStops(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	latitude := vars["latitude"]
	longitude := vars["longitude"]
	queryParams := r.URL.Query()
	searchNearestStopsLimit := queryParams.Get("searchNearestStopsLimit")

	limit := 30
	if searchNearestStopsLimit != "" {
		parsedLimit, err := strconv.Atoi(searchNearestStopsLimit)
		if err != nil {
			http.Error(w, "Invalid searchNearestStopsLimit", http.StatusBadRequest)
			return
		}
		limit = parsedLimit
	}

	var nearestStops []models.NearestStop
	nearestStops, _ = models.GetActiveNearestStops(r.Context(), impl.DB, 1, latitude, longitude, limit)

	response := map[string]interface{}{
		"data": map[string]interface{}{},
		"meta": map[string]interface{}{
			"request_info": map[string]interface{}{
				"timestamp": time.Now().UTC().Format(time.RFC3339),
				"latitude":  latitude,
				"longitude": longitude,
			},
		},
	}

	nearestRoutes := []map[string]interface{}{}
	for _, nearestStop := range nearestStops {
		stop, _ := models.GetNearestStopByStopId(r.Context(), impl.DB, nearestStop.Stop)

		var routeStops []*models.RouteStop
		routeStops, _ = models.GetActiveRouteStopsByStopId(r.Context(), impl.DB, 1, nearestStop.Stop, "1")
		for _, routeStop := range routeStops {
			route, _ := models.GetFirstActiveRouteByTransportTypeAndRouteStop(r.Context(), impl.DB, 1, routeStop)
			if route == nil {
				continue
			}

			company, _ := models.GetCompanyById(r.Context(), impl.DB, *route.CompanyID)

			route_info := map[string]interface{}{
				"route": route.Route,
				"bound": route.Bound,
				"service": map[string]interface{}{
					"type": route.ServiceType,
					"name": route.ServiceNameEn,
				},
				"origin":      route.OriginEn,
				"destination": route.DestinationEn,
				"company": map[string]interface{}{
					"slug":         company.SlugEn,
					"display_slug": company.SlugEn,
					"name":         company.NameEn,
					"meta": map[string]interface{}{
						"bg_color": company.BgColor,
						"color":    company.TextColor,
					},
				},
				"is_collaborated":  false,
				"is_available_now": map[string]interface{}{},
			}

			nearestStop := map[string]interface{}{
				"stop":               nearestStop.Stop,
				"sequence":           routeStop.Sequence,
				"name":               stop.NameEn,
				"latitude":           stop.Latitude,
				"longitude":          stop.Longitude,
				"distance_from_user": nearestStop.Distance,
			}

			nearestRoute := map[string]interface{}{
				"route_info":   route_info,
				"nearest_stop": nearestStop,
			}

			nearestRoutes = append(nearestRoutes, nearestRoute)
		}
	}
	response["data"] = nearestRoutes

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (impl *BusImpl) GetRouteByCompanyRouteDirection(w http.ResponseWriter, r *http.Request) {
	response := "test"

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
