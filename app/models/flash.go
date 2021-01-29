package models

import (
	"math"
)

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

const RR float64 = 8.3145
const FF float64 = 96485.33289 // [C/mol]
const E0 float64 = 1.6021766208e-19 // [C]
var dt = 1.0
var ds = 0.01*dt

func (dff *DataFrameFlash) CalcW(voltage, vacancyE, pressure, heatCap, molarVolume, electronDOS, formationEntropy,
	mobilityCoef, rearrangeE, localValency, localPotential, heatingRate, semiBandGap, semiEDOS float64, enableVol bool) {
	var Cp = heatCap //56.188 // モル熱容量
	var vm0 = molarVolume //モル体積：2.02×10^-5[m^3/mol] (=123.233e-03/(6.10e+03))
	var Cpv = Cp/vm0
	var N0 = electronDOS //7.0e+20 // 単位体積当たりの電子の有効状態密度(酸素イオン空孔形成に関わる)
	var dSv = formationEntropy //RR // 酸素イオン空孔の形成エントロピー
	var B0 = math.Pow(2.0 * math.Exp(dSv/RR), 1.0/3.0)
	var M0 = mobilityCoef //8.02e-02   // 温度の関数の場合もある
	var EM = rearrangeE //電子のホッピング伝導における活性化エネルギー：1.824×10^5 [J/mol]　(=1.89×96485.33289)
	var zz = localValency //-3       // 局所的な価数(電荷)
	var phi = localPotential //3.75       // 外部局所電位(電場ポテンシャル)
	var zp = zz*phi
	var extPote = zp*FF // [J/mol]
	dTFdt := heatingRate // 1.0/6.0[K/s]
	Ebg := semiBandGap //4.824×10^5 [J/mol] (=5.0×96485.33289) 真性半導体におけるバンドギャップ
	Neds := semiEDOS //1.31e+28[m-3] 真性半導体における電子の状態密度

	TF := 773.0
	T := TF
	// 外部入力3項
	vol := voltage
	dHv := vacancyE
	pO2 := pressure
	var tList []float64
	var wList []float64
	nd := int(math.Ceil(1000/dTFdt))
	for ii := 0; ii <= nd; ii++ {
		TF += dTFdt * dt
		tList = append(tList, TF - 273) //todo 初回のみの処理にしたい
		for i := 1; i <= 100; i++ {
			n1 := N0 * B0 * math.Pow(pO2, -1.0/6.0) * math.Exp(-(dHv+extPote)/(3.0*RR*T))
			n2 := 0.0
			if enableVol {
				n2 = Neds * math.Exp(-(Ebg-extPote)/(2.0*RR*T))
			} else {
				n2 = Neds * math.Exp(-Ebg/(2.0*RR*T))
			}
			Mn := M0 * math.Exp(-EM/(RR*T))
			dTdt := vol * vol * E0 /Cpv * (n1 + n2) * Mn
			T += dTdt * ds
		}
		if T <= TF { T = TF }
		n1 := N0 * B0 * math.Pow(pO2, -1.0/6.0) * math.Exp(-(dHv+extPote)/(3.0*RR*T))
		n2 := 0.0
		if enableVol {
			n2 = Neds * math.Exp(-(Ebg-extPote)/(2.0*RR*T))
		} else {
			n2 = Neds * math.Exp(-Ebg/(2.0*RR*T))
		}
		Mn := M0 * math.Exp(-EM/(RR*T))
		I := vol * E0 * (n1 + n2) * Mn
		W := I * vol
		wList = append(wList, W)
	}

	dff.Temperature = &Temperature{Values:tList}
	dff.PowerDissipation = &PowerDissipation{Values:wList}
}
