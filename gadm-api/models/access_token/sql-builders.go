package access_token

import "github.com/Masterminds/squirrel"

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)


func getAccessTokenSqlQuery(token string) (string, []interface{}, error) {
	sql, args, err := psql.
		Select("id", "token", "email", "created_at", "updated_at", "can_generate_access_tokens").
		From("access_tokens").
		Where(squirrel.Eq{"token": token}).
		ToSql()

	if err != nil {
		return "", nil, err
	}
	return sql, args, nil
}
