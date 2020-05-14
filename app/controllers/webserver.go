package controllers

import (
	"encoding/json"
	"flash-chart/app/models"
	"flash-chart/config"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

func StartWebServer() error {
	port := config.Config.Port
	http.HandleFunc("/calc/flash/", calcFlashHandler)
	http.HandleFunc("/chart/", viewChartHandler)
	log.Printf("start server at :%d, open on your browser -> http://localhost:%d/chart", port, port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
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
	ev, _ := strconv.ParseFloat(externalVoltage, 64)
	ve, _ := strconv.ParseFloat(vacancyEnergy, 64)
	pO2, _ := strconv.ParseFloat(oxygenPressure, 64)
	//log.Printf("ev=%g, ve=%g, pO2=%g\n", ev, ve, pO2)

	dff := CalcW(ev, ve, pO2)
	js, err := json.Marshal(dff)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

const NDP int = 6001
const RR float64 = 8.3145
const FF float64 = 96485.33289
var nd = NDP - 1
var Cp = 56.188
var vm0 = 123.233e-03/(6.10e+03)
var Cpv = Cp/vm0
var e0 = 1.6021766208e-19
var N0 = 7.0e+20
var B0 = math.Pow(2.0 * math.Exp(RR/RR), 1.0/3.0)
var MQ0 = 1.89*FF
var MM0 = 8.02e-02
var zp = -7.5
var ext_pote = zp*FF
var dt = 1.0
var dTFdt = 1.0/6.0
var ds = 0.01*dt

func CalcW(voltage, vacancyE, pressure float64) *models.DataFrameFlash {
	//nd := NDP - 1
	//Cp := 56.188
	//vm0 := 123.233e-03/(6.10e+03)
	//Cp := Cp/vm0
	//e0 := 1.6021766208e-19
	//N0 := 7.0e+20
	//B0 := math.Pow(2.0 * math.Exp(RR/RR), 1.0/3.0)
	//MQ0 := 1.89*FF
	//MM0 := 8.02e-02
	//zp := -7.5
	//ext_pote := zp*FF
	//dt := 1.0
	//dTFdt := 1.0/6.0
	//ds := 0.01*dt
	// ----------------------
	TF := 773.0
	T := TF
	// 外部入力3項
	vol := voltage
	dHv := vacancyE * FF
	pO2 := pressure
	var tList []float64
	var wList []float64
	for ii := 0; ii <= nd; ii++ {
		TF += dTFdt * dt
		tList = append(tList, TF - 273) //todo 初回のみの処理にしたい
		for i := 1; i <= 100; i++ {
			n1 := N0 * B0 * math.Pow(pO2, -1.0/6.0) * math.Exp(-(dHv+ext_pote)/3./RR/T)
			dTdt := vol * vol * e0/Cpv * n1 * MM0 * math.Exp(-MQ0/RR/T)
			T += dTdt*ds
		}
		if T <= TF { T = TF }
		n1 := N0 * B0 * math.Pow(pO2, -1.0/6.0) * math.Exp(-(dHv+ext_pote)/3./RR/T)
		I := vol * e0 * n1 * MM0 * math.Exp(-MQ0/RR/T)
		W := I * vol
		wList = append(wList, W)
	}
	return &models.DataFrameFlash{
		Temperature:      models.Temperature{Values: tList},
		PowerDissipation: models.PowerDissipation{Values: wList},
	}
}
