package db

func Tasks(limit int) ([]*Task, error) {
	row, err := DB.Query(
		"SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?",
		limit,
	)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	var tasks = make([]*Task, 0, limit)
	for row.Next() {
		var task Task
		if err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	if err := row.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
