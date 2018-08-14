package web

type ErrorResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Cause   string `json:"cause,omitempty"`
}

type GetMigrationRequest struct {
	IdMigration string `json:"id_migration" validate:"nonzero"`
}

type CreateMigrationRequest struct {
	Body struct {
		IdMigration string    `json:"id_migration" validate:"nonzero"`
	}
}

type DeleteMigrationRequest struct {
	IdMigration string `json:"id_migration"`
}
