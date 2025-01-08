package routes

import (
	"encoding/json"
	"net/http"
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

	var routes []*models.Route
	routes, _ = models.GetActiveRoutesByTransportType(r.Context(), impl.DB, 1)
	var nearestStops []models.NearestStop
	nearestStops, _ = models.GetActiveNearestStops(r.Context(), impl.DB, latitude, longitude, 10)

	response := map[string]interface{}{
		"data": routes,
		"meta": map[string]interface{}{
			"latitude":  latitude,
			"longitude": longitude,
		},
	}

	for index, nearestStop := range nearestStops {
		if index > 10 {
			continue
		}

		var routeStops []*models.RouteStop
		routeStops, _ = models.GetActiveRouteStopsByStopId(r.Context(), impl.DB, 1, nearestStop.Stop, "1")
		for _, routeStop := range routeStops {
			route, _ := models.GetFirstActiveRouteByTransportTypeAndRouteStop(r.Context(), impl.DB, 1, routeStop)

		}
	}

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
