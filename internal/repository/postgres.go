package repository

import (
	"database/sql"
	"fmt"
	"donor-monitor/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib" // не вызываем функции библиотеки, просто инициализируем
)

type PostgresRepository struct{
	db *sql.DB // указатель на объект, управляющий пулом соединений с бд. доступ к бд
}

// конкретных конструкторов в го нет как в c++ поэтому принято использовать функции с New 
func NewPostgresRepository(dsn string)(*PostgresRepository, error){
	// передаетсяф строка подключения к бд с пользователем, паролем, адресом и тд

	db, err := sql.Open("pgx", dsn)  // объект подключения
	if err != nil{
		return nil, fmt.Errorf("error to connect: %w", dsn)
	}
	
	// .Open() не гарантитрует подключение. вызывается метод .Ping()
	err = db.Ping() 
	if err != nil{
		return nil, fmt.Errorf("error ping: ", err)
	}

	return &PostgresRepository{
		db: db,
	}, nil
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}

func (r *PostgresRepository) SaveCenter(
	center model.DonorCenter,
) (int64, error){
	
	query := `
		INSERT INTO centers (name, city, source_url)
		VALUES ($1,$2, $3)
		ON CONFLICT (source_url)
		DO UPDATE SET
			name = EXCLUDED.name,
			city = EXCLUDED.city
		RETURNING id;
	`

	var id int64

	err := r.db.QueryRow(
		query,
		center.NameCenter,
		center.City,
		center.URL,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to save center: %w", err)
	}

	return id, nil

}

func (r * PostgresRepository) SaveCheck(
	CenterID int64,
	statuses []model.DonorStatus,
) error {

	if len(statuses) != 8 {
		return fmt.Errorf("not 8 statuses")
	}

	tx, err := r.db.Begin()// начать транзакцию 
	if err != nil{
		return fmt.Errorf("begin trensaction: %w", err)
	}

	defer tx.Rollback() // отмена изменений при ошибке

	queryCheck := `
		INSERT INTO checks (center_id, checked_at, success)
		VALUES ($1, $2, $3)
		RETURNING id;
	`
	
	var checkID int64

	err = tx.QueryRow(
		queryCheck,
		CenterID,
		statuses[0].CheckedAt,
		true,
	).Scan(&checkID)

	if err != nil{
		return fmt.Errorf("save check: %w", err)
	}

	queryStatus := `
		INSERT INTO statuses (check_id, blood_group, color)
		VALUES ($1, $2, $3)
	` 

	for _, status := range statuses{
		_, err := tx.Exec(
			queryStatus,
			checkID,
			status.BloodGroup,
			status.Color,
		)
		if err != nil{
			return fmt.Errorf("save status %s %w", status.BloodGroup, err)
		}
	}

	if err := tx.Commit(); err != nil{
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetLatest(centerID int64) ([]model.DonorStatus, error){
	query := `
	SELECT
		s.blood_group,
		s.color,
		ch.checked_at
	FROM statuses s
	JOIN checks ch ON ch.id = s.check_id
	WHERE ch.id = (
		SELECT id
		FROM checks
		WHERE center_id = $1
			AND success = TRUE
		ORDER BY checked_at DESC, id DESC
		LIMIT 1
	)
	ORDER BY s.id;
	`

	rows, err := r.db.Query(query, centerID)
	if err != nil{
		return nil, fmt.Errorf("get latest: %w", err)
	}

	defer rows.Close()

	statuses := []model.DonorStatus{}

	for rows.Next(){//проход по строкам полученным из запроса
		var status model.DonorStatus

		err := rows.Scan(	// запись в перемменную текущую строку
			&status.BloodGroup,
			&status.Color, 
			&status.CheckedAt,
		)
		if err != nil{
			return nil, fmt.Errorf("scan status: %w",err)
		}

		statuses = append(statuses, status)
	}
	if err := rows.Err(); err != nil{
		return  nil, fmt.Errorf("read statuses: %w" ,err)
	}
	return statuses, nil
}


func (r *PostgresRepository) GetHistory(
	centerID int64,
	group string,
)([]model.DonorStatus, error){
	query := `
	SELECT
		s.blood_group,
		s.color,
		ch.checked_at
	FROM statuses s
	JOIN checks ch ON ch.id = s.check_id
	WHERE ch.center_id = $1
		AND s.blood_group = $2
		AND ch.success = TRUE
	ORDER BY ch.checked_at DESC, ch.id DESC;
	`

	rows, err := r.db.Query(query, centerID, group)
	if err != nil{
		return nil, fmt.Errorf("get history: %w", err)
	}
	defer rows.Close()

	statuses := []model.DonorStatus{}

	for rows.Next(){
		var status model.DonorStatus

		err := rows.Scan(
			&status.BloodGroup,
			&status.Color,
			&status.CheckedAt,
		)
		if err != nil{
			return nil, fmt.Errorf("Scan status by group: %w",err)
		}
		statuses = append(statuses, status)
	}
	if err := rows.Err(); err != nil{
		return nil, fmt.Errorf("read statuses by group: %w", err)
	}
	return statuses, nil
}


func (r *PostgresRepository) GetCenters() ([]model.DonorCenter, error){
	query :=  `
		SELECT id, name, city, source_url
		FROM centers
		ORDER BY id;
	`
	rows, err := r.db.Query(query)
	if err != nil{
		return nil, fmt.Errorf("get centers: %w", err)
	}

	defer rows.Close()

	centers := []model.DonorCenter{}

	for rows.Next(){
		var center model.DonorCenter

		err := rows.Scan(
			&center.ID,
			&center.NameCenter,
			&center.City,
			&center.URL,
		)
		if err != nil{
			return nil, fmt.Errorf("Scan center: %w", err)
		}
		centers = append(centers, center)
	}
	if err := rows.Err(); err != nil{
		return nil, fmt.Errorf("read centers: %w", err)
	}
	return centers, nil
}
