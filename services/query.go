package services

const AddProjectQuery = `
INSERT INTO Projects 
	(name,description)
VALUES (?,?)
`

const AddWorkDetailsQuery = `
INSERT INTO work_details 
	(user_id,project_id,log_hours,work_details) 
VALUES (?,?,?,?)
`

const GetUsersQuery = `
SELECT 
	u.user_name 
FROM Users u
`
const GetUserByNameQuery = `
SELECT 
	id 
FROM users 
WHERE user_name = ?
`
const GetProjectByNameQuery = `
SELECT 
	id 
FROM projects 
WHERE Lower(name) = ?
`

const CheckProjectExistQuery = `
SELECT 
	count(*) 
FROM projects 
WHERE name = ?
`

const ActivateProjectQuery = `
UPDATE projects 
	SET is_active = true 
WHERE name = ?
`

const GetMemberProjectReport = `
SELECT 
	sum(wd.log_hours) as log_hours, p.id as project_id, p.name as project_name 
FROM users u
	LEFT JOIN work_details wd on wd.user_id = u.id
	LEFT JOIN Projects p on p.id = wd.project_id
WHERE u.user_name = ?
group by u.id,p.id
`
const RemoveProjectQuery = `
UPDATE Projects 
	SET is_active = false 
WHERE Lower(name) = ?
`

const GetProjectsQuery = `
SELECT 
	p.name 
FROM Projects p 
WHERE p.is_active = true
`

const GetProjectReportQuery = `
SELECT 
	u.id as user_id, u.user_name, sum(wd.log_hours) as log_hours 
FROM projects p
	LEFT JOIN work_details wd on wd.project_id = p.id
	LEFT JOIN users u on u.id = wd.user_id
WHERE Lower(p.name) = ?
GROUP By u.id, p.id
`

const GetMemberReportWeekly = `
SELECT 
	sum(wd.log_hours) as log_hours, p.id as project_id, p.name as project_name 
FROM users u
	LEFT JOIN work_details wd on wd.user_id = u.id
	LEFT JOIN Projects p on p.id = wd.project_id
WHERE u.user_name = ? and wd.created_at BETWEEN NOW() - INTERVAL '1 week' AND NOW()
group by u.id,p.id
`

const GetMemberReportMonthly = `
SELECT 
	sum(wd.log_hours) as log_hours, p.id as project_id, p.name as project_name 
FROM users u
	LEFT JOIN work_details wd on wd.user_id = u.id
	LEFT JOIN Projects p on p.id = wd.project_id
WHERE u.user_name = ? and wd.created_at >= (CURRENT_DATE - INTERVAL '1 month')::date
group by u.id,p.id
`

const GetMemberReportQuaterly = `
SELECT 
	sum(wd.log_hours) as log_hours, p.id as project_id, p.name as project_name 
FROM users u
	LEFT JOIN work_details wd on wd.user_id = u.id
	LEFT JOIN Projects p on p.id = wd.project_id
WHERE u.user_name = ? and wd.created_at >= (CURRENT_DATE - INTERVAL '3 month')::date
group by u.id,p.id
`

const GetProjectReportWeekly = `
SELECT 
	u.id as user_id, u.user_name, sum(wd.log_hours) as log_hours 
FROM projects p
	LEFT JOIN work_details wd on wd.project_id = p.id
	LEFT JOIN users u on u.id = wd.user_id
WHERE Lower(p.name) = ? and wd.created_at BETWEEN NOW() - INTERVAL '1 week' AND NOW()
GROUP By u.id, p.id
`
const GetProjectReportMonthly = `
SELECT 
	u.id as user_id, u.user_name, sum(wd.log_hours) as log_hours 
FROM projects p
	LEFT JOIN work_details wd on wd.project_id = p.id
	LEFT JOIN users u on u.id = wd.user_id
WHERE Lower(p.name) = ? and wd.created_at >= (CURRENT_DATE - INTERVAL '1 month')::date
GROUP By u.id, p.id
`
const GetProjectReportQuaterly = `
SELECT 
	u.id as user_id, u.user_name, sum(wd.log_hours) as log_hours 
FROM projects p
	LEFT JOIN work_details wd on wd.project_id = p.id
	LEFT JOIN users u on u.id = wd.user_id
WHERE Lower(p.name) = ? and wd.created_at >= (CURRENT_DATE - INTERVAL '3 month')::date
GROUP By u.id, p.id
`
