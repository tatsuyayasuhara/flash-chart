package models

import "math"

type DataFrameFlash struct {
	Temperature      *Temperature      `json:"temperature,omitempty"`
	PowerDissipation *PowerDissipation `json:"power_dissipation,omitempty"`
}

type Temperature struct {
	Values []float64 `json:"values,omitempty"`
}

type PowerDissipation struct {
	Values []float64 `json:"values,omitempty"`
}

const NDP int = 6001
const RR float64 = 8.3145
const FF float64 = 96485.33289 // [C/mol]
const E0 float64 = 1.6021766208e-19 // [C]
var nd = NDP - 1
var dt = 1.0
var dTFdt = 1.0/6.0
var ds = 0.01*dt

func (dff *DataFrameFlash) CalcW(voltage, vacancyE, pressure, heatCap, molWeight, density, electronDOS,
	formationEntropy, mobilityCoef, rearrangeE, localValency, localPotential float64) {
	var Cp = heatCap //56.188 // モル熱容量
	var Mw = molWeight*1e-03 //123.233e-03 // 分子量 [kg/mol]
	var den = density*1e+03 //6.10e+03 // 密度
	var vm0 = Mw/den
	var Cpv = Cp/vm0
	var N0 = electronDOS //7.0e+20 // 単位体積当たりの電子の有効状態密度(酸素イオン空孔形成に関わる)
	var dSv = formationEntropy //RR // 酸素イオン空孔の形成エントロピー
	var B0 = math.Pow(2.0 * math.Exp(dSv/RR), 1.0/3.0)
	var M0 = mobilityCoef //8.02e-02   // 温度の関数の場合もある
	var dE = rearrangeE //1.89       // [eV] エネルギー障壁的なもの？再配列エネルギー？
	var EM = dE*FF
	var zz = localValency //-1.5       // 局所的な価数(電荷)
	var phi = localPotential //5.0       // 外部局所電位(電場ポテンシャル)
	var zp = zz*phi
	var extPote = zp*FF // [J/mol]

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
			n1 := N0 * B0 * math.Pow(pO2, -1.0/6.0) * math.Exp(-(dHv+extPote)/(3.0*RR*T))
			Mn := M0 * math.Exp(-EM/(RR*T))
			dTdt := vol * vol * E0 /Cpv * n1 * Mn
			T += dTdt * ds
		}
		if T <= TF { T = TF }
		n1 := N0 * B0 * math.Pow(pO2, -1.0/6.0) * math.Exp(-(dHv+extPote)/(3.0*RR*T))
		Mn := M0 * math.Exp(-EM/(RR*T))
		I := vol * E0 * n1 * Mn
		W := I * vol
		wList = append(wList, W)
	}

	dff.Temperature = &Temperature{Values:tList}
	dff.PowerDissipation = &PowerDissipation{Values:wList}
}
