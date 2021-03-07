package bed

import (
	"opencensus/core/common"
	"opencensus/core/dkan"
	"strconv"
	"time"
)

type Record struct {
	Name string

	Institution string
	Code        string
	Category    string
	Correlative string

	CutDate      time.Time
	RegisterDate time.Time

	Region   string
	Province string
	// District string
	Ubigeo string

	ZNCBusy      int
	ZNCAvailable int
	ZNCTotal     int

	ZCBusy      int
	ZCAvailable int
	ZCTotal     int

	HospBusy      int
	HospAvailable int
	HospTotal     int

	UCINBusy      int
	UCINAvailable int
	UCINTotal     int

	UCIPediatraBusy      int
	UCIPediatraAvailable int
	UCIPediatraTotal     int

	UCIAdultBusy      int
	UCIAdultAvailable int
	UCIAdultTotal     int

	MainSourceKind string
}

const dateLayout = "20060102"

func Extract(lapses int) chan []Record {
	channelRecords := make(chan []Record)
	api, err := dkan.NewAPI(common.SuSaludDKANEndpoint)
	if err != nil {
		panic(err)
	}

	oxygenRes := dkan.ResourceWithID("0badb458-b44f-49ad-bd32-51acaaee7a05")
	oxygenRes.First100()

	go extractor(api, oxygenRes, lapses, channelRecords)

	return channelRecords
}

func extractor(api *dkan.API, res *dkan.Resource, lapses int, channel chan []Record) {
	for lapse := 0; lapse < lapses; lapse++ {
		data, err := api.ObtainResource(res)
		if err != nil {
			panic(err) // really?
		}

		records := data["records"].([]interface{})

		recordsArray := []Record{}

		for _, r := range records {
			rec := r.(map[string]interface{})

			name, _ := rec["NOMBRE"].(string)

			institution, _ := rec["INSTITUCION"].(string)
			category, _ := rec["CATEGORIA"].(string)
			correlative, _ := rec["CORRELATIVO"].(string)

			code, _ := rec["CODIGO"].(string)

			cutDate, _ := rec["FECHACORTE"].(int)
			registerDate, _ := rec["FECHAREGISTRO"].(int)

			cutDateTime, _ := time.Parse(dateLayout, string(cutDate))
			registerDateTime, _ := time.Parse(dateLayout, string(registerDate))
			// cutDate

			region, _ := rec["REGION"].(string)
			province, _ := rec["PROVINCIA"].(string)
			// district, _ := rec["DISTRITO"].(string)
			ubigeo, _ := rec["UBIGEO"].(string)

			hospAvailable, _ := rec["CAMAS_HOSP_DISPONIBLE"].(string)
			hospBusy, _ := rec["CAMAS_HOSP_OCUPADAS"].(string)
			hospTotal, _ := rec["CAMAS_HOSP_TOTAL"].(string)

			ucinAvailable, _ := rec["UCIN_CAMAS_DISPONIBLE"].(string)
			ucinBusy, _ := rec["UCIN_CAMAS_OCUPADAS"].(string)
			ucinTotal, _ := rec["UCIN_CAMAS_TOTAL"].(string)

			uciAdultAvailable, _ := rec["UCI_ADULTOS_CAMAS_DISPONIBLE"].(string)
			uciAdultBusy, _ := rec["UCI_ADULTOS_CAMAS_OCUPADAS"].(string)
			uciAdultTotal, _ := rec["UCI_ADULTOS_CAMAS_TOTAL"].(string)

			uciPediatraAvailable, _ := rec["UCI_PEDIATRIA_CAMAS_DISPONIBLE"].(string)
			uciPediatraBusy, _ := rec["UCI_PEDIATRIA_CAMAS_OCUPADAS"].(string)
			uciPediatraTotal, _ := rec["UCI_PEDIATRIA_CAMAS_TOTAL"].(string)

			zcAvailable, _ := rec["CAMAS_ZC_DISPONIBLES"].(string)
			zcBusy, _ := rec["CAMAS_ZC_OCUPADOS"].(string)
			zcTotal, _ := rec["CAMAS_ZC_TOTAL"].(string)

			zncAvailable, _ := rec["CAMAS_ZNC_DISPONIBLE"].(string)
			zncBusy, _ := rec["CAMAS_ZNC_OCUPADOS"].(string)
			zncTotal, _ := rec["CAMAS_ZNC_TOTAL"].(string)

			mainSourceKind := "idk"

			hospAvailableNumber, _ := strconv.Atoi(hospAvailable)
			hospBusyNumber, _ := strconv.Atoi(hospBusy)
			hospTotalNumber, _ := strconv.Atoi(hospTotal)

			ucinAvailableNumber, _ := strconv.Atoi(ucinAvailable)
			ucinBusyNumber, _ := strconv.Atoi(ucinBusy)
			ucinTotalNumber, _ := strconv.Atoi(ucinTotal)

			uciAdultAvailableNumber, _ := strconv.Atoi(uciAdultAvailable)
			uciAdultBusyNumber, _ := strconv.Atoi(uciAdultBusy)
			uciAdultTotalNumber, _ := strconv.Atoi(uciAdultTotal)

			uciPediatraAvailableNumber, _ := strconv.Atoi(uciPediatraAvailable)
			uciPediatraBusyNumber, _ := strconv.Atoi(uciPediatraBusy)
			uciPediatraTotalNumber, _ := strconv.Atoi(uciPediatraTotal)

			zcAvailableNumber, _ := strconv.Atoi(zcAvailable)
			zcBusyNumber, _ := strconv.Atoi(zcBusy)
			zcTotalNumber, _ := strconv.Atoi(zcTotal)

			zncAvailableNumber, _ := strconv.Atoi(zncAvailable)
			zncBusyNumber, _ := strconv.Atoi(zncBusy)
			zncTotalNumber, _ := strconv.Atoi(zncTotal)

			recordsArray = append(recordsArray, Record{
				Name:         name,
				Institution:  institution,
				Code:         code,
				Category:     category,
				Correlative:  correlative,
				CutDate:      cutDateTime,
				RegisterDate: registerDateTime,
				Region:       region,
				Province:     province,
				// District:           district,
				Ubigeo: ubigeo,

				ZNCBusy:      zncBusyNumber,
				ZNCAvailable: zncAvailableNumber,
				ZNCTotal:     zncTotalNumber,

				ZCBusy:      zcBusyNumber,
				ZCAvailable: zcAvailableNumber,
				ZCTotal:     zcTotalNumber,

				HospBusy:      hospBusyNumber,
				HospAvailable: hospAvailableNumber,
				HospTotal:     hospTotalNumber,

				UCINAvailable: ucinAvailableNumber,
				UCINBusy:      ucinBusyNumber,
				UCINTotal:     ucinTotalNumber,

				UCIAdultBusy:      uciAdultBusyNumber,
				UCIAdultAvailable: uciAdultAvailableNumber,
				UCIAdultTotal:     uciAdultTotalNumber,

				UCIPediatraBusy:      uciPediatraBusyNumber,
				UCIPediatraAvailable: uciPediatraAvailableNumber,
				UCIPediatraTotal:     uciPediatraTotalNumber,

				MainSourceKind: mainSourceKind,
			})
		}

		channel <- recordsArray
		res.NextN(100)
	}

	close(channel)
}
