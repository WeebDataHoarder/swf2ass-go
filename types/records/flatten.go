package records

func FlattenRecord(r Record, scale int64) []Record {
	if cr, ok := r.(CurvedRecord); ok {
		return cr.ToLineRecords(scale)
	} else if lr, ok := r.(LineRecord); ok {
		return []Record{lr}
	} else {
		panic("not supported")
		return nil
	}
}
