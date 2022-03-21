package meta

type ConsumeMeta struct {
	ErrorPredicates predicateList[error]
}

func (cm ConsumeMeta) clone() ConsumeMeta {
	return ConsumeMeta{
		ErrorPredicates: cm.ErrorPredicates.clone(),
	}
}
