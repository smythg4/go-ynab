package service

import (
	"context"
	"fmt"
	"go-ynab/ynab"
	"sort"
	"time"
)

// TODO: Gut and rework this

type MonthCategory struct {
	Month         string
	Cat           ynab.Category
	PlanId        string
	PlanLastMonth string
}

func (mc *MonthCategory) String() string {
	return fmt.Sprintf("[%s] %s: $%.2f / $%.2f", mc.Month, mc.Cat.Name, float64(mc.Cat.Activity)/1000, float64(mc.Cat.Budgeted)/1000)
}

type Service struct {
	client     *ynab.Client
	TimeFrame  int
	Plans      []ynab.PlanDetails
	Categories []MonthCategory
}

func NewService(key string, window int) *Service {
	return &Service{
		client:    ynab.NewClient(key).WithRateLimit(200, 10),
		TimeFrame: window,
	}
}

func (s *Service) FetchPlans(ctx context.Context) error {
	s.Plans = []ynab.PlanDetails{}
	allPlans, err := s.client.GetPlans(ctx)
	if err != nil {
		return err
	}
	for _, plan := range allPlans {
		if err := ctx.Err(); err != nil {
			return err
		}
		if plan.LastMonth.After(time.Now().AddDate(0, -s.TimeFrame, 0)) {
			fullPlan, err := s.client.GetPlan(ctx, plan.ID)
			if err != nil {
				return err
			}
			s.Plans = append(s.Plans, *fullPlan)
		}
	}
	return nil
}

func (s *Service) FetchCategories() error {
	if len(s.Plans) == 0 {
		return fmt.Errorf("no plans stored in this service")
	}
	s.Categories = []MonthCategory{}

	for _, plan := range s.Plans {
		for _, month := range plan.Months {
			if month.Month.Before(time.Now()) && month.Month.After(time.Now().AddDate(0, -s.TimeFrame, 0)) {
				for _, cat := range month.Categories {
					s.Categories = append(s.Categories, MonthCategory{
						Month:         month.Month.Format("2006-01-02"),
						Cat:           cat,
						PlanId:        plan.ID.String(),
						PlanLastMonth: plan.LastMonth.Time.Format("2006-01-02"),
					})
				}
			}
		}
	}
	sort.Slice(s.Categories, func(i, j int) bool {
		return s.Categories[i].PlanLastMonth > s.Categories[j].PlanLastMonth
	})
	return nil
}
