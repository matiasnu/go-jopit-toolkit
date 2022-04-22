package gosql

type QueryBuilder struct {
	Query  string
	Method string
	Models interface{}
}

type QueryBuilderResult struct {
	Models interface{}
	Error  error
}

func RunGenericQuery(queryBuilder QueryBuilder) QueryBuilderResult {
	var err error
	var models interface{}
	switch queryBuilder.Method {
	case "RAW":
		models, err = rawQueryBuild(queryBuilder.Query, queryBuilder.Models)
		return QueryBuilderResult{
			Models: models,
			Error:  err,
		}
	case "EXEC":
		err = execQueryBuild(queryBuilder.Query)
		return QueryBuilderResult{
			Error: err,
		}
	default:
		return QueryBuilderResult{}
	}
}
