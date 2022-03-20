/*
@Time : 2022/3/20 00:53
@Author : StudentCWZ
@File : job
@Software: GoLand
*/

package models

type Job struct {
	Type              string `json:"type"`
	JobID             string `json:"job_id"`
	CoID              string `json:"co_id"`
	JobHref           string `json:"job_href"`
	JobName           string `json:"job_name"`
	CompanyHref       string `json:"company_href"`
	CompanyName       string `json:"company_name"`
	ProvideSalaryText string `json:"provide_salary_text"`
	WorkAreaText      string `json:"work_area_text"`
	CompanyTypeText   string `json:"company_type_text"`
	IssueDate         string `json:"issue_date"`
	JobWelf           string `json:"job_welf"`
	AttributeText     string `json:"attribute_text"`
	CompanySizeText   string `json:"company_size_text"`
	CompanyIndText    string `json:"company_ind_text"`
}
