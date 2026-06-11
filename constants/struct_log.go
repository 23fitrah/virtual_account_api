package constants

import "time"

func HistoryLogParserCloseDataMt199(newStatus string, dataMaintenance DataMaintenance, dataSts DataSts, dataGPI DataGPI, isGpiEligible bool) map[string]map[string]interface{} {
	newVal := map[string]map[string]interface{}{
		"BRISWIFTMAINTENANCEMX": {
			"STATUS":      newStatus,
			"USERAPPROVE": dataMaintenance.UserApprove,
			"ACTION":      dataMaintenance.Action,
			"TGLAPPROVE":  dataMaintenance.TglApprove,
		},
		"BRISWIFTSTSTRXMX": {
			"CLOSEDATE":    dataSts.CloseDate,
			"KD_STATUS":    dataSts.KdStatus,
			"CLOSEREMARKS": dataSts.CloseRemarks,
			"PROCESS_DATE": dataSts.ProcessDate,
		},
	}

	if isGpiEligible {
		newVal["BRISWIFTGPI"] = map[string]interface{}{
			"ROWID":                     dataGPI.RowID,
			"REFF":                      dataGPI.Reff,
			"KD_STATUS":                 dataGPI.KdStatus,
			"FROM":                      "BRINIDJA",
			"BUSINESSSERVICE":           "001",
			"TRANSACTIONIDENTIFICATION": dataGPI.TransactionIdentification,
			"INSTRUCTIONIDENTIFICATION": dataGPI.Reff,
			"ORIGINATOR":                "BRINIDJA",
			"STATUS":                    dataGPI.Status,
			"REASON":                    dataGPI.Reason,
			"AMOUNT":                    dataGPI.Amount,
			"CURRENCY":                  dataGPI.Currency,
			"TRANSACTIONDATE":           dataGPI.TransactionDate,
			"CRAMOUNT":                  dataGPI.CrAmount,
			"CRCURR":                    dataGPI.CrCurrency,
			"RATE":                      dataGPI.Rate,
			"CHARGES":                   dataGPI.Charges,
			"SENDSTS":                   "0",
			"SENDDESC":                  "waiting send",
			"FUNDAVAILABLE":             time.Now().UTC().Format("2006-01-02 15:04:05"),
			"KEKINIAN":                  time.Now().Format("2006-01-02 15:04:05"),
		}
	}

	return newVal
}

func HistoryLogRejectDataMt199(id, userMaintenance string) map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		"BRISWIFTMAINTENANCEMX": {
			"ID":          id,
			"STATUS":      "0",
			"USERAPPROVE": userMaintenance,
			"ACTION":      "Reject Trx from WS",
			"TGLAPPROVE":  time.Now().Format("2006-01-02 15:04:05"),
		},
	}
}

func HistoryLogFlagDataMt199(dataMaintenance DataMaintenance, dataSts DataSts, dataGPI DataGPI, isGpiEligible bool) map[string]map[string]interface{} {
	newVal := map[string]map[string]interface{}{
		"BRISWIFTMAINTENANCEMX": {
			"STATUS":      "9",
			"USERAPPROVE": dataMaintenance.UserApprove,
			"ACTION":      dataMaintenance.Action,
			"TGLAPPROVE":  dataMaintenance.TglApprove,
		},
		"BRISWIFTSTSTRXMX": {
			"CLOSEDATE":    dataSts.CloseDate,
			"KD_STATUS":    dataSts.KdStatus,
			"CLOSEREMARKS": dataSts.CloseRemarks,
		},
	}

	if isGpiEligible {
		newVal["BRISWIFTGPI"] = map[string]interface{}{
			"ROWID":                     dataGPI.RowID,
			"REFF":                      dataGPI.Reff,
			"KD_STATUS":                 dataGPI.KdStatus,
			"FROM":                      "BRINIDJA",
			"BUSINESSSERVICE":           "001",
			"TRANSACTIONIDENTIFICATION": dataGPI.TransactionIdentification,
			"INSTRUCTIONIDENTIFICATION": dataGPI.Reff,
			"ORIGINATOR":                "BRINIDJA",
			"STATUS":                    dataGPI.Status,
			"REASON":                    dataGPI.Reason,
			"AMOUNT":                    dataGPI.Amount,
			"CURRENCY":                  dataGPI.Currency,
			"TRANSACTIONDATE":           dataGPI.TransactionDate,
			"CRAMOUNT":                  dataGPI.CrAmount,
			"CRCURR":                    dataGPI.CrCurrency,
			"RATE":                      dataGPI.Rate,
			"CHARGES":                   dataGPI.Charges,
			"SENDSTS":                   "0",
			"SENDDESC":                  "waiting send",
			"FUNDAVAILABLE":             time.Now().UTC().Format("2006-01-02 15:04:05"),
			"KEKINIAN":                  time.Now().Format("2006-01-02 15:04:05"),
		}
	}

	return newVal
}

func HistoryLogReleaseDataMt199(traceCounter string, dataMaintenance DataMaintenance, dataSts DataSts, typeSA string) map[string]map[string]interface{} {
	newVal := map[string]map[string]interface{}{
		"BRISWIFTMAINTENANCEMX": {
			"STATUS":      "0",
			"USERAPPROVE": dataMaintenance.UserApprove,
			"ACTION":      dataMaintenance.Action,
			"TGLAPPROVE":  dataMaintenance.TglApprove,
		},
	}

	newVal["BRISWIFTSTSTRXMX"] = map[string]interface{}{
		"REMARK2":      dataSts.Remark,
		"PROCESS_DATE": dataSts.ProcessDate,
	}

	if typeSA == "SAK" {
		newVal["BRISWIFTSTSTRXMX"] = map[string]interface{}{
			"NAMA_ASLI_TRANSAKSI": dataSts.NamaAsliTrx,
			"BENEFNAME":           dataSts.BenefName,
		}
	} else if typeSA == "SANK" {
		newVal["BRISWIFTSTSTRXMX"] = map[string]interface{}{
			"REFF_TRACER_AMEND":   dataSts.ReffTracerAmend,
			"PASSCHECKCHARGES":    dataSts.PassCheckCharges,
			"CHARGES":             dataSts.Charges,
			"CHARGESAMENDMENT":    dataSts.ChargesAmendment,
			"TRACE_COUNTER":       traceCounter,
			"KD_STATUS":           dataSts.KdStatus,
			"NAMA_ASLI_TRANSAKSI": dataSts.NamaAsliTrx,
			"BENEFNAME":           dataSts.BenefName,
		}
	}

	return newVal
}
