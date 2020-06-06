package controllers

import (
	"encoding/json"
	"flash-chart/app/models"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func StartWebServer(addr string) error {
	if addr == ":" {
		addr = fmt.Sprintf(":%d", 8080)
	}
	http.HandleFunc("/calc/flash/", calcFlashHandler)
	http.HandleFunc("/chart/", viewChartHandler)
	log.Printf("start server at %s, open on your browser -> http://localhost%s/chart", addr, addr)
	return http.ListenAndServe(fmt.Sprintf("%s", addr), nil)
}

var templates = template.Must(template.ParseFiles("app/views/chart.html"))

func viewChartHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "chart.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func calcFlashHandler(w http.ResponseWriter, r *http.Request) {
	externalVoltage := r.URL.Query().Get("external_voltage")
	vacancyEnergy := r.URL.Query().Get("vacancy_energy")
	oxygenPressure := r.URL.Query().Get("oxygen_pressure")
	heatCapacity := r.URL.Query().Get("heat_capacity")
	ev, _ := strconv.ParseFloat(externalVoltage, 64)
	ve, _ := strconv.ParseFloat(vacancyEnergy, 64)
	pO2, _ := strconv.ParseFloat(oxygenPressure, 64)
	Cp, _ := strconv.ParseFloat(heatCapacity, 64)
	log.Printf("ev=%g, ve=%g, pO2=%g, Cp=%g\n", ev, ve, pO2, Cp)

	dff := new(models.DataFrameFlash)
	dff.CalcW(ev, ve, pO2, Cp)
	js, err := json.Marshal(dff)
	if err != nil {
		log.Printf("marshal error: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(js); err != nil {
		log.Printf("write json error: %v", err)
	}
}
