package errand

const (
	ByQualification Restriction = iota
	ByVerification
	ByInsurance
)

func RestrictionType(value string) Restriction {
	if value == "qualification" {
		return ByQualification
	}
	if value == "verification" {
		return ByVerification
	}
	if value == "insurance" {
		return ByInsurance
	}
	return -1
}
