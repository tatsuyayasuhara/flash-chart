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
var Cp = 56.188 // モル熱容量 //todo 変数
var Mw = 123.233e-03 // 分子量 //todo 変数
var den = 6.10e+03 // 密度 //todo 変数
var vm0 = Mw/den
var Cpv = Cp/vm0
var N0 = 7.0e+20 // 単位体積当たりの電子の有効状態密度(酸素イオン空孔形成に関わる) //todo 変数
var dSv = RR // 酸素イオン空孔の形成エントロピー　//todo 変数
var B0 = math.Pow(2.0 * math.Exp(dSv/RR), 1.0/3.0)
var M0 = 8.02e-02   // 温度の関数の場合もある //todo 変数
var dE = 1.89       // [eV] エネルギー障壁的なもの？再配列エネルギー？ //todo 変数
var EM = dE*FF
var zz = -1.5       // 局所的な価数(電荷) //todo 変数
var phi = 5.0       // 外部局所電位(電場ポテンシャル) //todo 変数
var zp = zz*phi
var extPote = zp*FF // [J/mol]
var dt = 1.0
var dTFdt = 1.0/6.0
var ds = 0.01*dt

func (dff *DataFrameFlash) CalcW(voltage, vacancyE, pressure float64) {
	//nd := NDP - 1
	//Cp := 56.188
	//vm0 := 123.233e-03/(6.10e+03)
	//Cp := Cp/vm0
	//E0 := 1.6021766208e-19
	//N0 := 7.0e+20
	//B0 := math.Pow(2.0 * math.Exp(RR/RR), 1.0/3.0)
	//EM := 1.89*FF
	//M0 := 8.02e-02
	//zzphi := -7.5
	//extPote := zzphi*FF
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
