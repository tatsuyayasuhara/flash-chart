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
	molarVolume := r.URL.Query().Get("molar_volume")
	electronDOS := r.URL.Query().Get("electron_DOS")
	formationEntropy := r.URL.Query().Get("formation_entropy")
	mobilityCoef := r.URL.Query().Get("mobility_coef")
	rearrangeE := r.URL.Query().Get("rearrange_energy")
	localValency := r.URL.Query().Get("local_valency")
	localPotential := r.URL.Query().Get("local_potential")
	heatingRate := r.URL.Query().Get("heating_rate")
	semiBandGap := r.URL.Query().Get("semi_band_gap")
	semiEDOS := r.URL.Query().Get("semi_EDOS")
	ev, _ := strconv.ParseFloat(externalVoltage, 64)
	ve, _ := strconv.ParseFloat(vacancyEnergy, 64)
	pO2, _ := strconv.ParseFloat(oxygenPressure, 64)
	Cp, _ := strconv.ParseFloat(heatCapacity, 64)
	mv, _ := strconv.ParseFloat(molarVolume, 64)
	eDOS, _ := strconv.ParseFloat(electronDOS, 64)
	fe, _ := strconv.ParseFloat(formationEntropy, 64)
	mc, _ := strconv.ParseFloat(mobilityCoef, 64)
	re, _ := strconv.ParseFloat(rearrangeE, 64)
	lv, _ := strconv.ParseFloat(localValency, 64)
	lp, _ := strconv.ParseFloat(localPotential, 64)
	hr, _ := strconv.ParseFloat(heatingRate, 64)
	ebg, _ := strconv.ParseFloat(semiBandGap, 64)
	neds, _ := strconv.ParseFloat(semiEDOS, 64)
	log.Printf("ev=%g, ve=%g, pO2=%g, Cp=%g, mv=%g, eDOS=%g, fe=%g, mc=%g, re=%g, lv=%g, lp=%g, hr=%g, ebg=%g, neds=%g\n",
		ev, ve, pO2, Cp, mv, eDOS, fe, mc, re, lv, lp, hr, ebg, neds)

	dff := new(models.DataFrameFlash)
	dff.CalcW(ev, ve, pO2, Cp, mv, eDOS, fe, mc, re, lv, lp, hr, ebg, neds)
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
