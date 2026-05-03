package report

import (
	"encoding/csv"
	"fmt"
	"go-ynab/service"
	"io"
	"strings"
)

type ReportCategory int

const (
	Charity ReportCategory = iota
	NonEmployerSavings
	Mortgage
	EssHomeRepair
	HomeRenovations
	Utilities
	Phone
	HealthCare
	OtherInsurance
	Therapy
	CaneloCosts
	AutoMaintenance
	Groceries
	Gas
	Gifts
	DateNight
	SubServices
	FunFitness
	ClothesHygiene
	Vacation
	ChildCare
	CollegeSavings
)

func (rc ReportCategory) String() string {
	switch rc {
	case Charity:
		return "Charity"
	case NonEmployerSavings:
		return "Non-Employer Savings"
	case Mortgage:
		return "Mortgage"
	case EssHomeRepair:
		return "Essential Home Repair"
	case HomeRenovations:
		return "Home Renovations"
	case Utilities:
		return "Utilities"
	case Phone:
		return "Phone"
	case HealthCare:
		return "Health Care"
	case OtherInsurance:
		return "Other Insurance"
	case Therapy:
		return "Therapy"
	case CaneloCosts:
		return "Canelo Costs"
	case AutoMaintenance:
		return "Auto Maintenance"
	case Groceries:
		return "Groceries"
	case Gas:
		return "Gas"
	case Gifts:
		return "Gifts"
	case DateNight:
		return "Date Night"
	case SubServices:
		return "Subscription Services"
	case FunFitness:
		return "Fun & Fitness"
	case ClothesHygiene:
		return "Clothes & Hygiene"
	case Vacation:
		return "Vacation"
	case ChildCare:
		return "Child Care"
	case CollegeSavings:
		return "College Savings"
	default:
		return "Unknown"
	}
}

var keywordMap = []struct {
	keyword string
	cat     ReportCategory
}{
	{"Health", HealthCare},
	{"Gift", Gifts},
	{"Charit", Charity},
	{"Mortgage", Mortgage},
	{"Clothes", ClothesHygiene},
	{"Hygiene", ClothesHygiene},
	{"Date", DateNight},
	{"Fitness", FunFitness},
	{"Fun", FunFitness},
	{"Grocer", Groceries},
	{"Canelo", CaneloCosts},
	{"Auto", AutoMaintenance},
	{"Gas", Gas},
	{"Fuel", Gas},
	{"Vacation", Vacation},
	{"Trip", Vacation},
	{"College", CollegeSavings},
	{"Child", ChildCare},
	{"Subscrip", SubServices},
	{"Therapy", Therapy},
	{"Counseling", Therapy},
	{"Insurance", OtherInsurance},
	{"Renovation", HomeRenovations},
	{"Home", EssHomeRepair},
	{"Savings", NonEmployerSavings},
	{"Utilit", Utilities},
	{"Phone", Phone},
	{"Car", AutoMaintenance},
}

func FromString(s string) (ReportCategory, error) {
	cleaned := strings.ToLower(strings.Trim(s, " "))
	for _, entry := range keywordMap {
		if strings.Contains(cleaned, strings.ToLower(entry.keyword)) {
			return entry.cat, nil
		}
	}
	return -1, fmt.Errorf("unknown category: %s", s)
}

type Report struct {
	Categories map[ReportCategory][]service.MonthCategory
	Totals     map[ReportCategory]CategoryTotal
	Months     int
}

type CategoryTotal struct {
	Budgeted int64
	Activity int64
}

func (ct CategoryTotal) String() string {
	return fmt.Sprintf("Activity: $%.2f / Budgeted: $%.2f", float64(ct.Activity)/1000, float64(ct.Budgeted)/1000)
}

func NewReport(mc []service.MonthCategory, months int) Report {
	cats := make(map[ReportCategory][]service.MonthCategory)
	claimedBy := make(map[ReportCategory]map[string]string)
	for _, item := range mc {
		rc, err := FromString(item.Cat.Name)
		if err != nil {
			continue
		}
		if claimedBy[rc] == nil {
			claimedBy[rc] = make(map[string]string)
		}
		if owner := claimedBy[rc][item.Month]; owner != "" && owner != item.PlanId {
			continue
		}
		claimedBy[rc][item.Month] = item.PlanId
		cats[rc] = append(cats[rc], item)
	}
	return Report{Categories: cats, Months: months}
}

func (r *Report) CalculateTotals() {
	r.Totals = map[ReportCategory]CategoryTotal{}
	for rc, mcs := range r.Categories {
		for _, mc := range mcs {
			act := -mc.Cat.Activity
			bud := mc.Cat.Budgeted
			existing := r.Totals[rc]
			r.Totals[rc] = CategoryTotal{
				Budgeted: existing.Budgeted + bud,
				Activity: existing.Activity + act}
		}
	}
}

func (r *Report) WriteCsv(w io.Writer) error {
	cwriter := csv.NewWriter(w)
	err := cwriter.Write([]string{"Category", "Activity", "Budgeted"})
	if err != nil {
		return err
	}
	for rc, ct := range r.Totals {
		cat := rc.String()
		act := fmt.Sprintf("%.2f", float64(ct.Activity)/1000/float64(r.Months))
		budg := fmt.Sprintf("%.2f", float64(ct.Budgeted)/1000/float64(r.Months))
		err = cwriter.Write([]string{cat, act, budg})
		if err != nil {
			return err
		}
	}
	cwriter.Flush()
	return cwriter.Error()
}
